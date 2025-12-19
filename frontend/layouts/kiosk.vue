<template>
  <div id="kiosk-app" class="min-h-screen bg-background">
    <!-- Modal overlays -->
    <ModalConfirm />
    <KioskUnlockModal 
      v-model:open="showUnlockModal" 
      @unlocked="handleUnlocked"
    />
    
    <!-- Fixed header with admin controls -->
    <header class="fixed left-0 right-0 top-0 z-40 border-b bg-secondary shadow-sm">
      <div class="container mx-auto flex h-16 items-center justify-between px-4">
        <!-- Logo and title -->
        <div class="flex items-center gap-4">
          <NuxtLink to="/kiosk" class="flex items-center gap-3">
            <div class="flex size-10 items-center justify-center rounded-lg bg-white p-1.5 shadow-sm">
              <AppLogo />
            </div>
            <div class="hidden sm:block">
              <h1 class="text-lg font-bold leading-none">Equipment Self-Service</h1>
              <p class="text-xs text-muted-foreground">Hatchbridge Kiosk</p>
            </div>
          </NuxtLink>
        </div>

        <!-- Admin controls -->
        <div class="flex items-center gap-2">
          <!-- Unlock status indicator -->
          <Badge 
            v-if="kioskStatus?.isUnlocked" 
            variant="default"
            class="hidden sm:inline-flex"
          >
            <MdiShieldCheck class="mr-1 size-3" />
            Admin Access
          </Badge>

          <!-- Exit button (only when unlocked) -->
          <Button
            v-if="kioskStatus?.isUnlocked"
            variant="destructive"
            size="sm"
            @click="handleExitKiosk"
          >
            <MdiExitToApp class="mr-2 size-4" />
            Exit Kiosk
          </Button>

          <!-- Unlock button -->
          <Button
            variant="ghost"
            size="icon"
            class="rounded-full"
            @click="showUnlockModal = true"
            :title="kioskStatus?.isUnlocked ? 'Admin Access Active' : 'Admin Unlock'"
          >
            <MdiShieldLock 
              class="size-5" 
              :class="kioskStatus?.isUnlocked ? 'text-green-500' : 'text-muted-foreground'"
            />
          </Button>
        </div>
      </div>
    </header>

    <!-- Main content area with padding for fixed header and nav -->
    <main class="container mx-auto px-4 pb-24 pt-20">
      <slot />
    </main>

    <!-- Bottom navigation -->
    <nav class="fixed bottom-0 left-0 right-0 z-40 border-t bg-card shadow-lg">
      <div class="container mx-auto">
        <div class="grid grid-cols-4 gap-1 p-2">
          <NuxtLink 
            v-for="navItem in navItems" 
            :key="navItem.to"
            :to="navItem.to"
            class="flex flex-col items-center justify-center rounded-lg py-3 transition-colors"
            :class="isActive(navItem.to) 
              ? 'bg-primary/10 text-primary' 
              : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'"
          >
            <component :is="navItem.icon" class="size-6" />
            <span class="mt-1 text-xs font-medium">{{ navItem.label }}</span>
          </NuxtLink>
        </div>
      </div>
    </nav>
  </div>
</template>

<script lang="ts" setup>
import { toast } from "@/components/ui/sonner";
import MdiHome from "~icons/mdi/home";
import MdiPackageVariantClosedCheck from "~icons/mdi/package-variant-closed-check";
import MdiPackageVariantClosedRemove from "~icons/mdi/package-variant-closed-remove";
import MdiAccountPlus from "~icons/mdi/account-plus";
import MdiShieldLock from "~icons/mdi/shield-lock";
import MdiShieldCheck from "~icons/mdi/shield-check";
import MdiExitToApp from "~icons/mdi/exit-to-app";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import AppLogo from "~/components/App/Logo.vue";
import ModalConfirm from "~/components/ModalConfirm.vue";
import KioskUnlockModal from "~/components/Kiosk/UnlockModal.vue";

const route = useRoute();
const showUnlockModal = ref(false);
const { status: kioskStatus, deactivate, refreshStatus } = useKiosk();

// Navigation items for kiosk mode
const navItems = [
  {
    to: "/kiosk",
    label: "Home",
    icon: MdiHome,
  },
  {
    to: "/kiosk/checkout",
    label: "Check Out",
    icon: MdiPackageVariantClosedCheck,
  },
  {
    to: "/kiosk/return",
    label: "Return",
    icon: MdiPackageVariantClosedRemove,
  },
  {
    to: "/kiosk/register",
    label: "Register",
    icon: MdiAccountPlus,
  },
];

function isActive(path: string) {
  if (path === "/kiosk") {
    return route.path === "/kiosk";
  }
  return route.path.startsWith(path);
}

function handleUnlocked() {
  toast.success("Admin access granted");
  refreshStatus();
}

async function handleExitKiosk() {
  const confirmed = await useConfirm().openDialog({
    title: "Exit Kiosk Mode",
    message: "Are you sure you want to exit kiosk mode and return to the admin interface?",
    confirm: "Exit Kiosk Mode",
  });
  
  if (confirmed) {
    await deactivate();
  }
}

// Refresh status on mount
onMounted(() => {
  refreshStatus();
});

// Periodically check if unlock has expired
onMounted(() => {
  const interval = setInterval(() => {
    refreshStatus();
  }, 30000); // Check every 30 seconds

  onUnmounted(() => {
    clearInterval(interval);
  });
});
</script>
