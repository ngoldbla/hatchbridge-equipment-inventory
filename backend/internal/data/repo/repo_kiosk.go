package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/kiosksession"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/user"
)

// KioskSessionRepository is a repository for managing kiosk sessions.
type KioskSessionRepository struct {
	db *ent.Client
}

// KioskSessionOut represents the kiosk session output data
type KioskSessionOut struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"userId"`
	IsActive      bool       `json:"isActive"`
	UnlockedUntil *time.Time `json:"unlockedUntil,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// IsUnlocked returns true if the kiosk session has temporary admin access
func (s *KioskSessionOut) IsUnlocked() bool {
	if s.UnlockedUntil == nil {
		return false
	}
	return time.Now().Before(*s.UnlockedUntil)
}

func mapKioskSessionOut(session *ent.KioskSession) KioskSessionOut {
	return KioskSessionOut{
		ID:            session.ID,
		UserID:        session.Edges.User.ID,
		IsActive:      session.IsActive,
		UnlockedUntil: session.UnlockedUntil,
		CreatedAt:     session.CreatedAt,
		UpdatedAt:     session.UpdatedAt,
	}
}

// GetByUserID gets the kiosk session for a user
func (r *KioskSessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*KioskSessionOut, error) {
	session, err := r.db.KioskSession.
		Query().
		Where(kiosksession.HasUserWith(user.ID(userID))).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	out := mapKioskSessionOut(session)
	return &out, nil
}

// Activate creates or updates a kiosk session for a user
func (r *KioskSessionRepository) Activate(ctx context.Context, userID uuid.UUID) (*KioskSessionOut, error) {
	// Check if session already exists
	existing, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Update existing session to be active
		session, err := r.db.KioskSession.
			UpdateOneID(existing.ID).
			SetIsActive(true).
			SetUpdatedAt(time.Now()).
			ClearUnlockedUntil().
			Save(ctx)
		if err != nil {
			return nil, err
		}

		// Reload with edges
		session, err = r.db.KioskSession.
			Query().
			Where(kiosksession.ID(session.ID)).
			WithUser().
			Only(ctx)
		if err != nil {
			return nil, err
		}

		out := mapKioskSessionOut(session)
		return &out, nil
	}

	// Create new session
	session, err := r.db.KioskSession.
		Create().
		SetUserID(userID).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with edges
	session, err = r.db.KioskSession.
		Query().
		Where(kiosksession.ID(session.ID)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	out := mapKioskSessionOut(session)
	return &out, nil
}

// Deactivate sets the kiosk session to inactive
func (r *KioskSessionRepository) Deactivate(ctx context.Context, userID uuid.UUID) error {
	existing, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if existing == nil {
		return nil // No session to deactivate
	}

	_, err = r.db.KioskSession.
		UpdateOneID(existing.ID).
		SetIsActive(false).
		SetUpdatedAt(time.Now()).
		ClearUnlockedUntil().
		Save(ctx)

	return err
}

// Unlock temporarily unlocks the kiosk session for admin access
func (r *KioskSessionRepository) Unlock(ctx context.Context, userID uuid.UUID, duration time.Duration) (*KioskSessionOut, error) {
	existing, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if existing == nil || !existing.IsActive {
		return nil, nil // No active session to unlock
	}

	unlockUntil := time.Now().Add(duration)

	session, err := r.db.KioskSession.
		UpdateOneID(existing.ID).
		SetUnlockedUntil(unlockUntil).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with edges
	session, err = r.db.KioskSession.
		Query().
		Where(kiosksession.ID(session.ID)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	out := mapKioskSessionOut(session)
	return &out, nil
}

// Lock explicitly locks the kiosk session (clears unlock timer)
func (r *KioskSessionRepository) Lock(ctx context.Context, userID uuid.UUID) error {
	existing, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if existing == nil {
		return nil
	}

	_, err = r.db.KioskSession.
		UpdateOneID(existing.ID).
		ClearUnlockedUntil().
		SetUpdatedAt(time.Now()).
		Save(ctx)

	return err
}
