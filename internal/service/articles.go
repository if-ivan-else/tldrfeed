package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/if-ivan-else/tldrfeed/api"
)

func (s *Server) createFeedArticleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		articleRequest := api.CreateArticleRequest{}
		if err := decodeAndValidate(req, &articleRequest); err != nil {
			s.formatter.Text(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(req)
		articleID, err := s.repo.CreateFeedArticle(vars["feedID"], articleRequest.Title, articleRequest.Body)
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		response := api.CreateArticleResponse{
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
