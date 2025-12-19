<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import MdiPlus from "~icons/mdi/plus";
  import MdiMagnify from "~icons/mdi/magnify";
  import MdiAccountGroup from "~icons/mdi/account-group";
  import { Button } from "@/components/ui/button";
  import { Input } from "@/components/ui/input";
  import { Card } from "@/components/ui/card";
  import { Badge } from "@/components/ui/badge";
  import { useDialog } from "@/components/ui/dialog-provider";
  import { DialogID } from "~/components/ui/dialog-provider/utils";
  import BaseContainer from "@/components/Base/Container.vue";
  import BaseSectionHeader from "@/components/Base/SectionHeader.vue";
  import BorrowerCreateModal from "~/components/Borrower/CreateModal.vue";
  import type { BorrowerSummary } from "~~/lib/api/types/data-contracts";

  definePageMeta({
    middleware: ["auth"],
  });

  const { t } = useI18n();

  useHead({
    title: computed(() => `Hatchbridge Inventory | Borrowers`),
  });

  const api = useUserApi();
  const { openDialog } = useDialog();
  const searchQuery = ref("");

  const { data: borrowers, refresh } = useAsyncData("borrowers", async () => {
    const { data, error } = await api.borrowers.getAll();
    if (error) {
      toast.error("Failed to load borrowers");
      return [];
    }
    return data;
  });

  const filteredBorrowers = computed(() => {
    if (!borrowers.value) return [];
    if (!searchQuery.value) return borrowers.value;

    const query = searchQuery.value.toLowerCase();
    return borrowers.value.filter(
      (b: BorrowerSummary) =>
        b.name.toLowerCase().includes(query) ||
        b.email.toLowerCase().includes(query) ||
        b.organization?.toLowerCase().includes(query) ||
        b.studentId?.toLowerCase().includes(query)
    );
  });

  const activeBorrowers = computed(() => 
    filteredBorrowers.value.filter((b: BorrowerSummary) => b.isActive)
  );

  const inactiveBorrowers = computed(() =>
    filteredBorrowers.value.filter((b: BorrowerSummary) => !b.isActive)
  );

  const handleRefresh = () => refresh();
</script>

<template>
  <BaseContainer>
    <div class="mb-4 flex flex-wrap items-center justify-between gap-4">
      <BaseSectionHeader>
        <div class="flex items-center gap-2">
          <MdiAccountGroup class="size-7" />
          Borrowers
        </div>
      </BaseSectionHeader>
      <div class="flex items-center gap-2">
        <div class="relative">
          <MdiMagnify class="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            v-model="searchQuery"
            placeholder="Search borrowers..."
            class="w-64 pl-9"
          />
        </div>
        <Button @click="openDialog(DialogID.CreateBorrower)">
          <MdiPlus class="mr-2" />
          Add Borrower
        </Button>
      </div>
    </div>

    <!-- Active Borrowers -->
    <div v-if="activeBorrowers.length > 0" class="mb-6">
      <h3 class="mb-3 text-sm font-medium text-muted-foreground">
        Active Borrowers ({{ activeBorrowers.length }})
      </h3>
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <NuxtLink
          v-for="borrower in activeBorrowers"
          :key="borrower.id"
          :to="`/borrower/${borrower.id}`"
          class="block"
        >
          <Card class="cursor-pointer p-4 transition-colors hover:bg-accent/50">
            <div class="flex items-start justify-between">
              <div>
                <h4 class="font-medium">{{ borrower.name }}</h4>
                <p class="text-sm text-muted-foreground">{{ borrower.email }}</p>
                <p v-if="borrower.organization" class="text-sm text-muted-foreground">
                  {{ borrower.organization }}
                </p>
              </div>
              <Badge variant="default">Active</Badge>
            </div>
            <div v-if="borrower.studentId" class="mt-2 text-xs text-muted-foreground">
              Student ID: {{ borrower.studentId }}
            </div>
          </Card>
        </NuxtLink>
      </div>
    </div>

    <!-- Inactive Borrowers -->
    <div v-if="inactiveBorrowers.length > 0" class="mb-6">
      <h3 class="mb-3 text-sm font-medium text-muted-foreground">
        Inactive Borrowers ({{ inactiveBorrowers.length }})
      </h3>
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <NuxtLink
          v-for="borrower in inactiveBorrowers"
          :key="borrower.id"
          :to="`/borrower/${borrower.id}`"
          class="block"
        >
          <Card class="cursor-pointer p-4 opacity-60 transition-colors hover:bg-accent/50">
            <div class="flex items-start justify-between">
              <div>
                <h4 class="font-medium">{{ borrower.name }}</h4>
                <p class="text-sm text-muted-foreground">{{ borrower.email }}</p>
              </div>
              <Badge variant="secondary">Inactive</Badge>
            </div>
          </Card>
        </NuxtLink>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-if="(!borrowers || borrowers.length === 0)"
      class="flex flex-col items-center justify-center py-12 text-center"
    >
      <MdiAccountGroup class="mb-4 size-16 text-muted-foreground" />
      <p class="mb-4 text-muted-foreground">No borrowers found. Add your first borrower to start tracking equipment loans.</p>
      <Button @click="openDialog(DialogID.CreateBorrower)">
        <MdiPlus class="mr-2" />
        Add Borrower
      </Button>
    </div>

    <!-- No Search Results -->
    <div
      v-else-if="filteredBorrowers.length === 0 && searchQuery"
      class="flex flex-col items-center justify-center py-12 text-center"
    >
      <p class="text-muted-foreground">No borrowers match "{{ searchQuery }}"</p>
    </div>

    <!-- Create Modal -->
    <BorrowerCreateModal @created="handleRefresh" />
  </BaseContainer>
</template>
