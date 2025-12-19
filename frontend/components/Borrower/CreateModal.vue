<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import MdiLoading from "~icons/mdi/loading";
  import { Button } from "@/components/ui/button";
  import { Input } from "@/components/ui/input";
  import { Label } from "@/components/ui/label";
  import { Textarea } from "@/components/ui/textarea";
  import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog";
  import { useDialog } from "@/components/ui/dialog-provider";
  import { DialogID } from "~/components/ui/dialog-provider/utils";
  import type { BorrowerCreate } from "~~/lib/api/types/data-contracts";

  const emit = defineEmits<{
    created: [id: string];
  }>();

  const { t } = useI18n();
  const api = useUserApi();
  const { closeDialog, registerOpenDialogCallback } = useDialog();

  const loading = ref(false);
  const model = ref(false);

  const formData = reactive<BorrowerCreate>({
    name: "",
    email: "",
    phone: "",
    organization: "",
    studentId: "",
    notes: "",
  });

  function resetForm() {
    formData.name = "";
    formData.email = "";
    formData.phone = "";
    formData.organization = "";
    formData.studentId = "";
    formData.notes = "";
  }

  registerOpenDialogCallback(DialogID.CreateBorrower, () => {
    resetForm();
    model.value = true;
  });

  async function submit() {
    if (!formData.name || !formData.email) {
      toast.error("Name and email are required");
      return;
    }

    loading.value = true;
    const { data, error } = await api.borrowers.create(formData);
    loading.value = false;

    if (error) {
      toast.error("Failed to create borrower");
      return;
    }

    toast.success(`Borrower "${formData.name}" created successfully`);
    closeDialog(DialogID.CreateBorrower);
    model.value = false;
    emit("created", data.id);
  }

  function handleClose() {
    closeDialog(DialogID.CreateBorrower);
    model.value = false;
  }
</script>

<template>
  <Dialog :dialog-id="DialogID.CreateBorrower" v-model:open="model" @update:open="val => !val && handleClose()">
    <DialogContent class="sm:max-w-[500px]">
      <DialogHeader>
        <DialogTitle>Add Borrower</DialogTitle>
        <DialogDescription>
          Create a new borrower who can check out equipment.
        </DialogDescription>
      </DialogHeader>

      <form @submit.prevent="submit" class="space-y-4">
        <div class="grid gap-4">
          <div class="grid gap-2">
            <Label for="name">Name *</Label>
            <Input
              id="name"
              v-model="formData.name"
              placeholder="Full name"
              required
            />
          </div>

          <div class="grid gap-2">
            <Label for="email">Email *</Label>
            <Input
              id="email"
              v-model="formData.email"
              type="email"
              placeholder="email@example.com"
              required
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="grid gap-2">
              <Label for="phone">Phone</Label>
              <Input
                id="phone"
                v-model="formData.phone"
                type="tel"
                placeholder="(555) 123-4567"
              />
            </div>

            <div class="grid gap-2">
              <Label for="studentId">Student ID</Label>
              <Input
                id="studentId"
                v-model="formData.studentId"
                placeholder="Optional"
              />
            </div>
          </div>

          <div class="grid gap-2">
            <Label for="organization">Organization / Startup</Label>
            <Input
              id="organization"
              v-model="formData.organization"
              placeholder="Company or startup name"
            />
          </div>

          <div class="grid gap-2">
            <Label for="notes">Notes</Label>
            <Textarea
              id="notes"
              v-model="formData.notes"
              placeholder="Any additional notes about this borrower..."
              rows="3"
            />
          </div>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleClose">
            Cancel
          </Button>
          <Button type="submit" :disabled="loading">
            <MdiLoading v-if="loading" class="mr-2 animate-spin" />
            Create Borrower
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
