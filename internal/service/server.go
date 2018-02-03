package service

import (
	"net/http"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/if-ivan-else/tldrfeed/internal/db"
	"github.com/unrolled/render"
)

// Server captures runtime aspects of the tldrfeed
type Server struct {
	formatter *render.Render
	repo      db.Repository
	port      int
}

func NewServer(config Config) *Server {
	server := &Server{
		formatter: render.New(
			render.Options{
				IndentJSON: config.IndentJSON,
			},
		),
		port: config.Port,
		repo: db.NewRepository(),
	}

	return server
}

func (s *Server) Run() {
	n := negroni.Classic()
	// Run the server
	n.UseHandler(router(s))

	n.Run(":" + strconv.Itoa(s.port))
}

// Run starts and runs tldrfeed service
func Run(config Config) {
	NewServer(config).Run()
}

func router(s *Server) http.Handler {
	router := mux.NewRouter().PathPrefix(apiVersion).Subrouter()
	setupRoutes(router, s)
	return router
}

func setupRoutes(r *mux.Router, s *Server) {
	// User routes
	//
	// Create a User
	r.HandleFunc("/users", s.createUserHandler()).Methods("POST")
	// List Users
	r.HandleFunc("/users", s.getUserListHandler()).Methods("GET")
	// Get User
	r.HandleFunc("/users/{userID}", s.getUserHandler()).Methods("GET")

	// User feed and article retrieval
	//
	// Get all Feeds a Subscriber is following
	r.HandleFunc("/users/{userID}/feeds", s.getUserFeedListHandler()).Methods("GET")
	r.HandleFunc("/users/{userID}/feeds", s.addUserFeedHandler()).Methods("POST")

	r.HandleFunc("/users/{userID}/feeds/{feedID}/articles", s.getUserFeedArticleListHandler()).Methods("GET")
	r.HandleFunc("/users/{userID}/articles", s.getUserArticleListHandler()).Methods("GET")

	// Feed management routes
	//
	// List available Feeds
	r.HandleFunc("/feeds", s.getFeedListHandler()).Methods("GET")
	// Get a Feed
	r.HandleFunc("/feeds/{feedID}", s.getFeedHandler()).Methods("GET")

	// Create (sign up) a new Feed
	r.HandleFunc("/feeds", s.createFeedHandler()).Methods("POST")

	// Feed articles routes
	// List Articles in a Feed
	r.HandleFunc("/feeds/{feedID}/articles", s.getFeedArticleListHandler()).Methods("GET")
	// Add Articles to a Feed
	r.HandleFunc("/feeds/{feedID}/articles", s.createFeedArticleHandler()).Methods("POST")

}
