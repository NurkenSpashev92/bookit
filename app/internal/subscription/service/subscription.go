package service

import (
	"context"
	"errors"
	"time"

	"github.com/nurkenspashev92/bookit/internal/subscription/model"
	"github.com/nurkenspashev92/bookit/internal/subscription/port"
	"github.com/nurkenspashev92/bookit/internal/subscription/schema"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s model.Subscription) (model.Subscription, error)
	GetByID(ctx context.Context, id int) (model.Subscription, error)
	GetLatestByUserID(ctx context.Context, userID int) (model.Subscription, error)
	ListByUserID(ctx context.Context, userID int) ([]model.Subscription, error)
	ListAll(ctx context.Context) ([]model.Subscription, error)
	ListPaginated(ctx context.Context, limit, offset int) ([]model.Subscription, int, error)
	Update(ctx context.Context, id int, s model.Subscription) (model.Subscription, error)
	Delete(ctx context.Context, id int) (int, error)
}

type SubscriptionService struct {
	repository SubscriptionRepository
	userSync   port.UserSubscriptionSetter
}

func NewSubscriptionService(repo SubscriptionRepository, userSync port.UserSubscriptionSetter) *SubscriptionService {
	return &SubscriptionService{repository: repo, userSync: userSync}
}

func (s *SubscriptionService) GetMy(ctx context.Context, userID int) (schema.SubscriptionResponse, error) {
	sub, err := s.repository.GetLatestByUserID(ctx, userID)
	if err != nil {
		return schema.SubscriptionResponse{}, err
	}
	return toResponse(sub), nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id int) (schema.SubscriptionResponse, error) {
	sub, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return schema.SubscriptionResponse{}, err
	}
	return toResponse(sub), nil
}

func (s *SubscriptionService) GetByUserID(ctx context.Context, userID int) ([]schema.SubscriptionResponse, error) {
	subs, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return toResponses(subs), nil
}

func (s *SubscriptionService) GetAll(ctx context.Context) ([]schema.SubscriptionResponse, error) {
	subs, err := s.repository.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return toResponses(subs), nil
}

func (s *SubscriptionService) GetAllPaginated(ctx context.Context, limit, offset int) ([]schema.SubscriptionResponse, int, error) {
	subs, total, err := s.repository.ListPaginated(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return toResponses(subs), total, nil
}

func (s *SubscriptionService) Create(ctx context.Context, req schema.SubscriptionCreateRequest) (schema.SubscriptionResponse, error) {
	sub := model.Subscription{
		UserID:    req.UserID,
		Type:      orDefault(req.Type, model.TypeBasic),
		Status:    orDefault(req.Status, model.StatusActive),
		StartDate: time.Now(),
		EndDate:   req.EndDate,
	}

	created, err := s.repository.Create(ctx, sub)
	if err != nil {
		return schema.SubscriptionResponse{}, err
	}

	if err := s.syncUser(ctx, created.UserID, effectiveTier(created)); err != nil {
		return schema.SubscriptionResponse{}, err
	}
	return toResponse(created), nil
}

func (s *SubscriptionService) Update(ctx context.Context, id int, req schema.SubscriptionUpdateRequest) (schema.SubscriptionResponse, error) {
	sub, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return schema.SubscriptionResponse{}, err
	}

	if req.Type != nil {
		sub.Type = model.Type(*req.Type)
	}
	if req.Status != nil {
		sub.Status = model.Status(*req.Status)
	}
	if req.EndDate != nil {
		sub.EndDate = req.EndDate
	}

	updated, err := s.repository.Update(ctx, id, sub)
	if err != nil {
		return schema.SubscriptionResponse{}, err
	}

	if err := s.syncUser(ctx, updated.UserID, effectiveTier(updated)); err != nil {
		return schema.SubscriptionResponse{}, err
	}
	return toResponse(updated), nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id int) error {
	userID, err := s.repository.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Recompute the user's tier from their remaining latest subscription,
	// falling back to basic when none is left.
	tier := string(model.TypeBasic)
	if latest, err := s.repository.GetLatestByUserID(ctx, userID); err == nil {
		tier = effectiveTier(latest)
	} else if !errors.Is(err, model.ErrSubscriptionNotFound) {
		return err
	}

	return s.syncUser(ctx, userID, tier)
}

// syncUser mirrors the effective tier into users.subscription_type.
func (s *SubscriptionService) syncUser(ctx context.Context, userID int, tier string) error {
	if s.userSync == nil {
		return nil
	}
	return s.userSync.SetSubscriptionType(ctx, userID, tier)
}

// effectiveTier is the plan the user actually gets: the subscription type while
// active, otherwise basic.
func effectiveTier(sub model.Subscription) string {
	if sub.Status == model.StatusActive {
		return string(sub.Type)
	}
	return string(model.TypeBasic)
}

// orDefault converts a request enum string to its typed form, falling back to
// def when the field was omitted.
func orDefault[T ~string](v string, def T) T {
	if v == "" {
		return def
	}
	return T(v)
}

func toResponse(s model.Subscription) schema.SubscriptionResponse {
	return schema.SubscriptionResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		Type:      string(s.Type),
		Status:    string(s.Status),
		StartDate: s.StartDate,
		EndDate:   s.EndDate,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func toResponses(subs []model.Subscription) []schema.SubscriptionResponse {
	result := make([]schema.SubscriptionResponse, 0, len(subs))
	for i := range subs {
		result = append(result, toResponse(subs[i]))
	}
	return result
}
