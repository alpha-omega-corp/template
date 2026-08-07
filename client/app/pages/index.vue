<script setup lang="ts">
interface Meat {
  id: number
  name: string
  cut: string
  species: string
  weightG: number
  pricePerKg: number
  inStock: boolean
  imagePath: string
  createdAt: string
  updatedAt: string
}

const config = useRuntimeConfig()

const { data: meats, status, error, refresh } = await useFetch<Meat[]>('/meat', {
  baseURL: config.public.apiBase
})

function imageUrl(imagePath: string) {
  return `${config.public.apiBase}/storage/${imagePath}`
}

const search = ref('')

function emptyMeat() {
  return {
    name: '',
    cut: '',
    species: '',
    weightG: 0,
    pricePerKg: 0,
    inStock: true
  }
}

const isAddOpen = ref(false)
const isCreating = ref(false)
const createError = ref('')
const newMeat = ref(emptyMeat())

async function onCreate() {
  isCreating.value = true
  createError.value = ''

  try {
    await $fetch('/meat', {
      method: 'POST',
      baseURL: config.public.apiBase,
      body: newMeat.value
    })

    await refresh()
    newMeat.value = emptyMeat()
    isAddOpen.value = false
  } catch (err) {
    createError.value = err instanceof Error ? err.message : 'Failed to create meat'
  } finally {
    isCreating.value = false
  }
}

function onCancelCreate() {
  isAddOpen.value = false
  newMeat.value = emptyMeat()
}

const isEditOpen = ref(false)
const isEditing = ref(false)
const editError = ref('')
const editMeat = ref<Meat | null>(null)

function openEdit(meat: Meat) {
  editMeat.value = { ...meat }
  editError.value = ''
  isEditOpen.value = true
}

async function onEditSubmit() {
  if (!editMeat.value) return

  isEditing.value = true
  editError.value = ''

  try {
    await $fetch(`/meat/${editMeat.value.id}`, {
      method: 'PUT',
      baseURL: config.public.apiBase,
      body: editMeat.value
    })

    await refresh()
    isEditOpen.value = false
  } catch (err) {
    editError.value = err instanceof Error ? err.message : 'Failed to update meat'
  } finally {
    isEditing.value = false
  }
}

function onCancelEdit() {
  isEditOpen.value = false
}

const isDeleteOpen = ref(false)
const isDeleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref<Meat | null>(null)

function openDelete(meat: Meat) {
  deleteTarget.value = meat
  deleteError.value = ''
  isDeleteOpen.value = true
}

async function onDeleteConfirm() {
  if (!deleteTarget.value) return

  isDeleting.value = true
  deleteError.value = ''

  try {
    await $fetch(`/meat/${deleteTarget.value.id}`, {
      method: 'DELETE',
      baseURL: config.public.apiBase
    })

    await refresh()
    isDeleteOpen.value = false
  } catch (err) {
    deleteError.value = err instanceof Error ? err.message : 'Failed to delete meat'
  } finally {
    isDeleting.value = false
  }
}

function onCancelDelete() {
  isDeleteOpen.value = false
}

const filteredMeats = computed(() => {
  const term = search.value.trim().toLowerCase()
  if (!term) return meats.value ?? []

  return (meats.value ?? []).filter(meat =>
    meat.name.toLowerCase().includes(term)
    || meat.cut.toLowerCase().includes(term)
    || meat.species.toLowerCase().includes(term)
  )
})

const columns = [
  { accessorKey: 'imagePath', header: '' },
  { accessorKey: 'name', header: 'Name' },
  { accessorKey: 'species', header: 'Species' },
  { accessorKey: 'cut', header: 'Cut' },
  { accessorKey: 'weightG', header: 'Weight' },
  { accessorKey: 'pricePerKg', header: 'Price / kg' },
  { accessorKey: 'inStock', header: 'Stock' },
  { accessorKey: 'actions', header: '' }
]

const numberFormat = new Intl.NumberFormat('en-US', { maximumFractionDigits: 0 })
const priceFormat = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })
</script>

