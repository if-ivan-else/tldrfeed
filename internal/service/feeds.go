package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

// CreateUserRequest represents a request to create a new User
type CreateFeedRequest struct {
	Name string `json:"name" valid:"required~Feed name cannot be blank"`
}

// createFeedHandler creates a new Feed
func (s *Server) createFeedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		feedRequest := CreateFeedRequest{}
		if err := decodeAndValidate(req, &feedRequest); err != nil {
			s.formatter.Text(w, http.StatusBadRequest, err.Error())
			return
		}

		feed, err := s.repo.CreateFeed(feedRequest.Name)
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusCreated, feed)
	}
}

// getFeedListHandler returns the entire list of Feeds available for subscription
func (s *Server) getFeedListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		feeds, err := s.repo.ListFeeds()
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusOK, feeds)
	}
}

func (s *Server) getFeedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		feed, err := s.repo.GetFeedByID(vars["feedID"])
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusOK, feed)
	}
}

// getUserFeedListHandler returns all Feeds a User is following
func (s *Server) getUserFeedListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		feeds, err := s.repo.ListUserFeeds(vars["userID"])
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusOK, feeds)
	}
}

// addUserFeedHandler subscribes a User to a Feed
func (s *Server) addUserFeedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		err := s.repo.AddUserFeed(vars["userID"], vars["feedID"])
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}
		s.formatter.Text(w, http.StatusAccepted, "")

	}
}
