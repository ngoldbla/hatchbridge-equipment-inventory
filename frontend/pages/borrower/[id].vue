<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import MdiAccount from "~icons/mdi/account";
  import MdiEmail from "~icons/mdi/email";
  import MdiPhone from "~icons/mdi/phone";
  import MdiDomain from "~icons/mdi/domain";
  import MdiCardAccountDetails from "~icons/mdi/card-account-details";
  import MdiPencil from "~icons/mdi/pencil";
  import MdiDelete from "~icons/mdi/delete";
  import MdiLoading from "~icons/mdi/loading";
  import { Button } from "@/components/ui/button";
  import { Card } from "@/components/ui/card";
  import { Badge } from "@/components/ui/badge";
  import { Separator } from "@/components/ui/separator";
  import { useConfirm } from "~~/composables/use-confirm";
  import BaseContainer from "@/components/Base/Container.vue";
  import DateTime from "~/components/global/DateTime.vue";
  import type { BorrowerOut, LoanSummary } from "~~/lib/api/types/data-contracts";

  definePageMeta({
    middleware: ["auth"],
  });

  const { t } = useI18n();
  const route = useRoute();
  const router = useRouter();
  const api = useUserApi();
  const confirm = useConfirm();

  const borrowerId = computed(() => route.params.id as string);

  const { data: borrower, refresh } = useAsyncData(
    `borrower-${borrowerId.value}`,
    async () => {
      const { data, error } = await api.borrowers.get(borrowerId.value);
      if (error) {
        toast.error("Failed to load borrower");
        return null;
      }
      return data;
    },
    { watch: [borrowerId] }
  );

  const { data: loans } = useAsyncData(
    `borrower-loans-${borrowerId.value}`,
    async () => {
      const { data, error } = await api.borrowers.getLoans(borrowerId.value);
      if (error) {
        toast.error("Failed to load loan history");
        return [];
      }
      return data;
    },
    { watch: [borrowerId] }
  );

  useHead({
    title: computed(() => borrower.value ? `${borrower.value.name} | Hatchbridge Inventory` : "Borrower"),
  });

  const activeLoans = computed(() => 
    loans.value?.filter((l: LoanSummary) => !l.returnedAt) || []
  );

  const returnedLoans = computed(() =>
    loans.value?.filter((l: LoanSummary) => l.returnedAt) || []
  );

  const deleting = ref(false);

  async function confirmDelete() {
    if (!borrower.value) return;

    const result = await confirm.open(
      `Are you sure you want to delete "${borrower.value.name}"? This action cannot be undone.`
    );

    if (!result.data) return;

    deleting.value = true;
    const { error } = await api.borrowers.delete(borrowerId.value);
    deleting.value = false;

    if (error) {
      toast.error("Failed to delete borrower");
      return;
    }

    toast.success("Borrower deleted");
    router.push("/borrowers");
  }

  function formatDate(date: Date | string | null | undefined) {
    if (!date) return "N/A";
    return new Date(date).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  }
</script>

<template>
  <BaseContainer v-if="borrower">
    <Title>{{ borrower.name }}</Title>

    <Card class="p-4">
      <header class="mb-4">
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div class="flex items-center gap-4">
            <div class="flex size-16 items-center justify-center rounded-full bg-secondary text-secondary-foreground">
              <MdiAccount class="size-8" />
            </div>
            <div>
              <h1 class="text-2xl font-semibold">{{ borrower.name }}</h1>
              <div class="flex items-center gap-2 text-muted-foreground">
                <MdiEmail class="size-4" />
                <a :href="`mailto:${borrower.email}`" class="hover:underline">
                  {{ borrower.email }}
                </a>
              </div>
              <Badge :variant="borrower.isActive ? 'default' : 'secondary'" class="mt-1">
                {{ borrower.isActive ? 'Active' : 'Inactive' }}
              </Badge>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <Button variant="outline" disabled>
              <MdiPencil class="mr-2 size-4" />
              Edit
            </Button>
            <Button variant="destructive" :disabled="deleting" @click="confirmDelete">
              <MdiLoading v-if="deleting" class="mr-2 animate-spin" />
              <MdiDelete v-else class="mr-2 size-4" />
              Delete
            </Button>
          </div>
        </div>
      </header>

      <Separator class="my-4" />

      <!-- Borrower Details -->
      <div class="grid gap-4 sm:grid-cols-2">
        <div v-if="borrower.phone" class="flex items-center gap-2 text-sm">
          <MdiPhone class="size-4 text-muted-foreground" />
          <span>{{ borrower.phone }}</span>
        </div>
        <div v-if="borrower.organization" class="flex items-center gap-2 text-sm">
          <MdiDomain class="size-4 text-muted-foreground" />
          <span>{{ borrower.organization }}</span>
        </div>
        <div v-if="borrower.studentId" class="flex items-center gap-2 text-sm">
          <MdiCardAccountDetails class="size-4 text-muted-foreground" />
          <span>Student ID: {{ borrower.studentId }}</span>
        </div>
        <div class="text-xs text-muted-foreground">
          Created <DateTime :date="borrower.createdAt" />
        </div>
      </div>

      <div v-if="borrower.notes" class="mt-4 rounded-md bg-muted p-3 text-sm">
        {{ borrower.notes }}
      </div>
    </Card>

    <!-- Active Loans Section -->
    <section v-if="activeLoans.length > 0" class="mt-6">
      <h2 class="mb-3 text-lg font-medium">
        Currently Checked Out ({{ activeLoans.length }})
      </h2>
      <div class="space-y-2">
        <NuxtLink 
          v-for="loan in activeLoans" 
          :key="loan.id"
          :to="`/item/${loan.itemId}`"
          class="block"
        >
          <Card class="cursor-pointer p-3 transition-colors hover:bg-accent/50">
            <div class="flex items-center justify-between">
              <div>
                <span class="font-medium">{{ loan.itemName }}</span>
                <span v-if="loan.quantity > 1" class="ml-2 text-sm text-muted-foreground">
                  (Qty: {{ loan.quantity }})
                </span>
              </div>
              <Badge :variant="loan.isOverdue ? 'destructive' : 'secondary'">
                Due: {{ formatDate(loan.dueAt) }}
              </Badge>
            </div>
          </Card>
        </NuxtLink>
      </div>
    </section>

    <!-- Loan History Section -->
    <section v-if="returnedLoans.length > 0" class="mt-6">
      <h2 class="mb-3 text-lg font-medium text-muted-foreground">
        Loan History ({{ returnedLoans.length }})
      </h2>
      <div class="space-y-2">
        <Card 
          v-for="loan in returnedLoans" 
          :key="loan.id"
          class="p-3 opacity-60"
        >
          <div class="flex items-center justify-between">
            <div>
              <span>{{ loan.itemName }}</span>
            </div>
            <div class="text-xs text-muted-foreground">
              {{ formatDate(loan.checkedOutAt) }} - {{ formatDate(loan.returnedAt) }}
            </div>
          </div>
        </Card>
      </div>
    </section>

    <!-- No Loans -->
    <div 
      v-if="!loans || loans.length === 0" 
      class="mt-6 flex flex-col items-center justify-center py-8 text-center"
    >
      <p class="text-muted-foreground">No loan history for this borrower yet.</p>
    </div>
  </BaseContainer>

  <!-- Loading State -->
  <BaseContainer v-else>
    <div class="flex items-center justify-center py-12">
      <MdiLoading class="size-8 animate-spin text-muted-foreground" />
    </div>
  </BaseContainer>
</template>
