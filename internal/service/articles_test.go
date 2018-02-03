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

func TestCreateInvalidArticle(t *testing.T) {
	require := require.New(t)
	const textData = `not json`

	server := testServer()
	f, _ := server.repo.CreateFeed("Non-stop Tolstoy Fun Channel")

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/feeds/%s/articles", f.ID), strings.NewReader(textData))
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusBadRequest, require, rr)
}

func TestCreateBlankArticle(t *testing.T) {
	require := require.New(t)

	server := testServer()
	f, _ := server.repo.CreateFeed("Non-stop Tolstoy Fun Channel")

	const jsonData = `{}`
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/feeds/%s/articles", f.ID), strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusBadRequest, require, rr)

	t.Logf("Error message: %s", rr.Body.String())
}

func TestCreateValidArticle(t *testing.T) {
	require := require.New(t)
	server := testServer()
	f, _ := server.repo.CreateFeed("Non-stop Tolstoy Fun Channel")

	jsonData := `{
    "title": "War and Peace: Chapter 7",
    "body": "Anna Mikhailovna returns to her rich relations in Moscow in Chapter 7, where she and her son have lived for years..."
  }`

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/feeds/%s/articles", f.ID), strings.NewReader(jsonData))
	rr := httptest.NewRecorder()

	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusCreated, require, rr)

	var respJSON map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Contains(respJSON, "id")
}

func TestCreateUnknownFeedArticle(t *testing.T) {
	require := require.New(t)
	server := testServer()

	jsonData := `{
    "title": "War and Peace: Chapter 7",
    "body": "Anna Mikhailovna returns to her rich relations in Moscow in Chapter 7, where she and her son have lived for years..."
  }`

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/feeds/%s/articles", uuid.New().String()), strings.NewReader(jsonData))
	rr := httptest.NewRecorder()

	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusNotFound, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}
func TestListUnknownFeedArticles(t *testing.T) {
	require := require.New(t)

	server := testServer()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/feeds/%s/articles", uuid.New().String()), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusNotFound, require, rr)
	t.Logf("Error message (expected): %s", rr.Body.String())
}

func TestListFeedArticles(t *testing.T) {
	require := require.New(t)

	server := testServer()
	name := "Anton Chekhov Super Short Stories"
	f, _ := server.repo.CreateFeed(name)
	title := "Gooseberries"
	body := `Ivan Ivanovich Chimsha-Gimalayski, a veterinary surgeon,
tells the story of his younger brother Nikolai Ivanovich.`
	server.repo.CreateFeedArticle(f.ID, title, body)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/feeds/%s/articles", f.ID), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusOK, require, rr)

	var respJSON []map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.NoError(err)
	require.Len(respJSON, 1)
	require.Contains(respJSON[0], "id")
	require.Equal(title, respJSON[0]["title"])
	require.Equal(body, respJSON[0]["body"])
}

func TestListUnknownUserArticles(t *testing.T) {
	require := require.New(t)

	server := testServer()
	name := "Anton Chekhov Super Short Stories"
	f, _ := server.repo.CreateFeed(name)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s/feeds/%s/articles", uuid.New().String(), f.ID), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireStatus(http.StatusNotFound, require, rr)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s/articles", uuid.New().String()), nil)
	rr = httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)
	requireStatus(http.StatusNotFound, require, rr)
}

func TestListUserArticles(t *testing.T) {
	require := require.New(t)

	server := testServer()
	name := "Anton Chekhov Super Short Stories"
	f, _ := server.repo.CreateFeed(name)
	title := "A Boring Story"
	body := `Nikolai Stepanovich, a luminary in the world of medical science,
tormented by insomnia and bouts of devastating weakness,
lives in a kind of darkening haze.`

	server.repo.CreateFeedArticle(f.ID, title, body)
	u, _ := server.repo.CreateUser("alexey")
	server.repo.AddUserFeed(u.ID, f.ID)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s/feeds/%s/articles", u.ID, f.ID), nil)
	rr := httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireArticleJSON := func(rr *httptest.ResponseRecorder) {
		requireStatus(http.StatusOK, require, rr)
		var respJSON []map[string]string
		err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
		require.NoError(err)
		require.Len(respJSON, 1)
		require.Contains(respJSON[0], "id")
		require.Equal(title, respJSON[0]["title"])
		require.Equal(body, respJSON[0]["body"])
	}

	requireArticleJSON(rr)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s/articles", u.ID), nil)
	rr = httptest.NewRecorder()
	router(server).ServeHTTP(rr, req)

	requireArticleJSON(rr)
}
