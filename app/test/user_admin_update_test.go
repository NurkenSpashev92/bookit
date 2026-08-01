package test

import (
	"net/http"
	"testing"

	identityhandler "github.com/nurkenspashev92/bookit/internal/identity/handler"
	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	identityschema "github.com/nurkenspashev92/bookit/internal/identity/schema"
	identitysvc "github.com/nurkenspashev92/bookit/internal/identity/service"
)

func TestUserHandler_Update_PersistsProfileFields(t *testing.T) {
	repo := newMockUserRepo()
	phone := "+77001112233"
	repo.users[1] = identitymodel.User{ID: 1, Email: "old@mail.com", FirstName: "Old", IsActive: true, PhoneNumber: &phone}
	repo.byEmail["old@mail.com"] = repo.users[1]
	repo.users[2] = identitymodel.User{ID: 2, Email: "taken@mail.com", IsActive: true}
	repo.byEmail["taken@mail.com"] = repo.users[2]

	svc := identitysvc.NewUserService(repo, newTestJWTService(), testAwsConfig())
	h := identityhandler.NewUserHandler(svc)

	app := newTestApp()
	app.Patch("/users/:id", h.Update)

	resp := doRequest(t, app, http.MethodPatch, "/users/1", map[string]any{
		"first_name": "New",
		"last_name":  "Name",
		"email":      "new@mail.com",
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var body identityschema.AdminUser
	parseJSON(t, resp, &body)
	if body.FirstName != "New" {
		t.Errorf("response FirstName = %q, want New", body.FirstName)
	}
	if body.LastName != "Name" {
		t.Errorf("response LastName = %q, want Name", body.LastName)
	}
	if body.Email != "new@mail.com" {
		t.Errorf("response Email = %q, want new@mail.com", body.Email)
	}

	// persisted in the repository
	stored := repo.users[1]
	if stored.FirstName != "New" || stored.LastName != "Name" || stored.Email != "new@mail.com" {
		t.Errorf("stored user = %+v, want first_name=New last_name=Name email=new@mail.com", stored)
	}
}

func TestUserHandler_Update_DuplicateEmailConflict(t *testing.T) {
	repo := newMockUserRepo()
	repo.users[1] = identitymodel.User{ID: 1, Email: "old@mail.com", IsActive: true}
	repo.byEmail["old@mail.com"] = repo.users[1]
	repo.users[2] = identitymodel.User{ID: 2, Email: "taken@mail.com", IsActive: true}
	repo.byEmail["taken@mail.com"] = repo.users[2]

	svc := identitysvc.NewUserService(repo, newTestJWTService(), testAwsConfig())
	h := identityhandler.NewUserHandler(svc)

	app := newTestApp()
	app.Patch("/users/:id", h.Update)

	resp := doRequest(t, app, http.MethodPatch, "/users/1", map[string]any{
		"email": "taken@mail.com",
	})
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("status = %d, want 409", resp.StatusCode)
	}

	// original email unchanged
	if repo.users[1].Email != "old@mail.com" {
		t.Errorf("email changed on conflict = %q, want old@mail.com", repo.users[1].Email)
	}
}
