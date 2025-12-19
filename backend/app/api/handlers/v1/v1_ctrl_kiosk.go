package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/hay-kot/httpkit/errchain"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/web/adapters"
	"github.com/sysadminsmedia/homebox/backend/pkgs/hasher"
)

// KioskStatusResponse represents the current kiosk status for a user
type KioskStatusResponse struct {
	IsActive      bool       `json:"isActive"`
	IsUnlocked    bool       `json:"isUnlocked"`
	UnlockedUntil *time.Time `json:"unlockedUntil,omitempty"`
}

// KioskUnlockRequest represents the request to unlock kiosk mode
type KioskUnlockRequest struct {
	Password        string `json:"password" validate:"required"`
	DurationMinutes int    `json:"durationMinutes"`
}

// HandleKioskActivate godoc
//
//	@Summary	Activate Kiosk Mode
//	@Tags		Kiosk
//	@Produce	json
//	@Success	200	{object}	KioskStatusResponse
//	@Router		/v1/kiosk/activate [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleKioskActivate() errchain.HandlerFunc {
	fn := func(r *http.Request) (KioskStatusResponse, error) {
		auth := services.NewContext(r.Context())

		session, err := ctrl.repo.KioskSessions.Activate(auth, auth.UID)
		if err != nil {
			return KioskStatusResponse{}, err
		}

		return KioskStatusResponse{
			IsActive:      session.IsActive,
			IsUnlocked:    session.IsUnlocked(),
			UnlockedUntil: session.UnlockedUntil,
		}, nil
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandleKioskDeactivate godoc
//
//	@Summary	Deactivate Kiosk Mode
//	@Tags		Kiosk
//	@Produce	json
//	@Success	200	{object}	KioskStatusResponse
//	@Router		/v1/kiosk/deactivate [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleKioskDeactivate() errchain.HandlerFunc {
	fn := func(r *http.Request) (KioskStatusResponse, error) {
		auth := services.NewContext(r.Context())

		err := ctrl.repo.KioskSessions.Deactivate(auth, auth.UID)
		if err != nil {
			return KioskStatusResponse{}, err
		}

		return KioskStatusResponse{
			IsActive:   false,
			IsUnlocked: false,
		}, nil
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandleKioskStatus godoc
//
//	@Summary	Get Kiosk Status
//	@Tags		Kiosk
//	@Produce	json
//	@Success	200	{object}	KioskStatusResponse
//	@Router		/v1/kiosk/status [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleKioskStatus() errchain.HandlerFunc {
	fn := func(r *http.Request) (KioskStatusResponse, error) {
		auth := services.NewContext(r.Context())

		session, err := ctrl.repo.KioskSessions.GetByUserID(auth, auth.UID)
		if err != nil {
			return KioskStatusResponse{}, err
		}

		if session == nil {
			return KioskStatusResponse{
				IsActive:   false,
				IsUnlocked: false,
			}, nil
		}

		return KioskStatusResponse{
			IsActive:      session.IsActive,
			IsUnlocked:    session.IsUnlocked(),
			UnlockedUntil: session.UnlockedUntil,
		}, nil
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandleKioskUnlock godoc
//
//	@Summary	Unlock Kiosk Mode (Temporary Admin Access)
//	@Tags		Kiosk
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		KioskUnlockRequest	true	"Unlock Data"
//	@Success	200		{object}	KioskStatusResponse
//	@Router		/v1/kiosk/unlock [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleKioskUnlock() errchain.HandlerFunc {
	fn := func(r *http.Request, data KioskUnlockRequest) (KioskStatusResponse, error) {
		auth := services.NewContext(r.Context())

		// Get the user to verify password
		user, err := ctrl.repo.Users.GetOneID(auth, auth.UID)
		if err != nil {
			return KioskStatusResponse{}, err
		}

		// Verify password
		if user.PasswordHash == "" {
			return KioskStatusResponse{}, errors.New("password verification not available for this account")
		}

		passwordValid, _ := hasher.CheckPasswordHash(data.Password, user.PasswordHash)
		if !passwordValid {
			return KioskStatusResponse{}, errors.New("invalid password")
		}

		// Default to 5 minutes if not specified
		durationMinutes := data.DurationMinutes
		if durationMinutes <= 0 {
			durationMinutes = 5
		}
		// Cap at 30 minutes
		if durationMinutes > 30 {
			durationMinutes = 30
		}

		duration := time.Duration(durationMinutes) * time.Minute

		session, err := ctrl.repo.KioskSessions.Unlock(auth, auth.UID, duration)
		if err != nil {
			return KioskStatusResponse{}, err
		}

		if session == nil {
			return KioskStatusResponse{}, errors.New("no active kiosk session to unlock")
		}

		return KioskStatusResponse{
			IsActive:      session.IsActive,
			IsUnlocked:    session.IsUnlocked(),
			UnlockedUntil: session.UnlockedUntil,
		}, nil
	}

	return adapters.Action(fn, http.StatusOK)
}

// HandleKioskLock godoc
//
//	@Summary	Lock Kiosk Mode (Revoke Temporary Admin Access)
//	@Tags		Kiosk
//	@Produce	json
//	@Success	200	{object}	KioskStatusResponse
//	@Router		/v1/kiosk/lock [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleKioskLock() errchain.HandlerFunc {
	fn := func(r *http.Request) (KioskStatusResponse, error) {
		auth := services.NewContext(r.Context())

		err := ctrl.repo.KioskSessions.Lock(auth, auth.UID)
		if err != nil {
			return KioskStatusResponse{}, err
		}

		session, err := ctrl.repo.KioskSessions.GetByUserID(auth, auth.UID)
		if err != nil {
			return KioskStatusResponse{}, err
		}

		if session == nil {
			return KioskStatusResponse{
				IsActive:   false,
				IsUnlocked: false,
			}, nil
		}

		return KioskStatusResponse{
			IsActive:      session.IsActive,
			IsUnlocked:    false,
			UnlockedUntil: nil,
		}, nil
	}

	return adapters.Command(fn, http.StatusOK)
}
