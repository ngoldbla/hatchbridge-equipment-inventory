package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/data/repo"
	"github.com/sysadminsmedia/homebox/backend/internal/web/adapters"
)

// HandleBorrowersGetAll godoc
//
//	@Summary	Get All Borrowers
//	@Tags		Borrowers
//	@Produce	json
//	@Success	200	{object}	[]repo.BorrowerSummary
//	@Router		/v1/borrowers [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBorrowersGetAll() errchain.HandlerFunc {
	fn := func(r *http.Request) ([]repo.BorrowerSummary, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Borrowers.GetAll(auth, auth.GID)
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandleBorrowersGetActive godoc
//
//	@Summary	Get Active Borrowers
//	@Tags		Borrowers
//	@Produce	json
//	@Success	200	{object}	[]repo.BorrowerSummary
//	@Router		/v1/borrowers/active [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBorrowersGetActive() errchain.HandlerFunc {
	fn := func(r *http.Request) ([]repo.BorrowerSummary, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Borrowers.GetActive(auth, auth.GID)
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandleBorrowersCreate godoc
//
//	@Summary	Create Borrower
//	@Tags		Borrowers
//	@Produce	json
//	@Param		payload	body		repo.BorrowerCreate	true	"Borrower Data"
//	@Success	201		{object}	repo.BorrowerOut
//	@Router		/v1/borrowers [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBorrowersCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, data repo.BorrowerCreate) (repo.BorrowerOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Borrowers.Create(auth, auth.GID, data)
	}

	return adapters.Action(fn, http.StatusCreated)
}

// HandleBorrowerGet godoc
//
//	@Summary	Get Borrower
//	@Tags		Borrowers
//	@Produce	json
//	@Param		id	path		string	true	"Borrower ID"
//	@Success	200	{object}	repo.BorrowerOut
//	@Router		/v1/borrowers/{id} [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBorrowerGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) (repo.BorrowerOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Borrowers.GetOneByGroup(auth, auth.GID, ID)
	}

	return adapters.CommandID("id", fn, http.StatusOK)
}

// HandleBorrowerUpdate godoc
//
//	@Summary	Update Borrower
//	@Tags		Borrowers
//	@Produce	json
//	@Param		id		path		string				true	"Borrower ID"
//	@Param		payload	body		repo.BorrowerUpdate	true	"Borrower Data"
//	@Success	200		{object}	repo.BorrowerOut
//	@Router		/v1/borrowers/{id} [PUT]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBorrowerUpdate() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, data repo.BorrowerUpdate) (repo.BorrowerOut, error) {
		auth := services.NewContext(r.Context())
		data.ID = ID
		return ctrl.repo.Borrowers.UpdateByGroup(auth, auth.GID, data)
	}

	return adapters.ActionID("id", fn, http.StatusOK)
}

// HandleBorrowerDelete godoc
//
//	@Summary	Delete Borrower
//	@Tags		Borrowers
//	@Produce	json
//	@Param		id	path	string	true	"Borrower ID"
//	@Success	204
//	@Router		/v1/borrowers/{id} [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBorrowerDelete() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) (any, error) {
		auth := services.NewContext(r.Context())
		err := ctrl.repo.Borrowers.DeleteByGroup(auth, auth.GID, ID)
		return nil, err
	}

	return adapters.CommandID("id", fn, http.StatusNoContent)
}

// HandleBorrowerLoans godoc
//
//	@Summary	Get Borrower's Loans
//	@Tags		Borrowers
//	@Produce	json
//	@Param		id	path		string	true	"Borrower ID"
//	@Success	200	{object}	[]repo.LoanSummary
//	@Router		/v1/borrowers/{id}/loans [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBorrowerLoans() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) ([]repo.LoanSummary, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Loans.GetLoansByBorrower(auth, auth.GID, ID)
	}

	return adapters.CommandID("id", fn, http.StatusOK)
}
