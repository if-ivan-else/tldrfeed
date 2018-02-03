package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateInvalidUser(t *testing.T) {
	const textData = `not json`
	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(textData))
	rr := httptest.NewRecorder()
	server := testServer()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusBadRequest, rr.Result().StatusCode)
}

func TestCreateBlankNameUser(t *testing.T) {
	const jsonData = `{}`
	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	server := testServer()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusBadRequest, rr.Result().StatusCode)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestCreateNonAlphanumNameUser(t *testing.T) {
	const jsonData = `{
    "name" : "%$*&!!!"
  }`
	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	server := testServer()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusBadRequest, rr.Result().StatusCode)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestCreateValidUser(t *testing.T) {
	const jsonData = `{
    "name" : "boris"
  }
  `
	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(jsonData))
	rr := httptest.NewRecorder()

	server := testServer()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusCreated, rr.Result().StatusCode)

	var respJSON map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Contains(respJSON, "id")
	require.Equal("boris", respJSON["name"])
}

func TestListUsers(t *testing.T) {
	require := require.New(t)

	server := testServer()
	server.repo.CreateUser("natasha")

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusOK, rr.Result().StatusCode)

	var respJSON []map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Len(respJSON, 1)
	require.Contains(respJSON[0], "id")
	require.Equal("natasha", respJSON[0]["name"])
}

func TestGetUnknownUser(t *testing.T) {
	require := require.New(t)

	server := testServer()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", uuid.New().String()), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusNotFound, rr.Result().StatusCode)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestGetUser(t *testing.T) {
	require := require.New(t)

	server := testServer()
	user, _ := server.repo.CreateUser("alexey")

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", user.ID), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusOK, rr.Result().StatusCode)

	var respJSON map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Contains(respJSON, "id")
	require.Equal("alexey", respJSON["name"])
}
