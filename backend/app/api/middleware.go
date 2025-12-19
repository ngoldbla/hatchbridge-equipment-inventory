package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/hay-kot/httpkit/errchain"
	v1 "github.com/sysadminsmedia/homebox/backend/app/api/handlers/v1"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent"
	"github.com/sysadminsmedia/homebox/backend/internal/sys/validate"
)

type tokenHasKey struct {
	key string
}

var hashedToken = tokenHasKey{key: "hashedToken"}

type RoleMode int

const (
	RoleModeOr  RoleMode = 0
	RoleModeAnd RoleMode = 1
)

// mwRoles is a middleware that will validate the required roles are met. All roles
// are required to be met for the request to be allowed. If the user does not have
// the required roles, a 403 Forbidden will be returned.
//
// WARNING: This middleware _MUST_ be called after mwAuthToken or else it will panic
func (a *app) mwRoles(rm RoleMode, required ...string) errchain.Middleware {
	return func(next errchain.Handler) errchain.Handler {
		return errchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			maybeToken := ctx.Value(hashedToken)
			if maybeToken == nil {
				panic("mwRoles: token not found in context, you must call mwAuthToken before mwRoles")
			}

			token := maybeToken.(string)

			roles, err := a.repos.AuthTokens.GetRoles(r.Context(), token)
			if err != nil {
				return err
			}

		outer:
			switch rm {
			case RoleModeOr:
				for _, role := range required {
					if roles.Contains(role) {
						break outer
					}
				}
				return validate.NewRequestError(errors.New("Forbidden"), http.StatusForbidden)
			case RoleModeAnd:
				for _, req := range required {
					if !roles.Contains(req) {
						return validate.NewRequestError(errors.New("Unauthorized"), http.StatusForbidden)
					}
				}
			}

			return next.ServeHTTP(w, r)
		})
	}
}

type KeyFunc func(r *http.Request) (string, error)

func getBearer(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("authorization header is required")
	}

	return auth, nil
}

func getQuery(r *http.Request) (string, error) {
	token := r.URL.Query().Get("access_token")
	if token == "" {
		return "", errors.New("access_token query is required")
	}

	token, err := url.QueryUnescape(token)
	if err != nil {
		return "", errors.New("access_token query is required")
	}

	return token, nil
}

// mwAuthToken is a middleware that will check the database for a stateful token
// and attach it's user to the request context, or return an appropriate error.
// Authorization support is by token via Headers or Query Parameter
//
// Example:
//   - header = "Bearer 1234567890"
//   - query = "?access_token=1234567890"
func (a *app) mwAuthToken(next errchain.Handler) errchain.Handler {
	return errchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		var requestToken string

		// We ignore the error to allow the next strategy to be attempted
		{
			cookies, _ := v1.GetCookies(r)
			if cookies != nil {
				requestToken = cookies.Token
			}
		}

		if requestToken == "" {
			keyFuncs := [...]KeyFunc{
				getBearer,
				getQuery,
			}

			for _, keyFunc := range keyFuncs {
				token, err := keyFunc(r)
				if err == nil {
					requestToken = token
					break
				}
			}
		}

		if requestToken == "" {
			return validate.NewRequestError(errors.New("authorization header or query is required"), http.StatusUnauthorized)
		}

		requestToken = strings.TrimPrefix(requestToken, "Bearer ")

		r = r.WithContext(context.WithValue(r.Context(), hashedToken, requestToken))

		usr, err := a.services.User.GetSelf(r.Context(), requestToken)
		// Check the database for the token
		if err != nil {
			if ent.IsNotFound(err) {
				return validate.NewRequestError(errors.New("valid authorization token is required"), http.StatusUnauthorized)
			}

			return err
		}

		r = r.WithContext(services.SetUserCtx(r.Context(), &usr, requestToken))
		return next.ServeHTTP(w, r)
	})
}

// mwKioskContext is a middleware that checks if the user is in kiosk mode
// and sets the kiosk state in the context. This should be called after mwAuthToken.
func (a *app) mwKioskContext(next errchain.Handler) errchain.Handler {
	return errchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		user := services.UseUserCtx(r.Context())
		if user == nil {
			return next.ServeHTTP(w, r)
		}

		// Check if user has an active kiosk session
		session, err := a.repos.KioskSessions.GetByUserID(r.Context(), user.ID)
		if err != nil {
			// Don't fail the request, just proceed without kiosk mode
			return next.ServeHTTP(w, r)
		}

		if session != nil && session.IsActive {
			// User is in kiosk mode
			isUnlocked := session.IsUnlocked()
			r = r.WithContext(services.SetKioskCtx(r.Context(), true, isUnlocked))
		}

		return next.ServeHTTP(w, r)
	})
}

// mwKioskRestrict is a middleware that blocks destructive operations when in kiosk mode
// unless the kiosk is temporarily unlocked. This should be called after mwKioskContext.
func (a *app) mwKioskRestrict(next errchain.Handler) errchain.Handler {
	return errchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		isKiosk := services.UseKioskModeCtx(r.Context())
		isUnlocked := services.UseKioskUnlockedCtx(r.Context())

		// If in kiosk mode and not unlocked, block the request
		if isKiosk && !isUnlocked {
			return validate.NewRequestError(
				errors.New("this action requires admin access - please unlock to continue"),
				http.StatusForbidden,
			)
		}

		return next.ServeHTTP(w, r)
	})
}
