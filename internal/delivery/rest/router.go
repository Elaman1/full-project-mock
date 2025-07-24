package rest

import (
	"context"
	"github.com/Elaman1/full-project-mock/internal/domain/usecase"
	"github.com/Elaman1/full-project-mock/internal/middleware"
	"github.com/Elaman1/full-project-mock/internal/module"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func InitRouter(ctx context.Context, routeApp *RouteApp, allModules *module.Modules) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.LogMiddleware(routeApp.Logs))
	r.Use(middleware.ContextJoinMiddleware(ctx))

	r.Post("/register", allModules.UserHandler.RegisterHandler)
	r.Post("/login", allModules.UserHandler.LoginHandler)
	r.Post("/refresh", allModules.UserHandler.RefreshHandler)

	// auth group
	r.Route("/auth", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(routeApp.TokenService))

		r.Get("/me", allModules.UserHandler.MeHandler)
		r.Post("/logout", allModules.UserHandler.LogoutHandler)
		r.Post("/logout-all", allModules.UserHandler.LogoutAllHandler)
	})

	return r
}

type RouteApp struct {
	Logs         *slog.Logger
	TokenService usecase.TokenService
}
