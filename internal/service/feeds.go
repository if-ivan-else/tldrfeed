package service

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/if-ivan-else/tldrfeed/api"
)

// createFeedHandler creates a new Feed
func (s *Server) createFeedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		feedRequest := api.CreateFeedRequest{}
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

		feed, err := s.repo.GetFeed(vars["feedID"])
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

func (s *Server) getUserFeedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		feeds, err := s.repo.GetUserFeed(vars["userID"], vars["feedID"])
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
		addFeedRequest := api.AddUserFeedRequest{}
		if err := decodeAndValidate(req, &addFeedRequest); err != nil {
			s.formatter.Text(w, http.StatusBadRequest, err.Error())
			return
		}
		err := s.repo.AddUserFeed(vars["userID"], addFeedRequest.FeedID)
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}
		s.formatter.Text(w, http.StatusAccepted,
			fmt.Sprintf("Successfully subscribed User '%s' to Feed '%s'", vars["userID"], addFeedRequest.FeedID),
		)
	}
}
