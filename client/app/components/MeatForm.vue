<script setup lang="ts">
interface MeatDraft {
  name: string
  cut: string
  species: string
  weightG: number
  pricePerKg: number
  inStock: boolean
}

const model = defineModel<MeatDraft>({ required: true })

defineProps<{
  loading?: boolean
  error?: string
  submitLabel: string
}>()

const emit = defineEmits<{
  submit: []
  cancel: []
}>()
</script>

<template>
  <form
    class="space-y-4"
    @submit.prevent="emit('submit')"
  >
    <UAlert
      v-if="error"
      color="error"
      variant="subtle"
      icon="i-lucide-circle-alert"
      :description="error"
    />

    <UFormField
      label="Name"
      required
    >
      <UInput
        v-model="model.name"
        placeholder="Ribeye"
        class="w-full"
        required
      />
    </UFormField>

    <div class="grid grid-cols-2 gap-4">
      <UFormField
        label="Species"
        required
      >
        <UInput
          v-model="model.species"
          placeholder="Beef"
          class="w-full"
          required
        />
      </UFormField>

      <UFormField
        label="Cut"
        required
      >
        <UInput
          v-model="model.cut"
          placeholder="Steak"
          class="w-full"
          required
        />
      </UFormField>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <UFormField
        label="Weight (g)"
        required
      >
        <UInputNumber
          v-model="model.weightG"
          :min="0"
          class="w-full"
        />
      </UFormField>

      <UFormField
        label="Price / kg"
        required
      >
        <UInputNumber
          v-model="model.pricePerKg"
          :min="0"
          :step="0.01"
          class="w-full"
        />
      </UFormField>
    </div>

    <UFormField label="In stock">
      <USwitch v-model="model.inStock" />
    </UFormField>

    <div class="flex justify-end gap-2 pt-2">
      <UButton
        label="Cancel"
        color="neutral"
        variant="subtle"
        :disabled="loading"
        @click="emit('cancel')"
      />
      <UButton
        :label="submitLabel"
        type="submit"
        :loading="loading"
      />
    </div>
  </form>
</template>
