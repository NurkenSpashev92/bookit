package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

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

	var dateOfBirth interface{}
	if user.DateOfBirth != "" {
		dateOfBirth = user.DateOfBirth
	}

	var email interface{}
	if user.Email != "" {
		email = user.Email
	}

	var phoneNumber interface{}
	if user.PhoneNumber != "" {
		phoneNumber = user.PhoneNumber
	}

	err = r.db.QueryRow(ctx,
		`INSERT INTO users
			(email, first_name, last_name, middle_name, password, date_of_birth, phone_number, is_superuser, is_active, date_joined, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW(),NOW(),NOW())
		 RETURNING id,email,first_name,last_name,middle_name,password,date_of_birth,phone_number,is_superuser,is_active,date_joined,created_at,updated_at`,
		email, user.FirstName, user.LastName, user.MiddleName, string(hashedPassword), dateOfBirth, phoneNumber,
		false, true,
	).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.MiddleName, &u.Password, &u.DateOfBirth, &u.PhoneNumber,
		&u.IsSuperuser, &u.IsActive, &u.DateJoined, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				switch pgErr.ConstraintName {
				case "users_email_key", "users_email_unique", "ix_users_email":
					return u, fmt.Errorf("email already exists")
				case "users_phone_number_key", "users_phone_number_unique":
					return u, fmt.Errorf("phone number already exists")
				}
			}
		}
		return u, err
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, first_name, last_name, middle_name, password, phone_number, date_of_birth, COALESCE(avatar, ''), is_superuser, is_active
		 FROM users WHERE id=$1`, id,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.MiddleName, &user.Password, &user.PhoneNumber, &user.DateOfBirth, &user.Avatar, &user.IsSuperuser, &user.IsActive)
	if err != nil {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, userID int, req schemas.UserUpdateRequest) (models.User, error) {
	var user models.User

	// Fetch current user first
	err := r.db.QueryRow(ctx,
		`SELECT id, email, first_name, last_name, middle_name, password, phone_number, date_of_birth, COALESCE(avatar, ''), is_superuser, is_active
		 FROM users WHERE id=$1`, userID,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.MiddleName, &user.Password, &user.PhoneNumber, &user.DateOfBirth, &user.Avatar, &user.IsSuperuser, &user.IsActive)
	if err != nil {
		return user, errors.New("user not found")
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.MiddleName != nil {
		user.MiddleName = *req.MiddleName
	}
	if req.PhoneNumber != nil {
		if *req.PhoneNumber == "" {
			user.PhoneNumber = nil
		} else {
			user.PhoneNumber = req.PhoneNumber
		}
	}
	if req.DateOfBirth != nil {
		if *req.DateOfBirth == "" {
			user.DateOfBirth = nil
		} else {
			t, _ := time.Parse("2006-01-02", *req.DateOfBirth)
			user.DateOfBirth = &t
		}
	}

	_, err = r.db.Exec(ctx,
		`UPDATE users SET first_name=$1, last_name=$2, middle_name=$3, phone_number=$4, date_of_birth=$5, updated_at=NOW()
		 WHERE id=$6`,
		user.FirstName, user.LastName, user.MiddleName, user.PhoneNumber, user.DateOfBirth, userID,
	)
	if err != nil {
		return user, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByPhoneNumber(ctx context.Context, phone string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, first_name, last_name, middle_name, password, phone_number, date_of_birth, COALESCE(avatar, ''), is_superuser, is_active
		 FROM users WHERE phone_number=$1`, phone,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.MiddleName, &user.Password, &user.PhoneNumber, &user.DateOfBirth, &user.Avatar, &user.IsSuperuser, &user.IsActive)
	if err != nil {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET password=$1, updated_at=NOW() WHERE id=$2`,
		hashedPassword, userID,
	)
	return err
}

func (r *UserRepository) UpdateAvatar(ctx context.Context, userID int, avatar string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET avatar=$1, updated_at=NOW() WHERE id=$2`,
		avatar, userID,
	)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, first_name, last_name, middle_name, password, phone_number, date_of_birth, COALESCE(avatar, ''), is_superuser, is_active
		 FROM users WHERE email=$1`, email,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.MiddleName, &user.Password, &user.PhoneNumber, &user.DateOfBirth, &user.Avatar, &user.IsSuperuser, &user.IsActive)
	if err != nil {
		return user, errors.New("user not found")
	}
	return user, nil
}
