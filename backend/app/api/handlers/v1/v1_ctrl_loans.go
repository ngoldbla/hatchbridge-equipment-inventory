package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/data/repo"
	"github.com/sysadminsmedia/homebox/backend/internal/web/adapters"
)

// HandleLoansGetActive godoc
//
//	@Summary	Get Active Loans
//	@Tags		Loans
//	@Produce	json
//	@Success	200	{object}	[]repo.LoanSummary
//	@Router		/v1/loans [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleLoansGetActive() errchain.HandlerFunc {
	fn := func(r *http.Request) ([]repo.LoanSummary, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Loans.GetActiveLoans(auth, auth.GID)
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandleLoansGetOverdue godoc
//
//	@Summary	Get Overdue Loans
//	@Tags		Loans
//	@Produce	json
//	@Success	200	{object}	[]repo.LoanSummary
//	@Router		/v1/loans/overdue [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleLoansGetOverdue() errchain.HandlerFunc {
	fn := func(r *http.Request) ([]repo.LoanSummary, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Loans.GetOverdueLoans(auth, auth.GID)
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandleLoanCreate godoc
//
//	@Summary	Create Loan (Check Out Item)
//	@Tags		Loans
//	@Produce	json
//	@Param		payload	body		repo.LoanCreate	true	"Loan Data"
//	@Success	201		{object}	repo.LoanOut
//	@Router		/v1/loans [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleLoanCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, data repo.LoanCreate) (repo.LoanOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Loans.Create(auth, auth.GID, auth.UID, data)
	}

	return adapters.Action(fn, http.StatusCreated)
}

// HandleLoanGet godoc
//
//	@Summary	Get Loan
//	@Tags		Loans
//	@Produce	json
//	@Param		id	path		string	true	"Loan ID"
//	@Success	200	{object}	repo.LoanOut
//	@Router		/v1/loans/{id} [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleLoanGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) (repo.LoanOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Loans.GetOneByGroup(auth, auth.GID, ID)
	}

	return adapters.CommandID("id", fn, http.StatusOK)
}

// HandleLoanUpdate godoc
//
//	@Summary	Update Loan (Extend Due Date)
//	@Tags		Loans
//	@Produce	json
//	@Param		id		path		string			true	"Loan ID"
//	@Param		payload	body		repo.LoanUpdate	true	"Loan Data"
//	@Success	200		{object}	repo.LoanOut
//	@Router		/v1/loans/{id} [PUT]
//	@Security	Bearer
func (ctrl *V1Controller) HandleLoanUpdate() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, data repo.LoanUpdate) (repo.LoanOut, error) {
		auth := services.NewContext(r.Context())
		data.ID = ID
		return ctrl.repo.Loans.UpdateByGroup(auth, auth.GID, data)
	}

	return adapters.ActionID("id", fn, http.StatusOK)
}

// HandleLoanReturn godoc
//
//	@Summary	Return Loan (Check In Item)
//	@Tags		Loans
//	@Produce	json
//	@Param		id		path		string			true	"Loan ID"
//	@Param		payload	body		repo.LoanReturn	true	"Return Data"
//	@Success	200		{object}	repo.LoanOut
//	@Router		/v1/loans/{id}/return [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleLoanReturn() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, data repo.LoanReturn) (repo.LoanOut, error) {
		auth := services.NewContext(r.Context())
		data.ID = ID
		return ctrl.repo.Loans.Return(auth, auth.GID, auth.UID, data)
	}

	return adapters.ActionID("id", fn, http.StatusOK)
}

// HandleLoanDelete godoc
//
//	@Summary	Delete Loan
//	@Tags		Loans
//	@Produce	json
//	@Param		id	path	string	true	"Loan ID"
//	@Success	204
//	@Router		/v1/loans/{id} [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleLoanDelete() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) (any, error) {
		auth := services.NewContext(r.Context())
		err := ctrl.repo.Loans.DeleteByGroup(auth, auth.GID, ID)
		return nil, err
	}

	return adapters.CommandID("id", fn, http.StatusNoContent)
}

// HandleItemLoans godoc
//
//	@Summary	Get Item's Loan History
//	@Tags		Items
//	@Produce	json
//	@Param		id	path		string	true	"Item ID"
//	@Success	200	{object}	[]repo.LoanSummary
//	@Router		/v1/items/{id}/loans [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleItemLoans() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) ([]repo.LoanSummary, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Loans.GetLoansByItem(auth, auth.GID, ID)
	}

	return adapters.CommandID("id", fn, http.StatusOK)
}

// HandleItemCurrentLoan godoc
//
//	@Summary	Get Item's Current Loan
//	@Tags		Items
//	@Produce	json
//	@Param		id	path		string	true	"Item ID"
//	@Success	200	{object}	repo.LoanOut
//	@Router		/v1/items/{id}/current-loan [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleItemCurrentLoan() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) (*repo.LoanOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Loans.GetActiveLoanForItem(auth, auth.GID, ID)
	}

	return adapters.CommandID("id", fn, http.StatusOK)
}
