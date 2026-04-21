<script setup lang="ts">
import { ref } from 'vue'
import type { EditorSelection } from '../types'

const props = withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string
    readonly?: boolean
  }>(),
  {
    placeholder: 'Start writing notes, extracts, or refined summary here…',
    readonly: false,
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'selection-change': [selection: EditorSelection | null]
}>()

const textareaRef = ref<HTMLTextAreaElement | null>(null)

function onInput(e: Event) {
  if (props.readonly) {
    return
  }
  const target = e.target as HTMLTextAreaElement
  emit('update:modelValue', target.value)
  emitSelection(target)
}

function emitSelection(el: HTMLTextAreaElement | null) {
  if (!el) {
    emit('selection-change', null)
    return
  }
  const from = el.selectionStart ?? 0
  const to = el.selectionEnd ?? 0
  if (to <= from) {
    emit('selection-change', null)
    return
  }
  const text = el.value.slice(from, to)
  if (!text.trim()) {
    emit('selection-change', null)
    return
  }
  emit('selection-change', {
    from,
    to,
    text,
  })
}

function onSelect() {
  emitSelection(textareaRef.value)
}
</script>

<template>
  <div class="tiptap-editor" :class="{ 'tiptap-editor--readonly': readonly }">
    <textarea
      ref="textareaRef"
      class="tiptap-editor__textarea"
      :value="modelValue"
      :readonly="readonly"
      :placeholder="placeholder"
      @input="onInput"
      @select="onSelect"
      @mouseup="onSelect"
      @keyup="onSelect"
    />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.tiptap-editor {
  flex: 1 1 auto;
  min-height: 0;
  background: $color-bg-pane;
}

.tiptap-editor__textarea {
  width: 100%;
  height: 100%;
  min-height: 260px;
  resize: none;
  border: 0;
  outline: none;
  padding: $space-lg;
  background: transparent;
  color: $color-text;
  font-family: $font-ui;
  font-size: $font-size-sm;
  line-height: 1.58;
}

.tiptap-editor__textarea::placeholder {
  color: $color-text-secondary;
}

.tiptap-editor--readonly .tiptap-editor__textarea {
  color: $color-text-secondary;
}
</style>
