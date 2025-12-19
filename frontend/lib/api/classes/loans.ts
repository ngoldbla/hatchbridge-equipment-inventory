import { BaseAPI, route } from "../base";
import type { LoanCreate, LoanOut, LoanReturn, LoanSummary, LoanUpdate } from "../types/data-contracts";

export class LoansApi extends BaseAPI {
  getActive() {
    return this.http.get<LoanSummary[]>({ url: route("/loans") });
  }

  getOverdue() {
    return this.http.get<LoanSummary[]>({ url: route("/loans/overdue") });
  }

  create(body: LoanCreate) {
    return this.http.post<LoanCreate, LoanOut>({ url: route("/loans"), body });
  }

  get(id: string) {
    return this.http.get<LoanOut>({ url: route(`/loans/${id}`) });
  }

  update(id: string, body: LoanUpdate) {
    return this.http.put<LoanUpdate, LoanOut>({ url: route(`/loans/${id}`), body });
  }

  delete(id: string) {
    return this.http.delete<void>({ url: route(`/loans/${id}`) });
  }

  return(id: string, body: LoanReturn) {
    return this.http.post<LoanReturn, LoanOut>({ url: route(`/loans/${id}/return`), body });
  }

  // Item-specific loan methods
  getItemLoans(itemId: string) {
    return this.http.get<LoanSummary[]>({ url: route(`/items/${itemId}/loans`) });
  }

  getItemCurrentLoan(itemId: string) {
    return this.http.get<LoanOut | null>({ url: route(`/items/${itemId}/current-loan`) });
  }
}
