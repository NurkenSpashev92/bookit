package test

import (
	"context"
	"fmt"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

// --- User Repository Mock ---

type mockUserRepo struct {
	users    map[int]models.User
	byEmail  map[string]models.User
	nextID   int
	createFn func(ctx context.Context, req schemas.UserCreateRequest) (models.User, error)
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:   make(map[int]models.User),
		byEmail: make(map[string]models.User),
		nextID:  1,
	}
}

func (m *mockUserRepo) Create(ctx context.Context, req schemas.UserCreateRequest) (models.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	if _, exists := m.byEmail[req.Email]; exists {
		return models.User{}, fmt.Errorf("email %s already exists", req.Email)
	}
	user := models.User{
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

func (m *mockUserRepo) GetByID(ctx context.Context, id int) (models.User, error) {
	u, ok := m.users[id]
	if !ok {
		return models.User{}, fmt.Errorf("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return models.User{}, fmt.Errorf("user not found")
	}
	return u, nil
}

// --- House Like Repository Mock ---

type mockHouseLikeRepo struct {
	likes map[string]bool // "userID:houseID"
	counts map[int]int    // houseID -> count
}

func newMockHouseLikeRepo() *mockHouseLikeRepo {
	return &mockHouseLikeRepo{
		likes:  make(map[string]bool),
		counts: make(map[int]int),
	}
}

func (m *mockHouseLikeRepo) key(userID, houseID int) string {
	return fmt.Sprintf("%d:%d", userID, houseID)
}

func (m *mockHouseLikeRepo) LikeReturningCount(_ context.Context, userID, houseID int) (int, error) {
	k := m.key(userID, houseID)
	if !m.likes[k] {
		m.likes[k] = true
		m.counts[houseID]++
	}
	return m.counts[houseID], nil
}

func (m *mockHouseLikeRepo) UnlikeReturningCount(_ context.Context, userID, houseID int) (int, error) {
	k := m.key(userID, houseID)
	if m.likes[k] {
		delete(m.likes, k)
		m.counts[houseID]--
	}
	return m.counts[houseID], nil
}

func (m *mockHouseLikeRepo) StatusWithCount(_ context.Context, userID, houseID int) (bool, int, error) {
	k := m.key(userID, houseID)
	return m.likes[k], m.counts[houseID], nil
}

func (m *mockHouseLikeRepo) GetUserLikedHouses(_ context.Context, _ int) ([]schemas.HouseLikeItem, error) {
	return []schemas.HouseLikeItem{}, nil
}

// --- FAQ Repository Mock ---

type mockFAQRepo struct {
	faqs   map[int]schemas.FAQ
	nextID int
}

func newMockFAQRepo() *mockFAQRepo {
	return &mockFAQRepo{faqs: make(map[int]schemas.FAQ), nextID: 1}
}

func (m *mockFAQRepo) GetAll(_ context.Context) ([]schemas.FAQ, error) {
	var result []schemas.FAQ
	for _, f := range m.faqs {
		result = append(result, f)
	}
	return result, nil
}

func (m *mockFAQRepo) GetByID(_ context.Context, id int) (schemas.FAQ, error) {
	f, ok := m.faqs[id]
	if !ok {
		return schemas.FAQ{}, fmt.Errorf("FAQ not found")
	}
	return f, nil
}

func (m *mockFAQRepo) Create(_ context.Context, req schemas.FAQCreateRequest) (schemas.FAQ, error) {
	f := schemas.FAQ{
		ID: m.nextID, QuestionKz: req.QuestionKz, AnswerKz: req.AnswerKz,
		QuestionRu: req.QuestionRu, AnswerRu: req.AnswerRu,
		QuestionEn: req.QuestionEn, AnswerEn: req.AnswerEn,
	}
	m.faqs[f.ID] = f
	m.nextID++
	return f, nil
}

func (m *mockFAQRepo) Update(_ context.Context, id int, req schemas.FAQUpdateRequest) (schemas.FAQ, error) {
	f, ok := m.faqs[id]
	if !ok {
		return schemas.FAQ{}, fmt.Errorf("FAQ not found")
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
	countries map[int]models.Country
	nextID    int
}

func newMockCountryRepo() *mockCountryRepo {
	return &mockCountryRepo{countries: make(map[int]models.Country), nextID: 1}
}

func (m *mockCountryRepo) GetAll(_ context.Context) ([]models.Country, error) {
	var result []models.Country
	for _, c := range m.countries {
		result = append(result, c)
	}
	return result, nil
}

func (m *mockCountryRepo) GetByID(_ context.Context, id int) (models.Country, error) {
	c, ok := m.countries[id]
	if !ok {
		return models.Country{}, fmt.Errorf("country not found")
	}
	return c, nil
}

func (m *mockCountryRepo) Create(_ context.Context, req schemas.CountryCreateRequest) (models.Country, error) {
	c := models.Country{ID: m.nextID, NameKZ: req.NameKZ, NameEN: req.NameEN, NameRU: req.NameRU, Code: req.Code}
	m.countries[c.ID] = c
	m.nextID++
	return c, nil
}

func (m *mockCountryRepo) Update(_ context.Context, id int, req schemas.CountryUpdateRequest) (models.Country, error) {
	c, ok := m.countries[id]
	if !ok {
		return models.Country{}, fmt.Errorf("country not found")
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
