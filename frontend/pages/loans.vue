<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import MdiPackageVariantClosedCheck from "~icons/mdi/package-variant-closed-check";
  import MdiAlertCircle from "~icons/mdi/alert-circle";
  import { Button, ButtonGroup } from "@/components/ui/button";
  import { Card } from "@/components/ui/card";
  import { Badge } from "@/components/ui/badge";
  import BaseContainer from "@/components/Base/Container.vue";
  import BaseSectionHeader from "@/components/Base/SectionHeader.vue";
  import type { LoanSummary } from "~~/lib/api/types/data-contracts";

  definePageMeta({
    middleware: ["auth"],
  });

  const { t } = useI18n();

  useHead({
    title: computed(() => `Hatchbridge Inventory | Active Loans`),
  });

  const api = useUserApi();
  const activeTab = ref<"active" | "overdue">("active");

  const { data: activeLoans, refresh: refreshActive } = useAsyncData("activeLoans", async () => {
    const { data, error } = await api.loans.getActive();
    if (error) {
      toast.error("Failed to load active loans");
      return [];
    }
    return data;
  });

  const { data: overdueLoans, refresh: refreshOverdue } = useAsyncData("overdueLoans", async () => {
    const { data, error } = await api.loans.getOverdue();
    if (error) {
      toast.error("Failed to load overdue loans");
      return [];
    }
    return data;
  });

  function formatDate(date: Date | string | undefined) {
    if (!date) return "N/A";
    return new Date(date).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  }

  function getDaysUntilDue(dueAt: Date | string) {
    const now = new Date();
    const due = new Date(dueAt);
    const diffTime = due.getTime() - now.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  }

  function getDueBadgeVariant(dueAt: Date | string): "destructive" | "default" | "secondary" {
    const days = getDaysUntilDue(dueAt);
    if (days < 0) return "destructive";
    if (days <= 2) return "default"; // Warning
    return "secondary";
  }

  const displayedLoans = computed(() => {
    return activeTab.value === "active" ? activeLoans.value : overdueLoans.value;
  });
</script>

<template>
  <BaseContainer>
    <div class="mb-4 flex flex-wrap items-center justify-between gap-4">
      <BaseSectionHeader>
        <div class="flex items-center gap-2">
          <MdiPackageVariantClosedCheck class="size-7" />
          Equipment Loans
        </div>
      </BaseSectionHeader>
    </div>

    <!-- Tab Buttons -->
    <div class="mb-4">
      <ButtonGroup>
        <Button 
          :variant="activeTab === 'active' ? 'default' : 'outline'" 
          size="sm"
          @click="activeTab = 'active'"
        >
          Active
          <Badge v-if="activeLoans?.length" variant="secondary" class="ml-2">
            {{ activeLoans.length }}
          </Badge>
        </Button>
        <Button 
          :variant="activeTab === 'overdue' ? 'default' : 'outline'" 
          size="sm"
          @click="activeTab = 'overdue'"
        >
          <MdiAlertCircle class="mr-1 size-4" />
          Overdue
          <Badge v-if="overdueLoans?.length" variant="destructive" class="ml-2">
            {{ overdueLoans.length }}
          </Badge>
        </Button>
      </ButtonGroup>
    </div>

    <!-- Loans List -->
    <div v-if="displayedLoans && displayedLoans.length > 0" class="space-y-3">
      <Card 
        v-for="loan in displayedLoans" 
        :key="loan.id" 
        class="p-4"
        :class="{ 'border-destructive/50 bg-destructive/5': activeTab === 'overdue' }"
      >
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div class="flex-1">
            <div class="flex items-center gap-2">
              <MdiAlertCircle v-if="activeTab === 'overdue'" class="size-4 text-destructive" />
              <NuxtLink 
                :to="`/item/${loan.itemId}`" 
                class="font-medium hover:underline"
              >
                {{ loan.itemName }}
              </NuxtLink>
              <Badge :variant="getDueBadgeVariant(loan.dueAt)">
                {{ getDaysUntilDue(loan.dueAt) >= 0 
                  ? `Due in ${getDaysUntilDue(loan.dueAt)} days` 
                  : `${Math.abs(getDaysUntilDue(loan.dueAt))} days overdue` 
                }}
              </Badge>
            </div>
            <div class="mt-1 text-sm text-muted-foreground">
              Borrowed by 
              <NuxtLink 
                :to="`/borrower/${loan.borrowerId}`" 
                class="font-medium hover:underline"
              >
                {{ loan.borrowerName }}
              </NuxtLink>
            </div>
            <div class="mt-1 text-xs text-muted-foreground">
              Checked out: {{ formatDate(loan.checkedOutAt) }} · 
              Due: {{ formatDate(loan.dueAt) }}
              <span v-if="loan.quantity > 1"> · Qty: {{ loan.quantity }}</span>
            </div>
          </div>
          <NuxtLink :to="`/item/${loan.itemId}`">
            <Button :variant="activeTab === 'overdue' ? 'default' : 'outline'" size="sm">
              {{ activeTab === 'overdue' ? 'Process Return' : 'View Item' }}
            </Button>
          </NuxtLink>
        </div>
      </Card>
    </div>

    <!-- Empty States -->
    <div 
      v-else-if="activeTab === 'active'" 
      class="flex flex-col items-center justify-center py-12 text-center"
    >
      <MdiPackageVariantClosedCheck class="mb-4 size-16 text-muted-foreground" />
      <p class="text-muted-foreground">No active loans. All equipment is checked in!</p>
    </div>

    <div 
      v-else 
      class="flex flex-col items-center justify-center py-12 text-center"
    >
      <MdiAlertCircle class="mb-4 size-16 text-green-500" />
      <p class="text-muted-foreground">No overdue loans. Great job!</p>
    </div>
  </BaseContainer>
</template>
