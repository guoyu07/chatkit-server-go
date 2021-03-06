package chatkit

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRolePermissionsFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	rolePerms, err := testClient.GetRolePermissions("testRole", "testScope")
	assert.Error(t, err, "expected an error")
	assert.Empty(t, rolePerms.Permissions, "Should be empty")
}

func TestGetRolePermissionsSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `["update_user", "create_user"]`)
	}))
	defer testServer.Close()

	rolePerms, err := testClient.GetRolePermissions("testRole", "testScope")
	assert.NoError(t, err, "expected an error")
	assert.NotEmpty(t, rolePerms.Permissions, "Should be empty")
}

func TestUpdateRolePermissionsFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.UpdateRolePermissions("testRole", "testScope", UpdateRolePermissionsParams{
		[]string{"testPermission"},
		[]string{"perm"},
	})
	assert.Error(t, err, "expected an error")
}

func TestUpdateRolePermissionsSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	err := testClient.UpdateRolePermissions("testRole", "testScope", UpdateRolePermissionsParams{
		[]string{"testPermission"},
		nil,
	})
	assert.NoError(t, err, "expected an error")
}
