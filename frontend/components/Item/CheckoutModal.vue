<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import MdiLoading from "~icons/mdi/loading";
  import MdiCalendar from "~icons/mdi/calendar";
  import { Button } from "@/components/ui/button";
  import { Input } from "@/components/ui/input";
  import { Label } from "@/components/ui/label";
  import { Textarea } from "@/components/ui/textarea";
  import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
  } from "@/components/ui/select";
  import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog";
  import { useDialog } from "@/components/ui/dialog-provider";
  import { DialogID } from "~/components/ui/dialog-provider/utils";
  import type { BorrowerSummary, LoanCreate } from "~~/lib/api/types/data-contracts";

  const props = defineProps<{
    itemId: string;
    itemName: string;
    itemQuantity: number;
  }>();

  const emit = defineEmits<{
    created: [];
  }>();

  const { t } = useI18n();
  const api = useUserApi();
  const { closeDialog, registerOpenDialogCallback } = useDialog();

  const loading = ref(false);
  const model = ref(false);
  const borrowers = ref<BorrowerSummary[]>([]);
  const loadingBorrowers = ref(false);

  const formData = reactive<{
    borrowerId: string;
    dueAt: string;
    notes: string;
    quantity: number;
  }>({
    borrowerId: "",
    dueAt: "",
    notes: "",
    quantity: 1,
  });

  function resetForm() {
    formData.borrowerId = "";
    formData.dueAt = getDefaultDueDate();
    formData.notes = "";
    formData.quantity = 1;
  }

  function getDefaultDueDate(): string {
    const date = new Date();
    date.setDate(date.getDate() + 14); // Default 2 weeks
    const result = date.toISOString().split("T")[0];
    return result || "";
  }

  async function fetchBorrowers() {
    loadingBorrowers.value = true;
    const { data, error } = await api.borrowers.getActive();
    loadingBorrowers.value = false;

    if (error) {
      toast.error("Failed to load borrowers");
      return;
    }

    borrowers.value = data;
  }

  registerOpenDialogCallback(DialogID.CreateLoan, () => {
    resetForm();
    fetchBorrowers();
    model.value = true;
  });

  async function submit() {
    if (!formData.borrowerId || !formData.dueAt) {
      toast.error("Please select a borrower and due date");
      return;
    }

    loading.value = true;
    const payload: LoanCreate = {
      itemId: props.itemId,
      borrowerId: formData.borrowerId,
      dueAt: new Date(formData.dueAt).toISOString(),
      notes: formData.notes || undefined,
      quantity: formData.quantity,
    };

    const { error } = await api.loans.create(payload);
    loading.value = false;

    if (error) {
      toast.error("Failed to check out item");
      return;
    }

    toast.success(`"${props.itemName}" checked out successfully`);
    closeDialog(DialogID.CreateLoan);
    model.value = false;
    emit("created");
  }

  function handleClose() {
    closeDialog(DialogID.CreateLoan);
    model.value = false;
  }
</script>

<template>
  <Dialog :dialog-id="DialogID.CreateLoan" v-model:open="model" @update:open="val => !val && handleClose()">
    <DialogContent class="sm:max-w-[450px]">
      <DialogHeader>
        <DialogTitle>Check Out Equipment</DialogTitle>
        <DialogDescription>
          Loan "{{ itemName }}" to a borrower.
        </DialogDescription>
      </DialogHeader>

      <form @submit.prevent="submit" class="space-y-4">
        <div class="grid gap-4">
          <div class="grid gap-2">
            <Label for="borrower">Borrower *</Label>
            <Select v-model="formData.borrowerId">
              <SelectTrigger>
                <SelectValue placeholder="Select a borrower..." />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="borrower in borrowers"
                  :key="borrower.id"
                  :value="borrower.id"
                >
                  {{ borrower.name }}
                  <span v-if="borrower.organization" class="text-muted-foreground">
                    ({{ borrower.organization }})
                  </span>
                </SelectItem>
              </SelectContent>
            </Select>
            <p v-if="borrowers.length === 0 && !loadingBorrowers" class="text-sm text-muted-foreground">
              No borrowers found. 
              <NuxtLink to="/borrowers" class="underline">Add a borrower</NuxtLink> first.
            </p>
          </div>

          <div class="grid gap-2">
            <Label for="dueAt">
              <MdiCalendar class="mr-1 inline size-4" />
              Due Date *
            </Label>
            <Input
              id="dueAt"
              v-model="formData.dueAt"
              type="date"
              required
            />
          </div>

          <div v-if="itemQuantity > 1" class="grid gap-2">
            <Label for="quantity">Quantity (max {{ itemQuantity }})</Label>
            <Input
              id="quantity"
              v-model.number="formData.quantity"
              type="number"
              min="1"
              :max="itemQuantity"
            />
          </div>

          <div class="grid gap-2">
            <Label for="notes">Notes</Label>
            <Textarea
              id="notes"
              v-model="formData.notes"
              placeholder="Optional checkout notes..."
              rows="2"
            />
          </div>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleClose">
            Cancel
          </Button>
          <Button type="submit" :disabled="loading || borrowers.length === 0">
            <MdiLoading v-if="loading" class="mr-2 animate-spin" />
            Check Out
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
