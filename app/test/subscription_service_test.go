package test

import (
	"context"
	"testing"
	"time"

	submodel "github.com/nurkenspashev92/bookit/internal/subscription/model"
	subschema "github.com/nurkenspashev92/bookit/internal/subscription/schema"
	subsvc "github.com/nurkenspashev92/bookit/internal/subscription/service"
)

// mockSubscriptionRepo mirrors the real repository's single-active-subscription
// invariant: creating (or updating to) an active subscription deactivates every
// other active subscription for the same user.
type mockSubscriptionRepo struct {
	subs   map[int]submodel.Subscription
	nextID int
}

func newMockSubscriptionRepo() *mockSubscriptionRepo {
	return &mockSubscriptionRepo{subs: make(map[int]submodel.Subscription), nextID: 1}
}

func (m *mockSubscriptionRepo) deactivateOthers(userID, exceptID int) {
	for id, s := range m.subs {
		if s.UserID == userID && id != exceptID && s.Status == submodel.StatusActive {
			s.Status = submodel.StatusInactive
			m.subs[id] = s
		}
	}
}

func (m *mockSubscriptionRepo) Create(_ context.Context, s submodel.Subscription) (submodel.Subscription, error) {
	if s.Status == submodel.StatusActive {
		m.deactivateOthers(s.UserID, 0)
	}
	s.ID = m.nextID
	m.nextID++
	m.subs[s.ID] = s
	return s, nil
}

func (m *mockSubscriptionRepo) GetByID(_ context.Context, id int) (submodel.Subscription, error) {
	s, ok := m.subs[id]
	if !ok {
		return submodel.Subscription{}, submodel.ErrSubscriptionNotFound
	}
	return s, nil
}

func (m *mockSubscriptionRepo) GetLatestByUserID(_ context.Context, userID int) (submodel.Subscription, error) {
	var latest submodel.Subscription
	found := false
	for _, s := range m.subs {
		if s.UserID == userID && (!found || s.ID > latest.ID) {
			latest = s
			found = true
		}
	}
	if !found {
		return submodel.Subscription{}, submodel.ErrSubscriptionNotFound
	}
	return latest, nil
}

func (m *mockSubscriptionRepo) GetActiveByUserID(_ context.Context, userID int) (submodel.Subscription, error) {
	var active submodel.Subscription
	found := false
	for _, s := range m.subs {
		if s.UserID == userID && s.Status == submodel.StatusActive && (!found || s.ID > active.ID) {
			active = s
			found = true
		}
	}
	if !found {
		return submodel.Subscription{}, submodel.ErrSubscriptionNotFound
	}
	return active, nil
}

func (m *mockSubscriptionRepo) Deactivate(_ context.Context, id int) error {
	s, ok := m.subs[id]
	if !ok {
		return submodel.ErrSubscriptionNotFound
	}
	s.Status = submodel.StatusInactive
	m.subs[id] = s
	return nil
}

