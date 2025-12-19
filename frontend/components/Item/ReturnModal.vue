<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import MdiLoading from "~icons/mdi/loading";
  import { Button } from "@/components/ui/button";
  import { Textarea } from "@/components/ui/textarea";
  import { Label } from "@/components/ui/label";
  import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog";
  import { useDialog } from "@/components/ui/dialog-provider";
  import { DialogID } from "~/components/ui/dialog-provider/utils";
  import type { LoanReturn, LoanOut } from "~~/lib/api/types/data-contracts";

  const props = defineProps<{
    loan: LoanOut | null;
    itemName: string;
  }>();

  const emit = defineEmits<{
    returned: [];
  }>();

  const { t } = useI18n();
  const api = useUserApi();
  const { closeDialog, registerOpenDialogCallback } = useDialog();

  const loading = ref(false);
  const model = ref(false);

  const formData = reactive<{
    returnNotes: string;
  }>({
    returnNotes: "",
  });

  function resetForm() {
    formData.returnNotes = "";
  }

  registerOpenDialogCallback(DialogID.ReturnLoan, () => {
    resetForm();
    model.value = true;
  });

  async function submit() {
    if (!props.loan) {
      toast.error("No active loan found");
      return;
    }

    loading.value = true;
    const payload: LoanReturn = {
      id: props.loan.id,
      returnNotes: formData.returnNotes || undefined,
    };

    const { error } = await api.loans.return(props.loan.id, payload);
    loading.value = false;

    if (error) {
      toast.error("Failed to return item");
      return;
    }

    toast.success(`"${props.itemName}" returned successfully`);
    closeDialog(DialogID.ReturnLoan);
    model.value = false;
    emit("returned");
  }

  function handleClose() {
    closeDialog(DialogID.ReturnLoan);
    model.value = false;
  }

  function formatDate(date: Date | string | undefined) {
    if (!date) return "N/A";
    return new Date(date).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  }
</script>

<template>
  <Dialog :dialog-id="DialogID.ReturnLoan" v-model:open="model" @update:open="val => !val && handleClose()">
    <DialogContent class="sm:max-w-[400px]">
      <DialogHeader>
        <DialogTitle>Return Equipment</DialogTitle>
        <DialogDescription>
          Check in "{{ itemName }}" from {{ loan?.borrowerName || "borrower" }}.
        </DialogDescription>
      </DialogHeader>

      <div v-if="loan" class="rounded-md bg-muted p-3 text-sm">
        <div class="grid gap-1">
          <div class="flex justify-between">
            <span class="text-muted-foreground">Checked out:</span>
            <span>{{ formatDate(loan.checkedOutAt) }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-muted-foreground">Due date:</span>
            <span :class="{ 'text-destructive font-medium': loan.isOverdue }">
              {{ formatDate(loan.dueAt) }}
              <span v-if="loan.isOverdue"> (Overdue)</span>
            </span>
          </div>
          <div v-if="loan.quantity > 1" class="flex justify-between">
            <span class="text-muted-foreground">Quantity:</span>
            <span>{{ loan.quantity }}</span>
          </div>
        </div>
      </div>

      <form @submit.prevent="submit" class="space-y-4">
        <div class="grid gap-2">
          <Label for="returnNotes">Return Notes</Label>
          <Textarea
            id="returnNotes"
            v-model="formData.returnNotes"
            placeholder="Optional notes about the return condition..."
            rows="2"
          />
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleClose">
            Cancel
          </Button>
          <Button type="submit" :disabled="loading || !loan">
            <MdiLoading v-if="loading" class="mr-2 animate-spin" />
            Return Item
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
