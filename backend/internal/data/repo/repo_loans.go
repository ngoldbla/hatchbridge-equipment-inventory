package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services/reporting/eventbus"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/borrower"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/group"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/item"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/loan"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/predicate"
)

type LoanRepository struct {
	db  *ent.Client
	bus *eventbus.EventBus
}

type (
	LoanCreate struct {
		ItemID     uuid.UUID `json:"itemId"     validate:"required"`
		BorrowerID uuid.UUID `json:"borrowerId" validate:"required"`
		DueAt      time.Time `json:"dueAt"      validate:"required"`
		Notes      string    `json:"notes"      validate:"max=1000"`
		Quantity   int       `json:"quantity"   validate:"min=1"`
	}

	LoanUpdate struct {
		ID    uuid.UUID `json:"id"`
		DueAt time.Time `json:"dueAt"`
		Notes string    `json:"notes" validate:"max=1000"`
	}

	LoanReturn struct {
		ID          uuid.UUID `json:"id"`
		ReturnNotes string    `json:"returnNotes" validate:"max=1000"`
	}

	LoanSummary struct {
		ID           uuid.UUID  `json:"id"`
		CheckedOutAt time.Time  `json:"checkedOutAt"`
		DueAt        time.Time  `json:"dueAt"`
		ReturnedAt   *time.Time `json:"returnedAt"`
		Quantity     int        `json:"quantity"`
		IsOverdue    bool       `json:"isOverdue"`
		ItemID       uuid.UUID  `json:"itemId"`
		ItemName     string     `json:"itemName"`
		BorrowerID   uuid.UUID  `json:"borrowerId"`
		BorrowerName string     `json:"borrowerName"`
		CreatedAt    time.Time  `json:"createdAt"`
		UpdatedAt    time.Time  `json:"updatedAt"`
	}

	LoanOut struct {
		LoanSummary
		Notes         string     `json:"notes"`
		ReturnNotes   string     `json:"returnNotes"`
		ItemAssetID   int        `json:"itemAssetId"`
		BorrowerEmail string     `json:"borrowerEmail"`
		BorrowerPhone string     `json:"borrowerPhone"`
		CheckedOutBy  *uuid.UUID `json:"checkedOutBy"`
		ReturnedBy    *uuid.UUID `json:"returnedBy"`
	}
)

func mapLoanSummary(l *ent.Loan) LoanSummary {
	summary := LoanSummary{
		ID:           l.ID,
		CheckedOutAt: l.CheckedOutAt,
		DueAt:        l.DueAt,
		ReturnedAt:   l.ReturnedAt,
		Quantity:     l.Quantity,
		IsOverdue:    l.ReturnedAt == nil && time.Now().After(l.DueAt),
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}

	if l.Edges.Item != nil {
		summary.ItemID = l.Edges.Item.ID
		summary.ItemName = l.Edges.Item.Name
	}

	if l.Edges.Borrower != nil {
		summary.BorrowerID = l.Edges.Borrower.ID
		summary.BorrowerName = l.Edges.Borrower.Name
	}

	return summary
}

var (
	mapLoanOutErr   = mapTErrFunc(mapLoanOut)
	mapLoansSummary = mapTEachErrFunc(mapLoanSummary)
)

func mapLoanOut(l *ent.Loan) LoanOut {
	out := LoanOut{
		LoanSummary: mapLoanSummary(l),
		Notes:       l.Notes,
		ReturnNotes: l.ReturnNotes,
	}

	if l.Edges.Item != nil {
		out.ItemAssetID = l.Edges.Item.AssetID
	}

	if l.Edges.Borrower != nil {
		out.BorrowerEmail = l.Edges.Borrower.Email
		out.BorrowerPhone = l.Edges.Borrower.Phone
	}

	if l.Edges.CheckedOutBy != nil {
		out.CheckedOutBy = &l.Edges.CheckedOutBy.ID
	}

	if l.Edges.ReturnedBy != nil {
		out.ReturnedBy = &l.Edges.ReturnedBy.ID
	}

	return out
}

func (r *LoanRepository) publishMutationEvent(gid uuid.UUID) {
	if r.bus != nil {
		r.bus.Publish(eventbus.EventLoanMutation, eventbus.GroupMutationEvent{GID: gid})
	}
}

func (r *LoanRepository) getOne(ctx context.Context, where ...predicate.Loan) (LoanOut, error) {
	return mapLoanOutErr(r.db.Loan.Query().
		Where(where...).
		WithGroup().
		WithItem().
		WithBorrower().
		WithCheckedOutBy().
		WithReturnedBy().
		Only(ctx),
	)
}

func (r *LoanRepository) GetOne(ctx context.Context, id uuid.UUID) (LoanOut, error) {
	return r.getOne(ctx, loan.ID(id))
}

func (r *LoanRepository) GetOneByGroup(ctx context.Context, gid, lid uuid.UUID) (LoanOut, error) {
	return r.getOne(ctx, loan.ID(lid), loan.HasGroupWith(group.ID(gid)))
}

