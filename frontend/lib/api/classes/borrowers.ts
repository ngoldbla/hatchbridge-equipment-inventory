import { BaseAPI, route } from "../base";
import type { BorrowerCreate, BorrowerOut, BorrowerSummary, BorrowerUpdate, LoanSummary } from "../types/data-contracts";

export class BorrowersApi extends BaseAPI {
  getAll() {
    return this.http.get<BorrowerSummary[]>({ url: route("/borrowers") });
  }

  getActive() {
    return this.http.get<BorrowerSummary[]>({ url: route("/borrowers/active") });
  }

  create(body: BorrowerCreate) {
    return this.http.post<BorrowerCreate, BorrowerOut>({ url: route("/borrowers"), body });
  }

  get(id: string) {
    return this.http.get<BorrowerOut>({ url: route(`/borrowers/${id}`) });
  }

  update(id: string, body: BorrowerUpdate) {
    return this.http.put<BorrowerUpdate, BorrowerOut>({ url: route(`/borrowers/${id}`), body });
  }

  delete(id: string) {
    return this.http.delete<void>({ url: route(`/borrowers/${id}`) });
  }

  getLoans(id: string) {
    return this.http.get<LoanSummary[]>({ url: route(`/borrowers/${id}/loans`) });
  }
}