<template>
  <UContainer class="py-8">
    <div class="flex items-center justify-between gap-4 mb-6">
      <div>
        <h1 class="text-2xl font-bold text-highlighted">
          Meats
        </h1>
        <p class="text-muted">
          Everything currently in the butcher's inventory.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <UInput
          v-model="search"
          icon="i-lucide-search"
          placeholder="Search by name, species or cut..."
          class="w-72"
        />

        <UModal
          v-model:open="isAddOpen"
          title="New meat"
          description="Add a new item to the inventory."
        >
          <UButton
            icon="i-lucide-plus"
            label="Add meat"
          />

          <template #body>
            <MeatForm
              v-model="newMeat"
              submit-label="Create"
              :loading="isCreating"
              :error="createError"
              @submit="onCreate"
              @cancel="onCancelCreate"
            />
          </template>
        </UModal>
      </div>
    </div>

    <UModal
      v-model:open="isEditOpen"
      title="Edit meat"
      description="Update this item's details."
    >
      <template #body>
        <MeatForm
          v-if="editMeat"
          v-model="editMeat"
          submit-label="Save"
          :loading="isEditing"
          :error="editError"
          @submit="onEditSubmit"
          @cancel="onCancelEdit"
        />
      </template>
    </UModal>

    <UModal
      v-model:open="isDeleteOpen"
      title="Delete meat"
      :description="`Remove ${deleteTarget?.name ?? 'this item'} from the inventory. This cannot be undone.`"
    >
      <template #body>
        <div class="space-y-4">
          <UAlert
            v-if="deleteError"
            color="error"
            variant="subtle"
            icon="i-lucide-circle-alert"
            :description="deleteError"
          />

          <div class="flex justify-end gap-2">
            <UButton
              label="Cancel"
              color="neutral"
              variant="subtle"
              :disabled="isDeleting"
              @click="onCancelDelete"
            />
            <UButton
              label="Delete"
              color="error"
              :loading="isDeleting"
              @click="onDeleteConfirm"
            />
          </div>
        </div>
      </template>
    </UModal>

    <UAlert
      v-if="error"
      color="error"
      variant="subtle"
      icon="i-lucide-circle-alert"
      title="Failed to load meats"
      :description="error.message"
      class="mb-6"
    >
      <template #actions>
        <UButton
          color="error"
          variant="subtle"
          @click="refresh()"
        >
          Retry
        </UButton>
      </template>
    </UAlert>

    <UTable
      :data="filteredMeats"
      :columns="columns"
      :loading="status === 'pending'"
    >
      <template #imagePath-cell="{ row }">
        <img
          v-if="row.original.imagePath"
          :src="imageUrl(row.original.imagePath)"
          :alt="row.original.name"
          class="size-12 rounded-md object-cover ring ring-default"
        >
        <div
          v-else
          class="size-12 rounded-md bg-elevated flex items-center justify-center"
        >
          <UIcon
            name="i-lucide-image-off"
            class="text-dimmed size-5"
          />
        </div>
      </template>

      <template #weightG-cell="{ row }">
        {{ numberFormat.format(row.original.weightG) }} g
      </template>

      <template #pricePerKg-cell="{ row }">
        {{ priceFormat.format(row.original.pricePerKg) }}
      </template>

      <template #inStock-cell="{ row }">
        <UBadge
          :color="row.original.inStock ? 'success' : 'neutral'"
          variant="subtle"
        >
          {{ row.original.inStock ? 'In stock' : 'Out of stock' }}
        </UBadge>
      </template>

      <template #actions-cell="{ row }">
        <div class="flex justify-end gap-1">
          <UButton
            icon="i-lucide-pencil"
            color="neutral"
            variant="ghost"
            size="sm"
            aria-label="Edit meat"
            @click="openEdit(row.original)"
          />
          <UButton
            icon="i-lucide-trash-2"
            color="error"
            variant="ghost"
            size="sm"
            aria-label="Delete meat"
            @click="openDelete(row.original)"
          />
        </div>
      </template>

      <template #empty>
        <div class="py-8 text-center text-muted">
          No meats found.
        </div>
      </template>
    </UTable>
  </UContainer>
</template>