// GetActiveLoans returns all loans that have not been returned
func (r *LoanRepository) GetActiveLoans(ctx context.Context, groupID uuid.UUID) ([]LoanSummary, error) {
	return mapLoansSummary(r.db.Loan.Query().
		Where(
			loan.HasGroupWith(group.ID(groupID)),
			loan.ReturnedAtIsNil(),
		).
		Order(ent.Asc(loan.FieldDueAt)).
		WithItem().
		WithBorrower().
		All(ctx),
	)
}

// GetOverdueLoans returns all active loans that are past due
func (r *LoanRepository) GetOverdueLoans(ctx context.Context, groupID uuid.UUID) ([]LoanSummary, error) {
	return mapLoansSummary(r.db.Loan.Query().
		Where(
			loan.HasGroupWith(group.ID(groupID)),
			loan.ReturnedAtIsNil(),
			loan.DueAtLT(time.Now()),
		).
		Order(ent.Asc(loan.FieldDueAt)).
		WithItem().
		WithBorrower().
		All(ctx),
	)
}

// GetLoansByBorrower returns all loans for a specific borrower
func (r *LoanRepository) GetLoansByBorrower(ctx context.Context, gid, borrowerID uuid.UUID) ([]LoanSummary, error) {
	return mapLoansSummary(r.db.Loan.Query().
		Where(
			loan.HasGroupWith(group.ID(gid)),
			loan.HasBorrowerWith(borrower.ID(borrowerID)),
		).
		Order(ent.Desc(loan.FieldCheckedOutAt)).
		WithItem().
		WithBorrower().
		All(ctx),
	)
}

// GetLoansByItem returns all loans for a specific item
func (r *LoanRepository) GetLoansByItem(ctx context.Context, gid, itemID uuid.UUID) ([]LoanSummary, error) {
	return mapLoansSummary(r.db.Loan.Query().
		Where(
			loan.HasGroupWith(group.ID(gid)),
			loan.HasItemWith(item.ID(itemID)),
		).
		Order(ent.Desc(loan.FieldCheckedOutAt)).
		WithItem().
		WithBorrower().
		All(ctx),
	)
}

// GetActiveLoanForItem returns the active loan for an item if one exists
func (r *LoanRepository) GetActiveLoanForItem(ctx context.Context, gid, itemID uuid.UUID) (*LoanOut, error) {
	l, err := r.db.Loan.Query().
		Where(
			loan.HasGroupWith(group.ID(gid)),
			loan.HasItemWith(item.ID(itemID)),
			loan.ReturnedAtIsNil(),
		).
		WithItem().
		WithBorrower().
		WithCheckedOutBy().
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	out := mapLoanOut(l)
	return &out, nil
}

// Create creates a new loan (checkout)
func (r *LoanRepository) Create(ctx context.Context, gid uuid.UUID, userID uuid.UUID, data LoanCreate) (LoanOut, error) {
	quantity := data.Quantity
	if quantity == 0 {
		quantity = 1
	}

	l, err := r.db.Loan.Create().
		SetItemID(data.ItemID).
		SetBorrowerID(data.BorrowerID).
		SetGroupID(gid).
		SetCheckedOutAt(time.Now()).
		SetDueAt(data.DueAt).
		SetNotes(data.Notes).
		SetQuantity(quantity).
		SetCheckedOutByID(userID).
		Save(ctx)
	if err != nil {
		return LoanOut{}, err
	}

	r.publishMutationEvent(gid)
	return r.GetOne(ctx, l.ID)
}

// Return marks a loan as returned
func (r *LoanRepository) Return(ctx context.Context, gid uuid.UUID, userID uuid.UUID, data LoanReturn) (LoanOut, error) {
	now := time.Now()
	_, err := r.db.Loan.Update().
		Where(
			loan.ID(data.ID),
			loan.HasGroupWith(group.ID(gid)),
			loan.ReturnedAtIsNil(), // Ensure not already returned
		).
		SetReturnedAt(now).
		SetReturnNotes(data.ReturnNotes).
		SetReturnedByID(userID).
		Save(ctx)
	if err != nil {
		return LoanOut{}, err
	}

	r.publishMutationEvent(gid)
	return r.GetOne(ctx, data.ID)
}

// UpdateByGroup updates a loan's details (e.g., extend due date)
func (r *LoanRepository) UpdateByGroup(ctx context.Context, gid uuid.UUID, data LoanUpdate) (LoanOut, error) {
	_, err := r.db.Loan.Update().
		Where(
			loan.ID(data.ID),
			loan.HasGroupWith(group.ID(gid)),
		).
		SetDueAt(data.DueAt).
		SetNotes(data.Notes).
		Save(ctx)
	if err != nil {
		return LoanOut{}, err
	}

	r.publishMutationEvent(gid)
	return r.GetOne(ctx, data.ID)
}

// DeleteByGroup deletes a loan record (use with caution)
func (r *LoanRepository) DeleteByGroup(ctx context.Context, gid, id uuid.UUID) error {
	_, err := r.db.Loan.Delete().
		Where(
			loan.ID(id),
			loan.HasGroupWith(group.ID(gid)),
		).Exec(ctx)
	if err != nil {
		return err
	}

	r.publishMutationEvent(gid)
	return nil
}
