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

func TestCreateInvalidFeed(t *testing.T) {
	const textData = `not json`
	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/feeds", strings.NewReader(textData))
	rr := httptest.NewRecorder()
	server := testServer()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusBadRequest, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestCreateBlankNameFeed(t *testing.T) {
	const jsonData = `{}`
	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/feeds", strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	server := testServer()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusBadRequest, require, rr)

	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestCreateValidFeed(t *testing.T) {
	name := "Non-stop Tolstoy Fun Channel"
	jsonData := fmt.Sprintf(`{
    "name" : "%s"
  }`, name)

	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/feeds", strings.NewReader(jsonData))
	rr := httptest.NewRecorder()

	server := testServer()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusCreated, require, rr)

	var respJSON map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Contains(respJSON, "id")
	require.Equal(name, respJSON["name"])
}

func TestListFeeds(t *testing.T) {
	require := require.New(t)

	server := testServer()
	name := "Anton Chekhov News"
	server.repo.CreateFeed(name)

	req, _ := http.NewRequest("GET", "/api/v1/feeds", nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusOK, require, rr)

	var respJSON []map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Len(respJSON, 1)
	require.Contains(respJSON[0], "id")
	require.Equal(name, respJSON[0]["name"])
}

func TestGetUnknownFeed(t *testing.T) {
	require := require.New(t)

	server := testServer()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/feeds/%s", uuid.New().String()), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusNotFound, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestGetFeed(t *testing.T) {
	require := require.New(t)

	server := testServer()
	name := "N.V. Gogol's Personal Blog"
	f, _ := server.repo.CreateFeed(name)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/feeds/%s", f.ID), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusOK, require, rr)

	var respJSON map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Contains(respJSON, "id")
	require.Equal(name, respJSON["name"])
}

func TestGetUnknownUserFeedList(t *testing.T) {
	require := require.New(t)

	server := testServer()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s/feeds", uuid.New().String()), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusNotFound, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestGetUserFeedList(t *testing.T) {
	require := require.New(t)

	server := testServer()
	u, _ := server.repo.CreateUser("victor")
	f, _ := server.repo.CreateFeed("Dostoevsky Daily")
	_ = server.repo.AddUserFeed(u.ID, f.ID)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s/feeds", u.ID), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusOK, require, rr)

	var respJSON []map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Len(respJSON, 1)
	require.Contains(respJSON[0], "id")
	require.Equal(f.Name, respJSON[0]["name"])
}

func TestGetUserUnknownFeed(t *testing.T) {
	require := require.New(t)

	server := testServer()
	u, _ := server.repo.CreateUser("victor")

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s/feeds/%s", u.ID, uuid.New().String()), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusNotFound, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestGetUserFeed(t *testing.T) {
	require := require.New(t)

	server := testServer()
	u, _ := server.repo.CreateUser("olga")
	f, _ := server.repo.CreateFeed("Rakhmaninov Folk Fairy Tales")
	_ = server.repo.AddUserFeed(u.ID, f.ID)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s/feeds/%s", u.ID, f.ID), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusOK, require, rr)

	var respJSON map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Contains(respJSON, "id")
	require.Equal(f.Name, respJSON["name"])
}

func TestAddUnknownUserFeed(t *testing.T) {
	require := require.New(t)

	server := testServer()
	f, _ := server.repo.CreateFeed("Chaikovsky Breaking News")

	jsonData := fmt.Sprintf(`{
    "feedID" : "%s"
  }`, f.ID)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/users/%s/feeds", uuid.New().String()), strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusNotFound, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestAddUserUnknownFeed(t *testing.T) {
	require := require.New(t)

	server := testServer()
	u, _ := server.repo.CreateUser("alexandra")
	feedID := uuid.New().String()

	jsonData := fmt.Sprintf(`{
    "feedID" : "%s"
  }`, feedID)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/users/%s/feeds", u.ID), strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusNotFound, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestAddUserFeedBadRequest(t *testing.T) {
	require := require.New(t)

	server := testServer()
	u, _ := server.repo.CreateUser("alexandra")

	jsonData := "not json"

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/users/%s/feeds", u.ID), strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusBadRequest, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestAddUserFeedNoID(t *testing.T) {
	require := require.New(t)

	server := testServer()
	u, _ := server.repo.CreateUser("alexandra")

	jsonData := "{}"

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/users/%s/feeds", u.ID), strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusBadRequest, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestAddUserFeed(t *testing.T) {
	require := require.New(t)

	server := testServer()
	u, _ := server.repo.CreateUser("alexandra")
	f, _ := server.repo.CreateFeed("Chaikovsky Breaking News")

	jsonData := fmt.Sprintf(`{
    "feedID" : "%s"
  }`, f.ID)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/users/%s/feeds", u.ID), strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusAccepted, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}
