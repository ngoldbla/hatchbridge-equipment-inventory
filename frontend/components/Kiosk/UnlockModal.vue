<script setup lang="ts">
import { toast } from "@/components/ui/sonner";
import MdiShieldLock from "~icons/mdi/shield-lock";
import MdiEyeOff from "~icons/mdi/eye-off";
import MdiEye from "~icons/mdi/eye";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface Props {
  open: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "unlocked"): void;
}>();

const { unlock, loading } = useKiosk();

const password = ref("");
const showPassword = ref(false);
const errorMessage = ref("");

async function handleUnlock() {
  if (!password.value.trim()) {
    errorMessage.value = "Password is required";
    return;
  }

  errorMessage.value = "";
  const result = await unlock(password.value);

  if (result.success) {
    toast.success("Admin access granted for 5 minutes");
    password.value = "";
    emit("update:open", false);
    emit("unlocked");
  } else {
    errorMessage.value = result.error || "Invalid password";
    toast.error("Invalid password");
  }
}

function handleClose() {
  password.value = "";
  errorMessage.value = "";
  emit("update:open", false);
}
</script>

<template>
  <Dialog :open="props.open" @update:open="handleClose">
    <DialogContent class="sm:max-w-[400px]">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2">
          <MdiShieldLock class="size-5 text-yellow-500" />
          Admin Unlock
        </DialogTitle>
        <DialogDescription>
          Enter your admin password to temporarily access restricted features.
          Access will expire after 5 minutes.
        </DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleUnlock" class="space-y-4 py-4">
        <div class="space-y-2">
          <Label for="admin-password">Password</Label>
          <div class="relative">
            <Input
              id="admin-password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Enter your password"
              class="pr-10"
              :disabled="loading"
              autocomplete="current-password"
            />
            <button
              type="button"
              class="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              @click="showPassword = !showPassword"
            >
              <MdiEye v-if="!showPassword" class="size-5" />
              <MdiEyeOff v-else class="size-5" />
            </button>
          </div>
          <p v-if="errorMessage" class="text-sm text-destructive">
            {{ errorMessage }}
          </p>
        </div>
      </form>

      <DialogFooter>
        <Button variant="outline" @click="handleClose" :disabled="loading">
          Cancel
        </Button>
        <Button @click="handleUnlock" :disabled="loading || !password.trim()">
          <MdiShieldLock class="mr-2 size-4" />
          {{ loading ? "Unlocking..." : "Unlock" }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
