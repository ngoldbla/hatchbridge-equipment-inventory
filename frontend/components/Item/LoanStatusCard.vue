<script setup lang="ts">
  import MdiPackageVariantClosedCheck from "~icons/mdi/package-variant-closed-check";
  import MdiAlertCircle from "~icons/mdi/alert-circle";
  import MdiArrowRight from "~icons/mdi/arrow-right";
  import { Button } from "@/components/ui/button";
  import { Card } from "@/components/ui/card";
  import { Badge } from "@/components/ui/badge";
  import { useDialog } from "@/components/ui/dialog-provider";
  import { DialogID } from "~/components/ui/dialog-provider/utils";
  import type { LoanOut } from "~~/lib/api/types/data-contracts";

  const props = defineProps<{
    loading: boolean;
    currentLoan: LoanOut | null;
    itemName: string;
  }>();

  const emit = defineEmits<{
    refresh: [];
  }>();

  const { openDialog } = useDialog();

  function formatDate(date: Date | string | undefined) {
    if (!date) return "N/A";
    return new Date(date).toLocaleDateString("en-US", {
      weekday: "short",
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

  const daysUntilDue = computed(() => {
    if (!props.currentLoan) return 0;
    return getDaysUntilDue(props.currentLoan.dueAt);
  });

  const isOverdue = computed(() => {
    return props.currentLoan?.isOverdue || daysUntilDue.value < 0;
  });

  function handleCheckout() {
    openDialog(DialogID.CreateLoan);
  }

  function handleReturn() {
    openDialog(DialogID.ReturnLoan);
  }
</script>

<template>
  <Card class="p-4">
    <div class="flex items-start justify-between gap-4">
      <!-- Available State -->
      <div v-if="!currentLoan" class="flex items-center gap-3">
        <div class="flex size-10 items-center justify-center rounded-full bg-green-100 text-green-600 dark:bg-green-900/30">
          <MdiPackageVariantClosedCheck class="size-5" />
        </div>
        <div>
          <h3 class="font-medium">Available</h3>
          <p class="text-sm text-muted-foreground">This item is ready to be checked out</p>
        </div>
      </div>

      <!-- On Loan State -->
      <div v-else class="flex items-center gap-3">
        <div 
          class="flex size-10 items-center justify-center rounded-full"
          :class="isOverdue 
            ? 'bg-destructive/10 text-destructive' 
            : 'bg-amber-100 text-amber-600 dark:bg-amber-900/30'"
        >
          <MdiAlertCircle v-if="isOverdue" class="size-5" />
          <MdiPackageVariantClosedCheck v-else class="size-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h3 class="font-medium">On Loan</h3>
            <Badge :variant="isOverdue ? 'destructive' : 'secondary'">
              {{ isOverdue 
                ? `${Math.abs(daysUntilDue)} days overdue` 
                : `Due in ${daysUntilDue} days` 
              }}
            </Badge>
          </div>
          <p class="text-sm text-muted-foreground">
            Borrowed by 
            <NuxtLink 
              :to="`/borrower/${currentLoan.borrowerId}`" 
              class="font-medium underline hover:text-foreground"
            >
              {{ currentLoan.borrowerName }}
            </NuxtLink>
            <span v-if="currentLoan.quantity > 1"> (qty: {{ currentLoan.quantity }})</span>
          </p>
          <p class="text-xs text-muted-foreground">
            Checked out: {{ formatDate(currentLoan.checkedOutAt) }} · Due: {{ formatDate(currentLoan.dueAt) }}
          </p>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex-shrink-0">
        <Button v-if="!currentLoan" @click="handleCheckout">
          Check Out
          <MdiArrowRight class="ml-1 size-4" />
        </Button>
        <Button v-else variant="default" @click="handleReturn">
          Return Item
        </Button>
      </div>
    </div>
  </Card>
</template>
