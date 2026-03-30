package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user schemas.UserCreateRequest) (models.User, error) {
	var u models.User

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return u, fmt.Errorf("failed to hash password: %w", err)
	}

	err = r.db.QueryRow(ctx,
		`INSERT INTO users 
			(email, first_name, last_name, middle_name, password, date_of_birth, phone_number, is_superuser, is_active, date_joined, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW(),NOW(),NOW())
		 RETURNING id,email,first_name,last_name,middle_name,password,date_of_birth,phone_number,is_superuser,is_active,date_joined,created_at,updated_at`,
		user.Email, user.FirstName, user.LastName, user.MiddleName, string(hashedPassword), user.DateOfBirth, user.PhoneNumber,
		false, true,
	).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.MiddleName, &u.Password, &u.DateOfBirth, &u.PhoneNumber,
		&u.IsSuperuser, &u.IsActive, &u.DateJoined, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				return u, fmt.Errorf("email %s already exists", user.Email)
			}
		}
		return u, err
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, first_name, last_name, middle_name, password, phone_number, date_of_birth, is_superuser, is_active
		 FROM users WHERE id=$1`, id,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.MiddleName, &user.Password, &user.PhoneNumber, &user.DateOfBirth, &user.IsSuperuser, &user.IsActive)
	if err != nil {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, first_name, last_name, middle_name, password, phone_number, date_of_birth, is_superuser, is_active
		 FROM users WHERE email=$1`, email,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.MiddleName, &user.Password, &user.PhoneNumber, &user.DateOfBirth, &user.IsSuperuser, &user.IsActive)
	if err != nil {
		return user, errors.New("user not found")
	}
	return user, nil
}
