package service

import (
	"github.com/if-ivan-else/tldrfeed/internal/db/mock"
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
