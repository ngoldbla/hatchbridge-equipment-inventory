package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services/reporting/eventbus"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/borrower"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/group"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/predicate"
)

type BorrowerRepository struct {
	db  *ent.Client
	bus *eventbus.EventBus
}

type (
	BorrowerCreate struct {
		Name         string `json:"name"         validate:"required,min=1,max=255"`
		Email        string `json:"email"        validate:"required,email,max=255"`
		Phone        string `json:"phone"        validate:"max=50"`
		Organization string `json:"organization" validate:"max=255"`
		StudentID    string `json:"studentId"    validate:"max=100"`
		Notes        string `json:"notes"        validate:"max=1000"`
	}

	BorrowerUpdate struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"         validate:"required,min=1,max=255"`
		Email        string    `json:"email"        validate:"required,email,max=255"`
		Phone        string    `json:"phone"        validate:"max=50"`
		Organization string    `json:"organization" validate:"max=255"`
		StudentID    string    `json:"studentId"    validate:"max=100"`
		Notes        string    `json:"notes"        validate:"max=1000"`
		IsActive     bool      `json:"isActive"`
	}

	BorrowerSummary struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		Email        string    `json:"email"`
		Phone        string    `json:"phone"`
		Organization string    `json:"organization"`
		StudentID    string    `json:"studentId"`
		IsActive     bool      `json:"isActive"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	BorrowerOut struct {
		BorrowerSummary
		Notes       string `json:"notes"`
		ActiveLoans int    `json:"activeLoans"`
		TotalLoans  int    `json:"totalLoans"`
	}
)

func mapBorrowerSummary(b *ent.Borrower) BorrowerSummary {
	return BorrowerSummary{
		ID:           b.ID,
		Name:         b.Name,
		Email:        b.Email,
		Phone:        b.Phone,
		Organization: b.Organization,
		StudentID:    b.StudentID,
		IsActive:     b.IsActive,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}

var (
	mapBorrowerOutErr   = mapTErrFunc(mapBorrowerOut)
	mapBorrowersSummary = mapTEachErrFunc(mapBorrowerSummary)
)

func mapBorrowerOut(b *ent.Borrower) BorrowerOut {
	return BorrowerOut{
		BorrowerSummary: mapBorrowerSummary(b),
		Notes:           b.Notes,
		ActiveLoans:     0, // Will be populated by query if needed
		TotalLoans:      0, // Will be populated by query if needed
	}
}

func (r *BorrowerRepository) publishMutationEvent(gid uuid.UUID) {
	if r.bus != nil {
		r.bus.Publish(eventbus.EventBorrowerMutation, eventbus.GroupMutationEvent{GID: gid})
	}
}

func (r *BorrowerRepository) getOne(ctx context.Context, where ...predicate.Borrower) (BorrowerOut, error) {
	return mapBorrowerOutErr(r.db.Borrower.Query().
		Where(where...).
		WithGroup().
		Only(ctx),
	)
}

func (r *BorrowerRepository) GetOne(ctx context.Context, id uuid.UUID) (BorrowerOut, error) {
	return r.getOne(ctx, borrower.ID(id))
}

func (r *BorrowerRepository) GetOneByGroup(ctx context.Context, gid, bid uuid.UUID) (BorrowerOut, error) {
	return r.getOne(ctx, borrower.ID(bid), borrower.HasGroupWith(group.ID(gid)))
}

func (r *BorrowerRepository) GetAll(ctx context.Context, groupID uuid.UUID) ([]BorrowerSummary, error) {
	return mapBorrowersSummary(r.db.Borrower.Query().
		Where(borrower.HasGroupWith(group.ID(groupID))).
		Order(ent.Asc(borrower.FieldName)).
		WithGroup().
		All(ctx),
	)
}

func (r *BorrowerRepository) GetActive(ctx context.Context, groupID uuid.UUID) ([]BorrowerSummary, error) {
	return mapBorrowersSummary(r.db.Borrower.Query().
		Where(
			borrower.HasGroupWith(group.ID(groupID)),
			borrower.IsActive(true),
		).
		Order(ent.Asc(borrower.FieldName)).
		WithGroup().
		All(ctx),
	)
}

func (r *BorrowerRepository) Create(ctx context.Context, groupID uuid.UUID, data BorrowerCreate) (BorrowerOut, error) {
	b, err := r.db.Borrower.Create().
		SetName(data.Name).
		SetEmail(data.Email).
		SetPhone(data.Phone).
		SetOrganization(data.Organization).
		SetStudentID(data.StudentID).
		SetNotes(data.Notes).
		SetIsActive(true).
		SetGroupID(groupID).
		Save(ctx)
	if err != nil {
		return BorrowerOut{}, err
	}

	b.Edges.Group = &ent.Group{ID: groupID}
	r.publishMutationEvent(groupID)
	return mapBorrowerOut(b), nil
}

func (r *BorrowerRepository) update(ctx context.Context, data BorrowerUpdate, where ...predicate.Borrower) (int, error) {
	if len(where) == 0 {
		panic("empty where not supported")
	}

	return r.db.Borrower.Update().
		Where(where...).
		SetName(data.Name).
		SetEmail(data.Email).
		SetPhone(data.Phone).
		SetOrganization(data.Organization).
		SetStudentID(data.StudentID).
		SetNotes(data.Notes).
		SetIsActive(data.IsActive).
		Save(ctx)
}

func (r *BorrowerRepository) UpdateByGroup(ctx context.Context, gid uuid.UUID, data BorrowerUpdate) (BorrowerOut, error) {
	_, err := r.update(ctx, data, borrower.ID(data.ID), borrower.HasGroupWith(group.ID(gid)))
	if err != nil {
		return BorrowerOut{}, err
	}

	r.publishMutationEvent(gid)
	return r.GetOne(ctx, data.ID)
}

func (r *BorrowerRepository) DeleteByGroup(ctx context.Context, gid, id uuid.UUID) error {
	_, err := r.db.Borrower.Delete().
		Where(
			borrower.ID(id),
			borrower.HasGroupWith(group.ID(gid)),
		).Exec(ctx)
	if err != nil {
		return err
	}

	r.publishMutationEvent(gid)
	return nil
}

// SetActive sets the active status of a borrower
func (r *BorrowerRepository) SetActive(ctx context.Context, gid, id uuid.UUID, active bool) error {
	_, err := r.db.Borrower.Update().
		Where(
			borrower.ID(id),
			borrower.HasGroupWith(group.ID(gid)),
		).
		SetIsActive(active).
		Save(ctx)
	if err != nil {
		return err
	}

	r.publishMutationEvent(gid)
	return nil
}
