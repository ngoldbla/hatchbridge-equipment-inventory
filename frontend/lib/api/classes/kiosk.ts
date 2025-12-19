import { BaseAPI, route } from "../base";

export interface KioskStatus {
  isActive: boolean;
  isUnlocked: boolean;
  unlockedUntil?: string | null;
}

export interface KioskUnlockRequest {
  password: string;
  durationMinutes?: number;
}

export class KioskAPI extends BaseAPI {
  /**
   * Activate kiosk mode for the current user
   */
  activate() {
    return this.http.post<object, KioskStatus>({ url: route("/kiosk/activate"), body: {} });
  }

  /**
   * Deactivate kiosk mode for the current user
   */
  deactivate() {
    return this.http.post<object, KioskStatus>({ url: route("/kiosk/deactivate"), body: {} });
  }

  /**
   * Get current kiosk status for the current user
   */
  status() {
    return this.http.get<KioskStatus>({ url: route("/kiosk/status") });
  }

  /**
   * Temporarily unlock kiosk mode with admin password
   */
  unlock(data: KioskUnlockRequest) {
    return this.http.post<KioskUnlockRequest, KioskStatus>({ url: route("/kiosk/unlock"), body: data });
  }

  /**
   * Lock kiosk mode (revoke temporary admin access)
   */
  lock() {
    return this.http.post<object, KioskStatus>({ url: route("/kiosk/lock"), body: {} });
  }
}

