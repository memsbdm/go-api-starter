package middleware

import (
	"go-starter/internal/adapters/server/helpers"
	"go-starter/internal/adapters/server/responses"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"net/http"
	"slices"
)

// RoleMiddleware is a middleware function that checks if the user has the required role to access the resource.
func RoleMiddleware(userSvc ports.UserService, authMiddleware Middleware, roleIDs ...entities.RoleID) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		// apply auth middleware first to ensure we have a valid user id in the context
		handlerWithAuth := authMiddleware(func(w http.ResponseWriter, r *http.Request) {
			userID, err := helpers.GetUserIDFromContext(r.Context())
			if err != nil {
				responses.HandleError(w, domain.ErrUnauthorized)
				return
			}

			user, err := userSvc.GetByID(r.Context(), userID)
			if err != nil {
				responses.HandleError(w, domain.ErrUnauthorized)
				return
			}

			hasValidRole := slices.Contains(roleIDs, user.RoleID)

			if !hasValidRole {
				responses.HandleError(w, domain.ErrForbidden)
				return
			}

			f(w, r)
		})

		return handlerWithAuth
	}
}
