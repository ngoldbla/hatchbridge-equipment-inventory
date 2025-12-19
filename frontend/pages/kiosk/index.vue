<script setup lang="ts">
import MdiPackageVariantClosedCheck from "~icons/mdi/package-variant-closed-check";
import MdiPackageVariantClosedRemove from "~icons/mdi/package-variant-closed-remove";
import MdiAccountPlus from "~icons/mdi/account-plus";
import MdiClockOutline from "~icons/mdi/clock-outline";
import MdiAlertCircle from "~icons/mdi/alert-circle";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

definePageMeta({
  layout: "kiosk",
  middleware: ["auth"],
});

useHead({
  title: "Kiosk | Hatchbridge Equipment",
});

const api = useUserApi();

// Fetch summary data for the kiosk dashboard
const { data: activeLoans } = useAsyncData("kioskActiveLoans", async () => {
  const { data } = await api.loans.getActive();
  return data || [];
});

const { data: overdueLoans } = useAsyncData("kioskOverdueLoans", async () => {
  const { data } = await api.loans.getOverdue();
  return data || [];
});

const actionCards = [
  {
    to: "/kiosk/checkout",
    icon: MdiPackageVariantClosedCheck,
    title: "Check Out Equipment",
    description: "Borrow equipment for your project",
    color: "bg-blue-500",
  },
  {
    to: "/kiosk/return",
    icon: MdiPackageVariantClosedRemove,
    title: "Return Equipment",
    description: "Return borrowed items",
    color: "bg-green-500",
  },
  {
    to: "/kiosk/register",
    icon: MdiAccountPlus,
    title: "Register as Borrower",
    description: "Create your borrower profile",
    color: "bg-purple-500",
  },
];
</script>

<template>
  <div class="space-y-8">
    <!-- Welcome section -->
    <div class="text-center">
      <h1 class="text-3xl font-bold tracking-tight sm:text-4xl">
        Welcome to Equipment Self-Service
      </h1>
      <p class="mt-2 text-muted-foreground">
        Choose an action below to get started
      </p>
    </div>

    <!-- Action Cards -->
    <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <NuxtLink
        v-for="card in actionCards"
        :key="card.to"
        :to="card.to"
        class="block transition-transform hover:scale-[1.02]"
      >
        <Card class="h-full cursor-pointer transition-colors hover:bg-accent/50">
          <CardHeader class="text-center">
            <div
              class="mx-auto mb-4 flex size-16 items-center justify-center rounded-full text-white"
              :class="card.color"
            >
              <component :is="card.icon" class="size-8" />
            </div>
            <CardTitle class="text-xl">{{ card.title }}</CardTitle>
            <CardDescription>{{ card.description }}</CardDescription>
          </CardHeader>
        </Card>
      </NuxtLink>
    </div>

    <!-- Status Section -->
    <div class="grid gap-6 sm:grid-cols-2">
      <!-- Active Loans -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <div class="flex items-center gap-2">
            <MdiClockOutline class="size-5 text-blue-500" />
            <CardTitle class="text-base">Active Loans</CardTitle>
          </div>
          <Badge variant="secondary">
            {{ activeLoans?.length || 0 }}
          </Badge>
        </CardHeader>
        <CardContent>
          <p class="text-sm text-muted-foreground">
            {{ activeLoans?.length || 0 }} items currently checked out
          </p>
        </CardContent>
      </Card>

      <!-- Overdue Items -->
      <Card :class="overdueLoans?.length ? 'border-destructive/50' : ''">
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <div class="flex items-center gap-2">
            <MdiAlertCircle 
              class="size-5" 
              :class="overdueLoans?.length ? 'text-destructive' : 'text-green-500'" 
            />
            <CardTitle class="text-base">Overdue Items</CardTitle>
          </div>
          <Badge 
            :variant="overdueLoans?.length ? 'destructive' : 'secondary'"
          >
            {{ overdueLoans?.length || 0 }}
          </Badge>
        </CardHeader>
        <CardContent>
          <p class="text-sm text-muted-foreground">
            {{ overdueLoans?.length 
              ? `${overdueLoans.length} items need to be returned` 
              : 'No overdue items - great job!' 
            }}
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- Instructions -->
    <Card class="bg-muted/50">
      <CardContent class="py-4">
        <p class="text-center text-sm text-muted-foreground">
          Need help? Contact the equipment manager or tap the 
          <span class="font-medium">Admin Unlock</span> button in the top right corner.
        </p>
      </CardContent>
    </Card>
  </div>
</template>
