package services

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type UserRepository interface {
	Create(ctx context.Context, req schemas.UserCreateRequest) (models.User, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByPhoneNumber(ctx context.Context, phone string) (models.User, error)
	Update(ctx context.Context, userID int, req schemas.UserUpdateRequest) (models.User, error)
	UpdatePassword(ctx context.Context, userID int, hashedPassword string) error
	UpdateAvatar(ctx context.Context, userID int, avatar string) error
}

type HouseRepository interface {
	GetAll(ctx context.Context) ([]schemas.HouseListItem, error)
	GetAllPaginated(ctx context.Context, filter schemas.HouseFilter, limit, offset int) ([]schemas.HouseListItem, int, error)
	GetByOwnerPaginated(ctx context.Context, ownerID, limit, offset int) ([]schemas.HouseListItem, int, error)
	GetBySlug(ctx context.Context, slug string) (schemas.HouseDetailResponse, error)
	RecordView(ctx context.Context, slug string, userID *int, ip string)
	Create(ctx context.Context, req schemas.HouseCreateRequest) (models.House, error)
	Update(ctx context.Context, slug string, req schemas.HouseUpdateRequest) (models.House, error)
	Delete(ctx context.Context, slug string) error
	SlugExists(ctx context.Context, slug string) (bool, error)
}

type HouseLikeRepository interface {
	LikeReturningCount(ctx context.Context, userID int, slug string) (int, error)
	UnlikeReturningCount(ctx context.Context, userID int, slug string) (int, error)
	StatusWithCount(ctx context.Context, userID int, slug string) (bool, int, error)
	GetUserLikedHouses(ctx context.Context, userID int) ([]schemas.HouseListItem, error)
	GetUserLikedHousesPaginated(ctx context.Context, userID, limit, offset int) ([]schemas.HouseListItem, int, error)
	GetUserLikedHouseIDs(ctx context.Context, userID int) ([]int, error)
}

type CountryRepository interface {
	GetAll(ctx context.Context) ([]models.Country, error)
	GetByID(ctx context.Context, id int) (models.Country, error)
	Create(ctx context.Context, req schemas.CountryCreateRequest) (models.Country, error)
	Update(ctx context.Context, id int, req schemas.CountryUpdateRequest) (models.Country, error)
	Delete(ctx context.Context, id int) error
}

type CityRepository interface {
	GetAllWithCountry(ctx context.Context) ([]schemas.City, error)
	GetByIDWithCountry(ctx context.Context, id int) (schemas.City, error)
	Create(ctx context.Context, req schemas.CityCreateRequest) (models.City, error)
	Update(ctx context.Context, id int, req schemas.CityUpdateRequest) (models.City, error)
	Delete(ctx context.Context, id int) error
}

type FAQRepository interface {
	GetAll(ctx context.Context) ([]schemas.FAQ, error)
	GetByID(ctx context.Context, id int) (schemas.FAQ, error)
	Create(ctx context.Context, req schemas.FAQCreateRequest) (schemas.FAQ, error)
	Update(ctx context.Context, id int, req schemas.FAQUpdateRequest) (schemas.FAQ, error)
	Delete(ctx context.Context, id int) error
}

type InquiryRepository interface {
	GetAll(ctx context.Context) ([]schemas.Inquiry, error)
	GetByID(ctx context.Context, id int) (schemas.Inquiry, error)
	Create(ctx context.Context, req schemas.InquiryCreateRequest) (schemas.Inquiry, error)
	Update(ctx context.Context, id int, req schemas.InquiryUpdateRequest) (schemas.Inquiry, error)
	Delete(ctx context.Context, id int) error
}
