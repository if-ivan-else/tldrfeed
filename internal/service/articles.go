package service

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

// createArticleRequest defines a request to add an Article to a Feed
type createArticleRequest struct {
	Title string `json:"title" valid:"required~Article title cannot be blank"`
	Body  string `json:"body" valid:"required~Article title cannot be blank"`
}

// createArticleResponse defines a response to send for adding an Article to a Feed
type createArticleResponse struct {
	ID string `json:"id"`
}

func (s *Server) createFeedArticleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		articleRequest := createArticleRequest{}
		if err := decodeAndValidate(req, &articleRequest); err != nil {
			s.formatter.Text(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(req)
		spew.Dump(vars)
		articleID, err := s.repo.CreateFeedArticle(vars["feedID"], articleRequest.Title, articleRequest.Body)
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		response := createArticleResponse{
			ID: articleID,
		}
		s.formatter.JSON(w, http.StatusCreated, response)
	}
}

func (s *Server) getUserFeedArticleListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		articles, err := s.repo.ListUserFeedArticles(vars["userID"], vars["feedID"])
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusOK, articles)
	}
}

func (s *Server) getUserArticleListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		articles, err := s.repo.ListUserArticles(vars["userID"])
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusOK, articles)
	}
}

func (s *Server) getFeedArticleListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		articles, err := s.repo.ListFeedArticles(vars["feedID"])
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusOK, articles)
	}
}
