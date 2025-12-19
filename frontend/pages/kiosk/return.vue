<script setup lang="ts">
import { toast } from "@/components/ui/sonner";
import MdiPackageVariantClosedRemove from "~icons/mdi/package-variant-closed-remove";
import MdiMagnify from "~icons/mdi/magnify";
import MdiCheck from "~icons/mdi/check";
import MdiAlertCircle from "~icons/mdi/alert-circle";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import type { LoanSummary } from "~~/lib/api/types/data-contracts";

definePageMeta({
  layout: "kiosk",
  middleware: ["auth"],
});

useHead({
  title: "Return | Kiosk",
});

const api = useUserApi();

// State
const search = ref("");
const selectedLoan = ref<LoanSummary | null>(null);
const returnNotes = ref("");
const loading = ref(false);
const showSuccess = ref(false);

// Fetch active loans
const { data: activeLoans, refresh: refreshLoans } = useAsyncData("kioskReturnLoans", async () => {
  const { data } = await api.loans.getActive();
  return data || [];
});

const { data: overdueLoans } = useAsyncData("kioskOverdueLoans", async () => {
  const { data } = await api.loans.getOverdue();
  return data || [];
});

// Filter loans based on search
const filteredLoans = computed(() => {
  const allLoans = [...(activeLoans.value || []), ...(overdueLoans.value || [])];
  
  // Deduplicate by loan ID
  const uniqueLoans = allLoans.filter((loan, index, self) => 
    index === self.findIndex(l => l.id === loan.id)
  );
  
  if (!search.value) return uniqueLoans;
  
  const query = search.value.toLowerCase();
  return uniqueLoans.filter(loan => 
    loan.item?.name?.toLowerCase().includes(query) ||
    loan.borrower?.name?.toLowerCase().includes(query) ||
    loan.borrower?.email?.toLowerCase().includes(query)
  );
});

function isOverdue(loan: LoanSummary): boolean {
  if (!loan.dueAt) return false;
  return new Date(loan.dueAt) < new Date();
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

function selectLoan(loan: LoanSummary) {
  selectedLoan.value = loan;
  returnNotes.value = "";
}

async function handleReturn() {
  if (!selectedLoan.value) return;

  loading.value = true;
  try {
    const { error } = await api.loans.return(selectedLoan.value.id, {
      notes: returnNotes.value || "Returned via kiosk self-service",
    });

    if (error) {
      toast.error("Failed to process return");
      return;
    }

    showSuccess.value = true;
    setTimeout(() => {
      showSuccess.value = false;
      selectedLoan.value = null;
      returnNotes.value = "";
      search.value = "";
      refreshLoans();
    }, 2000);
  } finally {
    loading.value = false;
  }
}

function cancelReturn() {
  selectedLoan.value = null;
  returnNotes.value = "";
}
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold">Return Equipment</h1>
      <p class="text-muted-foreground">
        Find your borrowed item and process the return
      </p>
    </div>

    <!-- Success Message -->
    <Card v-if="showSuccess" class="border-green-500 bg-green-50 dark:bg-green-950">
      <CardContent class="flex items-center gap-4 py-6">
        <div class="flex size-12 items-center justify-center rounded-full bg-green-500 text-white">
          <MdiCheck class="size-6" />
        </div>
        <div>
          <p class="text-lg font-semibold text-green-800 dark:text-green-200">
            Item Returned Successfully!
          </p>
          <p class="text-green-600 dark:text-green-400">
            Thank you for returning the equipment.
          </p>
        </div>
      </CardContent>
    </Card>

    <!-- Confirm Return -->
    <Card v-else-if="selectedLoan">
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <MdiPackageVariantClosedRemove class="size-5" />
          Confirm Return
        </CardTitle>
      </CardHeader>
      <CardContent class="space-y-6">
        <!-- Loan Summary -->
        <div class="rounded-lg bg-muted p-4">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-sm font-medium text-muted-foreground">Item</p>
              <p class="text-lg font-semibold">{{ selectedLoan.item?.name }}</p>
            </div>
            <Badge 
              v-if="isOverdue(selectedLoan)" 
              variant="destructive"
            >
              Overdue
            </Badge>
          </div>
        </div>

        <div class="rounded-lg bg-muted p-4">
          <p class="text-sm font-medium text-muted-foreground">Borrower</p>
          <p class="text-lg font-semibold">{{ selectedLoan.borrower?.name }}</p>
          <p class="text-sm text-muted-foreground">{{ selectedLoan.borrower?.email }}</p>
        </div>

        <div class="grid grid-cols-2 gap-4 rounded-lg bg-muted p-4">
          <div>
            <p class="text-sm font-medium text-muted-foreground">Checked Out</p>
            <p>{{ formatDate(selectedLoan.checkedOutAt) }}</p>
          </div>
          <div>
            <p class="text-sm font-medium text-muted-foreground">Due Date</p>
            <p :class="isOverdue(selectedLoan) ? 'text-destructive font-semibold' : ''">
              {{ formatDate(selectedLoan.dueAt) }}
            </p>
          </div>
        </div>

        <!-- Return Notes -->
        <div class="space-y-2">
          <Label for="return-notes">Return Notes (Optional)</Label>
          <Textarea
            id="return-notes"
            v-model="returnNotes"
            placeholder="Any notes about the item condition..."
            rows="3"
          />
        </div>

        <!-- Actions -->
        <div class="flex gap-4">
          <Button 
            variant="outline" 
            class="flex-1" 
            @click="cancelReturn"
            :disabled="loading"
          >
            Cancel
          </Button>
          <Button 
            class="flex-1" 
            @click="handleReturn"
            :disabled="loading"
          >
            <MdiCheck class="mr-2 size-5" />
            {{ loading ? 'Processing...' : 'Confirm Return' }}
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Loan Search/List -->
    <Card v-else>
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <MdiMagnify class="size-5" />
          Find Your Loan
        </CardTitle>
        <CardDescription>
          Search by item name, borrower name, or email
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <Input
          v-model="search"
          placeholder="Type to search..."
          class="text-lg"
          autofocus
        />

        <!-- Loan List -->
        <div class="max-h-96 space-y-2 overflow-y-auto">
          <button
            v-for="loan in filteredLoans"
            :key="loan.id"
            class="w-full rounded-lg border p-4 text-left transition-colors hover:bg-accent"
            :class="isOverdue(loan) ? 'border-destructive/50' : ''"
            @click="selectLoan(loan)"
          >
            <div class="flex items-start justify-between gap-4">
              <div class="min-w-0 flex-1">
                <p class="font-medium">{{ loan.item?.name }}</p>
                <p class="text-sm text-muted-foreground">
                  {{ loan.borrower?.name }} • {{ loan.borrower?.email }}
                </p>
                <p class="mt-1 text-xs text-muted-foreground">
                  Due: {{ formatDate(loan.dueAt) }}
                </p>
              </div>
              <div class="flex flex-col items-end gap-1">
                <Badge 
                  v-if="isOverdue(loan)" 
                  variant="destructive"
                  class="flex items-center gap-1"
                >
                  <MdiAlertCircle class="size-3" />
                  Overdue
                </Badge>
                <Badge v-else variant="secondary">
                  Active
                </Badge>
              </div>
            </div>
          </button>
        </div>

        <div 
          v-if="filteredLoans.length === 0" 
          class="py-8 text-center text-muted-foreground"
        >
          {{ search ? 'No matching loans found' : 'No active loans' }}
        </div>
      </CardContent>
    </Card>
  </div>
</template>
