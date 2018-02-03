package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

// createUserRequest represents a request to create a new User
type createUserRequest struct {
	Name string `json:"name" valid:"required~User name cannot be blank,alphanum~User name should be alphanumeric"`
}

func (s *Server) createUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		userRequest := createUserRequest{}
		if err := decodeAndValidate(req, &userRequest); err != nil {
			s.formatter.Text(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := s.repo.CreateUser(userRequest.Name)
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusCreated, user)
	}
}

func (s *Server) getUserListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		users, err := s.repo.ListUsers()
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusOK, users)
	}
}

func (s *Server) getUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		user, err := s.repo.GetUser(vars["userID"])
		if err != nil {
			s.formatter.Text(w, errorToStatus(err), err.Error())
			return
		}

		s.formatter.JSON(w, http.StatusOK, user)
	}
}
