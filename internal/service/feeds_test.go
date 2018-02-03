package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateInvalidFeed(t *testing.T) {
	const textData = `not json`
	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/feeds", strings.NewReader(textData))
	rr := httptest.NewRecorder()
	server := testServer()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusBadRequest, rr.Result().StatusCode)
}

func TestCreateBlankNameFeed(t *testing.T) {
	const jsonData = `{}`
	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/feeds", strings.NewReader(jsonData))
	rr := httptest.NewRecorder()
	server := testServer()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusBadRequest, rr.Result().StatusCode)
	t.Logf("Response error returned: %s", rr.Body.String())
}
func TestCreateValidFeed(t *testing.T) {
	const jsonData = `{
    "name" : "Non-stop Tolstoy Channel"
  }
  `
	require := require.New(t)

	req, _ := http.NewRequest("POST", "/api/v1/feeds", strings.NewReader(jsonData))
	rr := httptest.NewRecorder()

	server := testServer()
	router(server).ServeHTTP(rr, req)

	require.Equal(http.StatusCreated, rr.Result().StatusCode)

	var respJSON map[string]string
	err := json.NewDecoder(rr.Result().Body).Decode(&respJSON)
	require.Nil(err)
	require.Contains(respJSON, "id")
	require.Equal("Non-stop Tolstoy Channel", respJSON["name"])
}
