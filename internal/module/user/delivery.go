package user

import (
	"encoding/json"
	"fmt"
	"full-project-mock/internal/domain/usecase"
	"full-project-mock/internal/middleware"
	"full-project-mock/internal/service"
	"full-project-mock/pkg/req"
	"full-project-mock/pkg/respond"
	"net/http"
)

type UserHandler struct {
	Usecase usecase.UserUsecase
}

func (u *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	lgr := service.LoggerFromContext(r.Context())

	var registerUser RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerUser)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Invalid request payload", lgr)
		return
	}

	if err = registerUser.Validate(); err != nil {
		msg := fmt.Sprintf("User validation error: %v", err)
		respond.WithError(w, http.StatusBadRequest, msg, lgr)
		return
	}

	newUserId, err := u.Usecase.Register(r.Context(), registerUser.Email, registerUser.Username, registerUser.Password)
	if err != nil {
		msg := fmt.Sprintf("User registration error: %v", err)
		respond.WithError(w, http.StatusInternalServerError, msg, lgr)
		return
	}

	respond.WithSuccess(w, http.StatusCreated, "Успешно создано")
	lgr.Info(fmt.Sprintf("user registered %d", newUserId))
}

func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	lgr := service.LoggerFromContext(r.Context())

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Invalid request payload", lgr)
		return
	}

	ip, userAgent := req.GetClientMeta(r)
	accessToken, refreshToken, err := u.Usecase.Login(r.Context(), loginRequest.Email, loginRequest.Password, ip, userAgent)
	if err != nil {
		msg := fmt.Sprintf("User login error: %v", err)
		respond.WithError(w, http.StatusInternalServerError, msg, lgr)
		return
	}

	respond.WithSuccessJSON(w, http.StatusCreated, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (u *UserHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.GetUserIDFromContext(r.Context())
	if ok {
		json.NewEncoder(w).Encode(map[string]string{
			"id": id,
		})
		return
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Not found",
	})
}
