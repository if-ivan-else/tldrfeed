package service

import (
	"net/http/httptest"

	"github.com/if-ivan-else/tldrfeed/internal/db/mock"
	"github.com/stretchr/testify/require"
)

func testServer() *Server {
	return newServer(
		Config{
			IndentJSON: true,
			Port:       8080,
		},
		mock.NewRepository(),
	)
}

func requireStatus(status int, require *require.Assertions, rr *httptest.ResponseRecorder) {
	require.Equal(status, rr.Result().StatusCode, "HTTP Error: %s: %s", rr.Result().Status, rr.Body.String())
}
