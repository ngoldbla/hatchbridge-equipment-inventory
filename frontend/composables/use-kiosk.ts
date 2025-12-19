import type { KioskStatus } from "~~/lib/api/classes/kiosk";

/**
 * Composable for managing kiosk mode state and operations
 */
export function useKiosk() {
  const api = useUserApi();
  
  // Reactive kiosk status
  const status = ref<KioskStatus | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  /**
   * Check if currently in kiosk mode
   */
  const isKioskMode = computed(() => status.value?.isActive ?? false);

  /**
   * Check if kiosk is temporarily unlocked with admin access
   */
  const isUnlocked = computed(() => status.value?.isUnlocked ?? false);

  /**
   * When the unlock expires (if unlocked)
   */
  const unlockedUntil = computed(() => {
    if (!status.value?.unlockedUntil) return null;
    return new Date(status.value.unlockedUntil);
  });

  /**
   * Refresh the current kiosk status from the server
   */
  async function refreshStatus() {
    loading.value = true;
    error.value = null;
    try {
      const { data, error: apiError } = await api.kiosk.status();
      if (apiError) {
        error.value = "Failed to get kiosk status";
        return;
      }
      status.value = data;
    } finally {
      loading.value = false;
    }
  }

  /**
   * Activate kiosk mode and navigate to kiosk interface
   */
  async function activate() {
    loading.value = true;
    error.value = null;
    try {
      const { data, error: apiError } = await api.kiosk.activate();
      if (apiError) {
        error.value = "Failed to activate kiosk mode";
        return false;
      }
      status.value = data;
      // Navigate to kiosk interface
      await navigateTo("/kiosk");
      return true;
    } finally {
      loading.value = false;
    }
  }

  /**
   * Deactivate kiosk mode and return to normal interface
   */
  async function deactivate() {
    loading.value = true;
    error.value = null;
    try {
      const { data, error: apiError } = await api.kiosk.deactivate();
      if (apiError) {
        error.value = "Failed to deactivate kiosk mode";
        return false;
      }
      status.value = data;
      // Navigate back to home
      await navigateTo("/home");
      return true;
    } finally {
      loading.value = false;
    }
  }

  /**
   * Temporarily unlock kiosk mode with admin password
   * @param password Admin user's password
   * @param durationMinutes How long to stay unlocked (default: 5 minutes)
   */
  async function unlock(password: string, durationMinutes = 5): Promise<{ success: boolean; error?: string }> {
    loading.value = true;
    error.value = null;
    try {
      const { data, error: apiError } = await api.kiosk.unlock({
        password,
        durationMinutes,
      });
      if (apiError) {
        const errorMessage = "Invalid password or unlock failed";
        error.value = errorMessage;
        return { success: false, error: errorMessage };
      }
      status.value = data;
      return { success: true };
    } finally {
      loading.value = false;
    }
  }

  /**
   * Lock the kiosk (revoke temporary admin access)
   */
  async function lock() {
    loading.value = true;
    error.value = null;
    try {
      const { data, error: apiError } = await api.kiosk.lock();
      if (apiError) {
        error.value = "Failed to lock kiosk";
        return false;
      }
      status.value = data;
      return true;
    } finally {
      loading.value = false;
    }
  }

  // Auto-refresh status on initial load
  onMounted(() => {
    refreshStatus();
  });

  return {
    // State
    status,
    loading,
    error,
    
    // Computed
    isKioskMode,
    isUnlocked,
    unlockedUntil,
    
    // Actions
    refreshStatus,
    activate,
    deactivate,
    unlock,
    lock,
  };
}
