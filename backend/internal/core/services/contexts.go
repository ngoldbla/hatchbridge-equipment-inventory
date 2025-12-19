package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/repo"
)

type contextKeys struct {
	name string
}

var (
	ContextUser          = &contextKeys{name: "User"}
	ContextUserToken     = &contextKeys{name: "UserToken"}
	ContextKioskMode     = &contextKeys{name: "KioskMode"}
	ContextKioskUnlocked = &contextKeys{name: "KioskUnlocked"}
)

type Context struct {
	context.Context

	// UID is a unique identifier for the acting user.
	UID uuid.UUID

	// GID is a unique identifier for the acting users group.
	GID uuid.UUID

	// User is the acting user.
	User *repo.UserOut

	// IsKiosk indicates whether the session is in kiosk mode.
	IsKiosk bool

	// IsKioskUnlocked indicates whether the kiosk has temporary admin access.
	IsKioskUnlocked bool
}

// NewContext is a helper function that returns the service context from the context.
// This extracts the users from the context and embeds it into the ServiceContext struct
func NewContext(ctx context.Context) Context {
	user := UseUserCtx(ctx)
	return Context{
		Context:         ctx,
		UID:             user.ID,
		GID:             user.GroupID,
		User:            user,
		IsKiosk:         UseKioskModeCtx(ctx),
		IsKioskUnlocked: UseKioskUnlockedCtx(ctx),
	}
}

// SetUserCtx is a helper function that sets the ContextUser and ContextUserToken
// values within the context of a web request (or any context).
func SetUserCtx(ctx context.Context, user *repo.UserOut, token string) context.Context {
	ctx = context.WithValue(ctx, ContextUser, user)
	ctx = context.WithValue(ctx, ContextUserToken, token)
	return ctx
}

// SetKioskCtx sets the kiosk mode state in the context.
func SetKioskCtx(ctx context.Context, isKiosk bool, isUnlocked bool) context.Context {
	ctx = context.WithValue(ctx, ContextKioskMode, isKiosk)
	ctx = context.WithValue(ctx, ContextKioskUnlocked, isUnlocked)
	return ctx
}

// UseUserCtx is a helper function that returns the user from the context.
func UseUserCtx(ctx context.Context) *repo.UserOut {
	if val := ctx.Value(ContextUser); val != nil {
		return val.(*repo.UserOut)
	}
	return nil
}

// UseTokenCtx is a helper function that returns the user token from the context.
func UseTokenCtx(ctx context.Context) string {
	if val := ctx.Value(ContextUserToken); val != nil {
		return val.(string)
	}
	return ""
}

// UseKioskModeCtx returns whether the session is in kiosk mode.
func UseKioskModeCtx(ctx context.Context) bool {
	if val := ctx.Value(ContextKioskMode); val != nil {
		return val.(bool)
	}
	return false
}

// UseKioskUnlockedCtx returns whether the kiosk has temporary admin access.
func UseKioskUnlockedCtx(ctx context.Context) bool {
	if val := ctx.Value(ContextKioskUnlocked); val != nil {
		return val.(bool)
	}
	return false
}
