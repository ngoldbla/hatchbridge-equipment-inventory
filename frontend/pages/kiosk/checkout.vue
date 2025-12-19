<script setup lang="ts">
import { toast } from "@/components/ui/sonner";
import MdiPackageVariantClosedCheck from "~icons/mdi/package-variant-closed-check";
import MdiMagnify from "~icons/mdi/magnify";
import MdiQrcodeScan from "~icons/mdi/qrcode-scan";
import MdiChevronLeft from "~icons/mdi/chevron-left";
import MdiCheck from "~icons/mdi/check";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import type { ItemOut, BorrowerSummary, LoanCreate } from "~~/lib/api/types/data-contracts";

definePageMeta({
  layout: "kiosk",
  middleware: ["auth"],
});

useHead({
  title: "Check Out | Kiosk",
});

const api = useUserApi();

// State
const step = ref<"search-item" | "select-borrower" | "confirm">("search-item");
const itemSearch = ref("");
const borrowerSearch = ref("");
const selectedItem = ref<ItemOut | null>(null);
const selectedBorrower = ref<BorrowerSummary | null>(null);
const loading = ref(false);
const dueDate = ref<string>("");

// Set default due date to 7 days from now
onMounted(() => {
  const date = new Date();
  date.setDate(date.getDate() + 7);
  dueDate.value = date.toISOString().split("T")[0];
});

// Search results
const { data: itemResults, pending: searchingItems } = useAsyncData(
  "kioskItemSearch",
  async () => {
    if (!itemSearch.value || itemSearch.value.length < 2) return [];
    const { data } = await api.items.getAll({ q: itemSearch.value, pageSize: 10 });
    return data?.items || [];
  },
  { watch: [itemSearch] }
);

const { data: borrowers } = useAsyncData("kioskBorrowers", async () => {
  const { data } = await api.borrowers.getAll();
  return data || [];
});

const filteredBorrowers = computed(() => {
  if (!borrowers.value) return [];
  if (!borrowerSearch.value) return borrowers.value.filter(b => b.isActive);
  
  const query = borrowerSearch.value.toLowerCase();
  return borrowers.value.filter(
    b => b.isActive && (
      b.name.toLowerCase().includes(query) ||
      b.email.toLowerCase().includes(query) ||
      b.studentId?.toLowerCase().includes(query)
    )
  );
});

function selectItem(item: ItemOut) {
  selectedItem.value = item;
  step.value = "select-borrower";
}

function selectBorrower(borrower: BorrowerSummary) {
  selectedBorrower.value = borrower;
  step.value = "confirm";
}

function goBack() {
  if (step.value === "select-borrower") {
    step.value = "search-item";
    selectedItem.value = null;
  } else if (step.value === "confirm") {
    step.value = "select-borrower";
    selectedBorrower.value = null;
  }
}

async function handleCheckout() {
  if (!selectedItem.value || !selectedBorrower.value || !dueDate.value) {
    toast.error("Please complete all fields");
    return;
  }

  loading.value = true;
  try {
    const loanData: LoanCreate = {
      itemId: selectedItem.value.id,
      borrowerId: selectedBorrower.value.id,
      dueAt: new Date(dueDate.value).toISOString(),
      quantity: 1,
      notes: "Checked out via kiosk self-service",
    };

    const { error } = await api.loans.create(loanData);
    
    if (error) {
      toast.error("Failed to checkout item. It may already be on loan.");
      return;
    }

    toast.success(`Successfully checked out ${selectedItem.value.name}`);
    
    // Reset and go back to start
    selectedItem.value = null;
    selectedBorrower.value = null;
    itemSearch.value = "";
    borrowerSearch.value = "";
    step.value = "search-item";
  } finally {
    loading.value = false;
  }
}

