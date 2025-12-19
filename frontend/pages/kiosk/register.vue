<script setup lang="ts">
import { toast } from "@/components/ui/sonner";
import MdiAccountPlus from "~icons/mdi/account-plus";
import MdiCheck from "~icons/mdi/check";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { BorrowerCreate } from "~~/lib/api/types/data-contracts";

definePageMeta({
  layout: "kiosk",
  middleware: ["auth"],
});

useHead({
  title: "Register | Kiosk",
});

const api = useUserApi();

// Form state
const formData = reactive<BorrowerCreate>({
  name: "",
  email: "",
  phone: "",
  organization: "",
  studentId: "",
  notes: "",
});

const loading = ref(false);
const showSuccess = ref(false);
const errors = reactive<Record<string, string>>({});

function validateForm(): boolean {
  // Clear previous errors
  Object.keys(errors).forEach(key => delete errors[key]);

  if (!formData.name.trim()) {
    errors.name = "Name is required";
  }

  if (!formData.email.trim()) {
    errors.email = "Email is required";
  } else if (!isValidEmail(formData.email)) {
    errors.email = "Please enter a valid email address";
  }

  return Object.keys(errors).length === 0;
}

function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

async function handleSubmit() {
  if (!validateForm()) return;

  loading.value = true;
  try {
    // Note: We're setting selfRegistered flag via the notes for now
    // In a full implementation, this would be a separate field
    const dataToSubmit = {
      ...formData,
      notes: formData.notes 
        ? `${formData.notes} [Self-registered via kiosk]`
        : "[Self-registered via kiosk]",
    };

    const { error } = await api.borrowers.create(dataToSubmit);

    if (error) {
      if (error.message?.includes("duplicate") || error.message?.includes("unique")) {
        toast.error("A borrower with this email already exists");
        errors.email = "This email is already registered";
      } else {
        toast.error("Failed to register. Please try again.");
      }
      return;
    }

    // Show success
    showSuccess.value = true;
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  formData.name = "";
  formData.email = "";
  formData.phone = "";
  formData.organization = "";
  formData.studentId = "";
  formData.notes = "";
  showSuccess.value = false;
  Object.keys(errors).forEach(key => delete errors[key]);
}

function goToCheckout() {
  navigateTo("/kiosk/checkout");
}
</script>

<template>
  <div class="mx-auto max-w-lg space-y-6">
    <!-- Header -->
    <div class="text-center">
      <h1 class="text-2xl font-bold">Register as Borrower</h1>
      <p class="text-muted-foreground">
        Create your profile to borrow equipment
      </p>
    </div>

    <!-- Success State -->
    <Card v-if="showSuccess" class="border-green-500 bg-green-50 dark:bg-green-950">
      <CardContent class="py-8 text-center">
        <div class="mx-auto mb-4 flex size-16 items-center justify-center rounded-full bg-green-500 text-white">
          <MdiCheck class="size-8" />
        </div>
        <h2 class="text-xl font-semibold text-green-800 dark:text-green-200">
          Registration Successful!
        </h2>
        <p class="mt-2 text-green-600 dark:text-green-400">
          You can now check out equipment using your profile.
        </p>
        <div class="mt-6 flex gap-4">
          <Button variant="outline" class="flex-1" @click="resetForm">
            Register Another
          </Button>
          <Button class="flex-1" @click="goToCheckout">
            Check Out Equipment
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Registration Form -->
    <Card v-else>
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <MdiAccountPlus class="size-5" />
          Your Information
        </CardTitle>
        <CardDescription>
          Fields marked with * are required
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <!-- Name -->
          <div class="space-y-2">
            <Label for="name">Full Name *</Label>
            <Input
              id="name"
              v-model="formData.name"
              placeholder="John Doe"
              :class="errors.name ? 'border-destructive' : ''"
              autofocus
            />
            <p v-if="errors.name" class="text-sm text-destructive">
              {{ errors.name }}
            </p>
          </div>

          <!-- Email -->
          <div class="space-y-2">
            <Label for="email">Email Address *</Label>
            <Input
              id="email"
              v-model="formData.email"
              type="email"
              placeholder="john@example.com"
              :class="errors.email ? 'border-destructive' : ''"
            />
            <p v-if="errors.email" class="text-sm text-destructive">
              {{ errors.email }}
            </p>
          </div>

          <!-- Phone -->
          <div class="space-y-2">
            <Label for="phone">Phone Number</Label>
            <Input
              id="phone"
              v-model="formData.phone"
              type="tel"
              placeholder="(555) 123-4567"
            />
          </div>

          <!-- Organization -->
          <div class="space-y-2">
            <Label for="organization">Organization / Company</Label>
            <Input
              id="organization"
              v-model="formData.organization"
              placeholder="Your company or startup name"
            />
          </div>

          <!-- Student ID -->
          <div class="space-y-2">
            <Label for="studentId">Student ID</Label>
            <Input
              id="studentId"
              v-model="formData.studentId"
              placeholder="If applicable"
            />
          </div>

          <!-- Submit -->
          <Button 
            type="submit" 
            class="w-full" 
            size="lg"
            :disabled="loading"
          >
            <MdiAccountPlus class="mr-2 size-5" />
            {{ loading ? 'Registering...' : 'Complete Registration' }}
          </Button>
        </form>
      </CardContent>
    </Card>

    <!-- Already Registered -->
    <Card class="bg-muted/50">
      <CardContent class="py-4 text-center">
        <p class="text-sm text-muted-foreground">
          Already registered? 
          <NuxtLink to="/kiosk/checkout" class="font-medium underline">
            Go to checkout
          </NuxtLink>
        </p>
      </CardContent>
    </Card>
  </div>
</template>
