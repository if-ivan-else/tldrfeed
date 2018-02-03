package service

import (
	"net/http/httptest"

	"github.com/if-ivan-else/tldrfeed/internal/db/mock"
	"github.com/stretchr/testify/require"
	"github.com/unrolled/render"
)

func testServer() *Server {
	server := &Server{
		formatter: render.New(
			render.Options{
				IndentJSON: true,
			},
		),
		port: 8080,
		repo: mock.NewRepository(),
	}

	return server

}

func requireStatus(status int, require *require.Assertions, rr *httptest.ResponseRecorder) {
	require.Equal(status, rr.Result().StatusCode, "HTTP Error: %s: %s", rr.Result().Status, rr.Body.String())
}