function resetFlow() {
  selectedItem.value = null;
  selectedBorrower.value = null;
  itemSearch.value = "";
  borrowerSearch.value = "";
  step.value = "search-item";
}
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-6">
    <!-- Header -->
    <div class="flex items-center gap-4">
      <Button 
        v-if="step !== 'search-item'" 
        variant="ghost" 
        size="icon"
        @click="goBack"
      >
        <MdiChevronLeft class="size-6" />
      </Button>
      <div>
        <h1 class="text-2xl font-bold">Check Out Equipment</h1>
        <p class="text-muted-foreground">
          {{ step === 'search-item' ? 'Step 1: Find your item' : 
             step === 'select-borrower' ? 'Step 2: Select your profile' : 
             'Step 3: Confirm checkout' }}
        </p>
      </div>
    </div>

    <!-- Step 1: Search for Item -->
    <Card v-if="step === 'search-item'">
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <MdiMagnify class="size-5" />
          Find Equipment
        </CardTitle>
        <CardDescription>
          Search by name, description, or asset tag
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="flex gap-2">
          <Input
            v-model="itemSearch"
            placeholder="Type to search..."
            class="text-lg"
            autofocus
          />
          <Button variant="outline" size="icon" disabled>
            <MdiQrcodeScan class="size-5" />
          </Button>
        </div>

        <!-- Search Results -->
        <div v-if="searchingItems" class="py-8 text-center text-muted-foreground">
          Searching...
        </div>
        <div v-else-if="itemResults && itemResults.length > 0" class="space-y-2">
          <button
            v-for="item in itemResults"
            :key="item.id"
            class="w-full rounded-lg border p-4 text-left transition-colors hover:bg-accent"
            @click="selectItem(item)"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <p class="font-medium">{{ item.name }}</p>
                <p class="text-sm text-muted-foreground">
                  {{ item.location?.name || 'No location' }}
                </p>
              </div>
              <Badge v-if="item.assetId" variant="outline">
                #{{ item.assetId }}
              </Badge>
            </div>
          </button>
        </div>
        <div 
          v-else-if="itemSearch.length >= 2" 
          class="py-8 text-center text-muted-foreground"
        >
          No items found matching "{{ itemSearch }}"
        </div>
        <div 
          v-else-if="itemSearch.length > 0" 
          class="py-8 text-center text-muted-foreground"
        >
          Type at least 2 characters to search
        </div>
      </CardContent>
    </Card>

    <!-- Step 2: Select Borrower -->
    <Card v-if="step === 'select-borrower'">
      <CardHeader>
        <CardTitle>Select Your Profile</CardTitle>
        <CardDescription>
          Choose your borrower profile. If you don't have one, 
          <NuxtLink to="/kiosk/register" class="underline">register here</NuxtLink>.
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <Input
          v-model="borrowerSearch"
          placeholder="Search by name, email, or student ID..."
          class="text-lg"
        />

        <div class="max-h-64 space-y-2 overflow-y-auto">
          <button
            v-for="borrower in filteredBorrowers"
            :key="borrower.id"
            class="w-full rounded-lg border p-4 text-left transition-colors hover:bg-accent"
            @click="selectBorrower(borrower)"
          >
            <p class="font-medium">{{ borrower.name }}</p>
            <p class="text-sm text-muted-foreground">{{ borrower.email }}</p>
            <p v-if="borrower.studentId" class="text-xs text-muted-foreground">
              Student ID: {{ borrower.studentId }}
            </p>
          </button>
        </div>

        <div 
          v-if="filteredBorrowers.length === 0" 
          class="py-4 text-center text-muted-foreground"
        >
          No matching profiles found
        </div>
      </CardContent>
    </Card>

    <!-- Step 3: Confirm -->
    <Card v-if="step === 'confirm'">
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <MdiCheck class="size-5 text-green-500" />
          Confirm Checkout
        </CardTitle>
      </CardHeader>
      <CardContent class="space-y-6">
        <!-- Item Summary -->
        <div class="rounded-lg bg-muted p-4">
          <p class="text-sm font-medium text-muted-foreground">Item</p>
          <p class="text-lg font-semibold">{{ selectedItem?.name }}</p>
          <p class="text-sm text-muted-foreground">
            {{ selectedItem?.location?.name || 'No location' }}
          </p>
        </div>

        <!-- Borrower Summary -->
        <div class="rounded-lg bg-muted p-4">
          <p class="text-sm font-medium text-muted-foreground">Borrower</p>
          <p class="text-lg font-semibold">{{ selectedBorrower?.name }}</p>
          <p class="text-sm text-muted-foreground">{{ selectedBorrower?.email }}</p>
        </div>

        <!-- Due Date -->
        <div class="space-y-2">
          <Label for="due-date">Due Date</Label>
          <Input 
            id="due-date"
            v-model="dueDate" 
            type="date" 
            class="text-lg"
          />
        </div>

        <!-- Actions -->
        <div class="flex gap-4">
          <Button 
            variant="outline" 
            class="flex-1" 
            @click="resetFlow"
            :disabled="loading"
          >
            Cancel
          </Button>
          <Button 
            class="flex-1" 
            @click="handleCheckout"
            :disabled="loading || !dueDate"
          >
            <MdiPackageVariantClosedCheck class="mr-2 size-5" />
            {{ loading ? 'Processing...' : 'Confirm Checkout' }}
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
