package test

import (
	"context"
	"fmt"
	"strings"

	contentschema "github.com/nurkenspashev92/bookit/internal/content/schema"
	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	identityschema "github.com/nurkenspashev92/bookit/internal/identity/schema"
	locationmodel "github.com/nurkenspashev92/bookit/internal/location/model"
	locationschema "github.com/nurkenspashev92/bookit/internal/location/schema"
	propertyschema "github.com/nurkenspashev92/bookit/internal/property/schema"
)

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
		Password: "$2a$10$fakehash",
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
		return identitymodel.User{}, identitymodel.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (identitymodel.User, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return identitymodel.User{}, identitymodel.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByPhoneNumber(_ context.Context, phone string) (identitymodel.User, error) {
	for _, u := range m.users {
		if u.PhoneNumber != nil && *u.PhoneNumber == phone {
			return u, nil
		}
	}
	return identitymodel.User{}, identitymodel.ErrUserNotFound
}

func (m *mockUserRepo) Update(_ context.Context, userID int, _ identityschema.UserUpdateRequest) (identitymodel.User, error) {
	u, ok := m.users[userID]
	if !ok {
		return identitymodel.User{}, identitymodel.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) UpdateUser(_ context.Context, id int, req identityschema.UserAdminUpdateRequest) (identitymodel.User, error) {
	u, ok := m.users[id]
	if !ok {
		return identitymodel.User{}, identitymodel.ErrUserNotFound
	}
	if req.Email != nil && *req.Email != u.Email {
		if _, exists := m.byEmail[*req.Email]; exists {
			return identitymodel.User{}, identitymodel.ErrEmailExists
		}
		delete(m.byEmail, u.Email)
		u.Email = *req.Email
	}
	if req.FirstName != nil {
		u.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		u.LastName = *req.LastName
	}
	if req.MiddleName != nil {
		u.MiddleName = *req.MiddleName
	}
	if req.PhoneNumber != nil {
		if *req.PhoneNumber == "" {
			u.PhoneNumber = nil
		} else {
			p := *req.PhoneNumber
			u.PhoneNumber = &p
		}
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}
	if req.IsSuperuser != nil {
		u.IsSuperuser = *req.IsSuperuser
	}
	m.users[id] = u
	m.byEmail[u.Email] = u
	return u, nil
}

func (m *mockUserRepo) UpdatePassword(_ context.Context, _ int, _ string) error {
	return nil
}

func (m *mockUserRepo) UpdateAvatar(_ context.Context, _ int, _ string) error {
	return nil
}

func (m *mockUserRepo) ListAll(_ context.Context, search string) ([]identitymodel.User, error) {
	var result []identitymodel.User
	for _, u := range m.users {
		if userMatchesSearch(u, search) {
			result = append(result, u)
		}
	}
	return result, nil
}

func (m *mockUserRepo) ListPaginated(_ context.Context, search string, limit, offset int) ([]identitymodel.User, int, error) {
	users, _ := m.ListAll(context.Background(), search)
	total := len(users)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return users[offset:end], total, nil
}

func userMatchesSearch(u identitymodel.User, search string) bool {
	if search == "" {
		return true
	}
	q := strings.ToLower(search)
	phone := ""
	if u.PhoneNumber != nil {
		phone = *u.PhoneNumber
	}
	for _, field := range []string{u.FirstName, u.LastName, u.Email, phone} {
		if strings.Contains(strings.ToLower(field), q) {
			return true
		}
	}
	return false
}

type mockHouseLikeRepo struct {
	likes  map[string]bool
	counts map[string]int
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

func (m *mockHouseLikeRepo) StatusWithCountByID(_ context.Context, _, _ int) (bool, int, error) {
	return false, 0, nil
}

func (m *mockHouseLikeRepo) GetUserLikedHouses(_ context.Context, _ int, _ string) ([]propertyschema.HouseListItem, error) {
	return []propertyschema.HouseListItem{}, nil
}

func (m *mockHouseLikeRepo) GetUserLikedHousesPaginated(_ context.Context, _ int, _ string, _, _ int) ([]propertyschema.HouseListItem, int, error) {
	return []propertyschema.HouseListItem{}, 0, nil
}

func (m *mockHouseLikeRepo) GetUserLikedHouseIDs(_ context.Context, _ int, _ []int) ([]int, error) {
	return nil, nil
}

type mockFAQRepo struct {
	faqs   map[int]contentschema.FAQ
	nextID int
}

func newMockFAQRepo() *mockFAQRepo {
	return &mockFAQRepo{faqs: make(map[int]contentschema.FAQ), nextID: 1}
}

func (m *mockFAQRepo) GetAll(_ context.Context, _ string) ([]contentschema.FAQ, error) {
	var result []contentschema.FAQ
	for _, f := range m.faqs {
		result = append(result, f)
	}
	return result, nil
}

func (m *mockFAQRepo) GetAllPaginated(_ context.Context, _ string, limit, offset int) ([]contentschema.FAQ, int, error) {
	var all []contentschema.FAQ
	for _, f := range m.faqs {
		all = append(all, f)
	}
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
	if req.QuestionKz != nil {
		f.QuestionKz = *req.QuestionKz
	}
	if req.AnswerKz != nil {
		f.AnswerKz = *req.AnswerKz
	}
	if req.QuestionRu != nil {
		f.QuestionRu = *req.QuestionRu
	}
	if req.AnswerRu != nil {
		f.AnswerRu = *req.AnswerRu
	}
	if req.QuestionEn != nil {
		f.QuestionEn = *req.QuestionEn
	}
	if req.AnswerEn != nil {
		f.AnswerEn = *req.AnswerEn
	}
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

type mockCountryRepo struct {
	countries map[int]locationmodel.Country
	nextID    int
}

func newMockCountryRepo() *mockCountryRepo {
	return &mockCountryRepo{countries: make(map[int]locationmodel.Country), nextID: 1}
}

func (m *mockCountryRepo) GetAll(_ context.Context, _ string) ([]locationmodel.Country, error) {
	var result []locationmodel.Country
	for _, c := range m.countries {
		result = append(result, c)
	}
	return result, nil
}

func (m *mockCountryRepo) GetAllPaginated(_ context.Context, _ string, limit, offset int) ([]locationmodel.Country, int, error) {
	var all []locationmodel.Country
	for _, c := range m.countries {
		all = append(all, c)
	}
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
	if req.NameKZ != nil {
		c.NameKZ = *req.NameKZ
	}
	if req.NameEN != nil {
		c.NameEN = *req.NameEN
	}
	if req.NameRU != nil {
		c.NameRU = *req.NameRU
	}
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