func (m *mockSubscriptionRepo) ListByUserID(_ context.Context, userID int) ([]submodel.Subscription, error) {
	var result []submodel.Subscription
	for _, s := range m.subs {
		if s.UserID == userID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockSubscriptionRepo) ListAll(_ context.Context, _ string) ([]submodel.Subscription, error) {
	var result []submodel.Subscription
	for _, s := range m.subs {
		result = append(result, s)
	}
	return result, nil
}

func (m *mockSubscriptionRepo) ListPaginated(_ context.Context, _ string, limit, offset int) ([]submodel.Subscription, int, error) {
	all, _ := m.ListAll(context.Background(), "")
	total := len(all)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (m *mockSubscriptionRepo) Update(_ context.Context, id int, s submodel.Subscription) (submodel.Subscription, error) {
	if _, ok := m.subs[id]; !ok {
		return submodel.Subscription{}, submodel.ErrSubscriptionNotFound
	}
	s.ID = id
	m.subs[id] = s
	if s.Status == submodel.StatusActive {
		m.deactivateOthers(s.UserID, id)
	}
	return m.subs[id], nil
}

func (m *mockSubscriptionRepo) Delete(_ context.Context, id int) (int, error) {
	s, ok := m.subs[id]
	if !ok {
		return 0, submodel.ErrSubscriptionNotFound
	}
	delete(m.subs, id)
	return s.UserID, nil
}

// mockUserSync records the last subscription_type synced per user.
type mockUserSync struct {
	types map[int]string
}

func newMockUserSync() *mockUserSync {
	return &mockUserSync{types: make(map[int]string)}
}

func (m *mockUserSync) SetSubscriptionType(_ context.Context, userID int, subType string) error {
	m.types[userID] = subType
	return nil
}

func aboutOneMonth(t *testing.T, start time.Time, end *time.Time) {
	t.Helper()
	if end == nil {
		t.Fatal("end_date is nil, want start_date + 1 month")
	}
	want := start.AddDate(0, 1, 0)
	if diff := end.Sub(want); diff > time.Second || diff < -time.Second {
		t.Errorf("end_date = %v, want ~%v (start %v + 1 month)", *end, want, start)
	}
}

func TestSubscriptionService_Create_DefaultsEndDateOneMonth(t *testing.T) {
	repo := newMockSubscriptionRepo()
	svc := subsvc.NewSubscriptionService(repo, newMockUserSync())

	res, err := svc.Create(context.Background(), subschema.SubscriptionCreateRequest{
		UserID: 53, Type: "max", Status: "active",
	})
	if err != nil {
		t.Fatal(err)
	}
	aboutOneMonth(t, res.StartDate, res.EndDate)
}

func TestSubscriptionService_Create_RespectsProvidedEndDate(t *testing.T) {
	repo := newMockSubscriptionRepo()
	svc := subsvc.NewSubscriptionService(repo, newMockUserSync())

	provided := time.Date(2030, 6, 15, 12, 0, 0, 0, time.UTC)
	res, err := svc.Create(context.Background(), subschema.SubscriptionCreateRequest{
		UserID: 7, Type: "pro", Status: "active", EndDate: &provided,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.EndDate == nil || !res.EndDate.Equal(provided) {
		t.Errorf("end_date = %v, want provided %v", res.EndDate, provided)
	}
}

func TestSubscriptionService_Create_SecondActiveDeactivatesFirst(t *testing.T) {
	repo := newMockSubscriptionRepo()
	sync := newMockUserSync()
	svc := subsvc.NewSubscriptionService(repo, sync)
	ctx := context.Background()

	first, err := svc.Create(ctx, subschema.SubscriptionCreateRequest{UserID: 53, Type: "pro", Status: "active"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.Create(ctx, subschema.SubscriptionCreateRequest{UserID: 53, Type: "max", Status: "active"})
	if err != nil {
		t.Fatal(err)
	}

	firstAfter, err := repo.GetByID(ctx, first.ID)
	if err != nil {
		t.Fatal(err)
	}
	if firstAfter.Status != submodel.StatusInactive {
		t.Errorf("first subscription status = %q, want in_active", firstAfter.Status)
	}
	secondAfter, err := repo.GetByID(ctx, second.ID)
	if err != nil {
		t.Fatal(err)
	}
	if secondAfter.Status != submodel.StatusActive {
		t.Errorf("second subscription status = %q, want active", secondAfter.Status)
	}

	// The freshly created active subscription wins the users.subscription_type sync.
	if sync.types[53] != "max" {
		t.Errorf("synced subscription_type = %q, want max", sync.types[53])
	}
}

func TestSubscriptionService_Update_ToActiveDeactivatesOthers(t *testing.T) {
	repo := newMockSubscriptionRepo()
	sync := newMockUserSync()
	svc := subsvc.NewSubscriptionService(repo, sync)
	ctx := context.Background()

	// Existing active subscription for the user.
	active, err := svc.Create(ctx, subschema.SubscriptionCreateRequest{UserID: 9, Type: "pro", Status: "active"})
	if err != nil {
		t.Fatal(err)
	}
	// A second, inactive subscription that we will flip to active via Update.
	inactive, err := svc.Create(ctx, subschema.SubscriptionCreateRequest{UserID: 9, Type: "max", Status: "in_active"})
	if err != nil {
		t.Fatal(err)
	}

	activeStatus := string(submodel.StatusActive)
	if _, err := svc.Update(ctx, inactive.ID, subschema.SubscriptionUpdateRequest{Status: &activeStatus}); err != nil {
		t.Fatal(err)
	}

	firstAfter, _ := repo.GetByID(ctx, active.ID)
	if firstAfter.Status != submodel.StatusInactive {
		t.Errorf("previously active subscription status = %q, want in_active", firstAfter.Status)
	}
	secondAfter, _ := repo.GetByID(ctx, inactive.ID)
	if secondAfter.Status != submodel.StatusActive {
		t.Errorf("updated subscription status = %q, want active", secondAfter.Status)
	}
	if sync.types[9] != "max" {
		t.Errorf("synced subscription_type = %q, want max", sync.types[9])
	}
}
