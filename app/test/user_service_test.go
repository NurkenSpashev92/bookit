package test

import (
	"testing"

	"github.com/nurkenspashev92/bookit/configs"
	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	identitysvc "github.com/nurkenspashev92/bookit/internal/identity/service"
)

func testAwsConfig() *configs.AwsConfig {
	return configs.NewAwsConfig()
}

func TestUserMapper_ToAuthUser(t *testing.T) {
	mapper := identitysvc.UserMapper{}
	awsCfg := testAwsConfig()
	user := identitymodel.User{
		ID:         10,
		Email:      "test@mail.com",
		FirstName:  "Jane",
		LastName:   "Smith",
		MiddleName: "A",
		Avatar:     "pic.jpg",
		Password:   "should-not-appear",
	}

	auth := mapper.ToAuthUser(user, awsCfg)

	if auth.ID != 10 {
		t.Errorf("ID = %d, want 10", auth.ID)
	}
	if auth.Email != "test@mail.com" {
		t.Errorf("Email = %q, want test@mail.com", auth.Email)
	}
	if auth.FirstName != "Jane" {
		t.Errorf("FirstName = %q, want Jane", auth.FirstName)
	}
	if auth.LastName != "Smith" {
		t.Errorf("LastName = %q, want Smith", auth.LastName)
	}
	if auth.MiddleName != "A" {
		t.Errorf("MiddleName = %q, want A", auth.MiddleName)
	}
	if auth.Avatar == "" {
		t.Error("Avatar should not be empty")
	}
}

func TestUserMapper_ToAuthUser_EmptyFields(t *testing.T) {
	mapper := identitysvc.UserMapper{}
	awsCfg := testAwsConfig()
	user := identitymodel.User{ID: 1, Email: "x@y.com"}

	auth := mapper.ToAuthUser(user, awsCfg)

	if auth.FirstName != "" {
		t.Errorf("FirstName should be empty, got %q", auth.FirstName)
	}
	if auth.Avatar != "" {
		t.Errorf("Avatar should be empty, got %q", auth.Avatar)
	}
}
