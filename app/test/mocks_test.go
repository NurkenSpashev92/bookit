package test

import (
	"context"
	"fmt"

	contentschema "github.com/nurkenspashev92/bookit/internal/content/schema"
	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	identityschema "github.com/nurkenspashev92/bookit/internal/identity/schema"
	locationmodel "github.com/nurkenspashev92/bookit/internal/location/model"
	locationschema "github.com/nurkenspashev92/bookit/internal/location/schema"
	propertyschema "github.com/nurkenspashev92/bookit/internal/property/schema"
)

// --- User Repository Mock ---

type mockUserRepo struct {
	users    map[int]identitymodel.User
	byEmail  map[string]identitymodel.User
	nextID   int
	createFn func(ctx context.Context, req identityschema.UserCreateRequest) (identitymodel.User, error)
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:   make(map[int]identitymodel.User),
		byEmail: make(map[string]identitymodel.User),
		nextID:  1,
	}
}

func (m *mockUserRepo) Create(ctx context.Context, req identityschema.UserCreateRequest) (identitymodel.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	if _, exists := m.byEmail[req.Email]; exists {
		return identitymodel.User{}, fmt.Errorf("email %s already exists", req.Email)
	}
	user := identitymodel.User{
		ID:       m.nextID,
		Email:    req.Email,
		Password: "$2a$10$fakehash", // fake bcrypt
		IsActive: true,
	}
	m.nextID++
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	return user, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int) (identitymodel.User, error) {
	u, ok := m.users[id]
	if !ok {
		return identitymodel.User{}, fmt.Errorf("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (identitymodel.User, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return identitymodel.User{}, fmt.Errorf("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetByPhoneNumber(_ context.Context, phone string) (identitymodel.User, error) {
	for _, u := range m.users {
		if u.PhoneNumber != nil && *u.PhoneNumber == phone {
			return u, nil
		}
	}
	return identitymodel.User{}, fmt.Errorf("user not found")
}

func (m *mockUserRepo) Update(_ context.Context, userID int, _ identityschema.UserUpdateRequest) (identitymodel.User, error) {
	u, ok := m.users[userID]
	if !ok {
		return identitymodel.User{}, fmt.Errorf("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) UpdatePassword(_ context.Context, _ int, _ string) error {
	return nil
}

func (m *mockUserRepo) UpdateAvatar(_ context.Context, _ int, _ string) error {
	return nil
}

// --- House Like Repository Mock ---

type mockHouseLikeRepo struct {
	likes  map[string]bool   // "userID:slug"
	counts map[string]int    // slug -> count
}

func newMockHouseLikeRepo() *mockHouseLikeRepo {
	return &mockHouseLikeRepo{
		likes:  make(map[string]bool),
		counts: make(map[string]int),
	}
}

func (m *mockHouseLikeRepo) key(userID int, slug string) string {
	return fmt.Sprintf("%d:%s", userID, slug)
}

func (m *mockHouseLikeRepo) LikeReturningCount(_ context.Context, userID int, slug string) (int, error) {
	k := m.key(userID, slug)
	if !m.likes[k] {
		m.likes[k] = true
		m.counts[slug]++
	}
	return m.counts[slug], nil
}

func (m *mockHouseLikeRepo) UnlikeReturningCount(_ context.Context, userID int, slug string) (int, error) {
	k := m.key(userID, slug)
	if m.likes[k] {
		delete(m.likes, k)
		m.counts[slug]--
	}
	return m.counts[slug], nil
}

func (m *mockHouseLikeRepo) StatusWithCount(_ context.Context, userID int, slug string) (bool, int, error) {
	k := m.key(userID, slug)
	return m.likes[k], m.counts[slug], nil
}

func (m *mockHouseLikeRepo) GetUserLikedHouses(_ context.Context, _ int) ([]propertyschema.HouseListItem, error) {
	return []propertyschema.HouseListItem{}, nil
}

func (m *mockHouseLikeRepo) GetUserLikedHousesPaginated(_ context.Context, _ int, _, _ int) ([]propertyschema.HouseListItem, int, error) {
	return []propertyschema.HouseListItem{}, 0, nil
}

func (m *mockHouseLikeRepo) GetUserLikedHouseIDs(_ context.Context, userID int) ([]int, error) {
	var ids []int
	for k := range m.likes {
		var uid int
		var slug string
		fmt.Sscanf(k, "%d:%s", &uid, &slug)
		if uid == userID {
			ids = append(ids, uid)
		}
	}
	return ids, nil
}

// --- FAQ Repository Mock ---

type mockFAQRepo struct {
	faqs   map[int]contentschema.FAQ
	nextID int
}

func newMockFAQRepo() *mockFAQRepo {
	return &mockFAQRepo{faqs: make(map[int]contentschema.FAQ), nextID: 1}
}

func (m *mockFAQRepo) GetAll(_ context.Context) ([]contentschema.FAQ, error) {
	var result []contentschema.FAQ
	for _, f := range m.faqs {
		result = append(result, f)
	}
	return result, nil
}

func (m *mockFAQRepo) GetByID(_ context.Context, id int) (contentschema.FAQ, error) {
	f, ok := m.faqs[id]
	if !ok {
		return contentschema.FAQ{}, fmt.Errorf("FAQ not found")
	}
	return f, nil
}

func (m *mockFAQRepo) Create(_ context.Context, req contentschema.FAQCreateRequest) (contentschema.FAQ, error) {
	f := contentschema.FAQ{
		ID: m.nextID, QuestionKz: req.QuestionKz, AnswerKz: req.AnswerKz,
		QuestionRu: req.QuestionRu, AnswerRu: req.AnswerRu,
		QuestionEn: req.QuestionEn, AnswerEn: req.AnswerEn,
	}
	m.faqs[f.ID] = f
	m.nextID++
	return f, nil
}

func (m *mockFAQRepo) Update(_ context.Context, id int, req contentschema.FAQUpdateRequest) (contentschema.FAQ, error) {
	f, ok := m.faqs[id]
	if !ok {
		return contentschema.FAQ{}, fmt.Errorf("FAQ not found")
	}
	if req.QuestionKz != nil { f.QuestionKz = *req.QuestionKz }
	if req.AnswerKz != nil { f.AnswerKz = *req.AnswerKz }
	if req.QuestionRu != nil { f.QuestionRu = *req.QuestionRu }
	if req.AnswerRu != nil { f.AnswerRu = *req.AnswerRu }
	if req.QuestionEn != nil { f.QuestionEn = *req.QuestionEn }
	if req.AnswerEn != nil { f.AnswerEn = *req.AnswerEn }
	m.faqs[id] = f
	return f, nil
}

func (m *mockFAQRepo) Delete(_ context.Context, id int) error {
	if _, ok := m.faqs[id]; !ok {
		return fmt.Errorf("FAQ not found")
	}
	delete(m.faqs, id)
	return nil
}

// --- Country Repository Mock ---

type mockCountryRepo struct {
	countries map[int]locationmodel.Country
	nextID    int
}

func newMockCountryRepo() *mockCountryRepo {
	return &mockCountryRepo{countries: make(map[int]locationmodel.Country), nextID: 1}
}

func (m *mockCountryRepo) GetAll(_ context.Context) ([]locationmodel.Country, error) {
	var result []locationmodel.Country
	for _, c := range m.countries {
		result = append(result, c)
	}
	return result, nil
}

func (m *mockCountryRepo) GetByID(_ context.Context, id int) (locationmodel.Country, error) {
	c, ok := m.countries[id]
	if !ok {
		return locationmodel.Country{}, fmt.Errorf("country not found")
	}
	return c, nil
}

func (m *mockCountryRepo) Create(_ context.Context, req locationschema.CountryCreateRequest) (locationmodel.Country, error) {
	c := locationmodel.Country{ID: m.nextID, NameKZ: req.NameKZ, NameEN: req.NameEN, NameRU: req.NameRU, Code: req.Code}
	m.countries[c.ID] = c
	m.nextID++
	return c, nil
}

func (m *mockCountryRepo) Update(_ context.Context, id int, req locationschema.CountryUpdateRequest) (locationmodel.Country, error) {
	c, ok := m.countries[id]
	if !ok {
		return locationmodel.Country{}, fmt.Errorf("country not found")
	}
	if req.NameKZ != nil { c.NameKZ = *req.NameKZ }
	if req.NameEN != nil { c.NameEN = *req.NameEN }
	if req.NameRU != nil { c.NameRU = *req.NameRU }
	m.countries[id] = c
	return c, nil
}

func (m *mockCountryRepo) Delete(_ context.Context, id int) error {
	if _, ok := m.countries[id]; !ok {
		return fmt.Errorf("country not found")
	}
	delete(m.countries, id)
	return nil
}
