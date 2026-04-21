<script setup lang="ts">
import { computed, ref, watch } from 'vue'

const props = defineProps<{
  modelValue: string
  focusToken?: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  submit: []
  clear: []
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const canClear = computed(() => props.modelValue.trim().length > 0)

watch(
  () => props.focusToken,
  () => {
    inputRef.value?.focus()
    inputRef.value?.select()
  },
)

function onInput(event: Event) {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}

function onClear() {
  emit('clear')
}

function focusInput() {
  inputRef.value?.focus()
}

defineExpose({
  focusInput,
})
</script>

<template>
  <div class="search-toolbar">
    <input
      ref="inputRef"
      type="text"
      :value="props.modelValue"
      class="search-toolbar__input"
      placeholder="Search articles, extracts, nodes, reviews..."
      @input="onInput"
      @keydown.enter="emit('submit')"
    />
    <button v-if="canClear" type="button" class="search-toolbar__clear" @click="onClear">Clear</button>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.search-toolbar {
  display: flex;
  align-items: center;
  gap: $space-sm;
}

.search-toolbar__input {
  flex: 1 1 auto;
  min-width: 0;
  padding: 9px 12px;
  border-radius: $radius-sm;
  border: 1px solid $color-border;
  background: $color-bg-pane;
  color: $color-text;
  font-size: $font-size-sm;
  outline: none;
}

.search-toolbar__input:focus {
  border-color: $color-active-border;
}

.search-toolbar__clear {
  flex-shrink: 0;
  padding: 6px 10px;
  border-radius: $radius-sm;
  border: 1px solid $color-border;
  background: $color-bg-pane-elevated;
  color: $color-text-secondary;
  font-size: $font-size-xs;
}

.search-toolbar__clear:hover {
  border-color: $color-active-border;
  color: $color-text;
}
</style>
