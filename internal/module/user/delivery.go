package user

import (
	"encoding/json"
	"fmt"
	"github.com/Elaman1/full-project-mock/internal/domain/usecase"
	"github.com/Elaman1/full-project-mock/internal/metrics"
	"github.com/Elaman1/full-project-mock/internal/middleware"
	"github.com/Elaman1/full-project-mock/internal/service"
	"github.com/Elaman1/full-project-mock/pkg/req"
	"github.com/Elaman1/full-project-mock/pkg/respond"
	"net/http"
)

type UserHandler struct {
	Usecase         usecase.UserUsecase
	MetricCollector metrics.MetricsCollector
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

	respond.WithSuccess(w, http.StatusCreated, "user registered")
	lgr.Info(fmt.Sprintf("user registered %d", newUserId))
}

func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	lgr := service.LoggerFromContext(r.Context())

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Invalid request payload", lgr)
		u.MetricCollector.LoginFailureCounter("Invalid request payload")
		return
	}

	err = loginRequest.Validate()
	if err != nil {
		msg := fmt.Sprintf("User validation error: %v", err)
		respond.WithError(w, http.StatusBadRequest, msg, lgr)
		u.MetricCollector.LoginFailureCounter("User validation error")
		return
	}

	ip, userAgent := req.GetClientMeta(r)
	accessToken, refreshToken, httpStatus, err := u.Usecase.Login(r.Context(), loginRequest.Email, loginRequest.Password, ip, userAgent)
	if err != nil {
		msg := fmt.Sprintf("User login error: %v", err)
		respond.WithError(w, httpStatus, msg, lgr)
		u.MetricCollector.LoginFailureCounter("User login error")
		return
	}

	u.MetricCollector.LoginSuccessCounter()
	respond.WithSuccessJSON(w, httpStatus, map[string]string{
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

func (u *UserHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	lgr := service.LoggerFromContext(r.Context())

	var refreshRequest RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&refreshRequest)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Invalid request payload", lgr)
		return
	}

	ip, userAgent := req.GetClientMeta(r)
	accessToken, refreshToken, httpStatus, err := u.Usecase.Refresh(r.Context(), refreshRequest.AccessToken, refreshRequest.RefreshToken, ip, userAgent)
	if err != nil {
		msg := fmt.Sprintf("Refresh error: %v", err)
		respond.WithError(w, httpStatus, msg, lgr)
		return
	}

	respond.WithSuccessJSON(w, httpStatus, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (u *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	lgr := service.LoggerFromContext(r.Context())

	var refreshRequest RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&refreshRequest)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Invalid request payload", lgr)
		return
	}

	ip, userAgent := req.GetClientMeta(r)
	err = u.Usecase.Logout(r.Context(), refreshRequest.RefreshToken, ip, userAgent)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "error logout", lgr)
		return
	}

	respond.WithSuccessJSON(w, http.StatusCreated, map[string]string{
		"message": "Успешно вышли из аккаунта",
	})
}

func (u *UserHandler) LogoutAllHandler(w http.ResponseWriter, r *http.Request) {
	lgr := service.LoggerFromContext(r.Context())

	var refreshRequest RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&refreshRequest)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "Invalid request payload", lgr)
		return
	}

	ip, userAgent := req.GetClientMeta(r)
	err = u.Usecase.LogoutAllDevices(r.Context(), refreshRequest.RefreshToken, ip, userAgent)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "error logout all", lgr)
		return
	}
}
