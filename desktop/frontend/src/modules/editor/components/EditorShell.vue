<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import { buildPendingExtractFromSelection } from '@/modules/extract/services/extract.mapper'
import { useSettingsStore } from '@/modules/settings/stores/settings.store'
import { useToast } from '@/shared/composables/useToast'
import { useEditorStore } from '../stores/editor.store'
import EditorToolbar from './EditorToolbar.vue'
import TiptapEditor from './TiptapEditor.vue'
import ExtractToolbar from './ExtractToolbar.vue'
import ExtractCreateDialog from './ExtractCreateDialog.vue'
import ArticleExtractList from '@/modules/extract/components/ArticleExtractList.vue'

const editor = useEditorStore()
const extractStore = useExtractStore()
const treeStore = useTreeStore()
const settings = useSettingsStore()
const toast = useToast()
const { currentDocument, currentContent, dirty, saving, saveStatusText, selection } =
  storeToRefs(editor)

const updatedText = computed(() => {
  const ts = currentDocument.value?.updatedAt
  if (!ts) {
    return '—'
  }
  try {
    return new Date(ts).toLocaleString()
  } catch {
    return ts
  }
})

const hasSelection = computed(() => {
  return Boolean(selection.value?.text && selection.value.text.trim().length > 0)
})

const contentWidthClass = computed(() => settings.contentWidthClass)
const fontSizeClass = computed(() => settings.editorFontSizeClass)
const showRelatedExtracts = computed(() => settings.readingPreferences.showRelatedExtracts)

function onCreateExtract() {
  if (!currentDocument.value || !selection.value) {
    return
  }
  const payload = buildPendingExtractFromSelection({
    sourceArticleId: currentDocument.value.id,
    sourceArticleTitle: currentDocument.value.title,
    selection: selection.value,
    defaultParentNodeId: treeStore.selectedNodeId ?? 'kn-reading',
  })
  extractStore.openCreateDialog(payload)
}

async function onSave() {
  if (!dirty.value || saving.value) {
    return
  }
  await editor.markSaved()
  toast.success('Document saved')
}
</script>

<template>
  <div class="editor-shell" :class="[contentWidthClass, fontSizeClass]">
    <AppEmpty v-if="!currentDocument" message="Select a document to open editor shell" />
    <template v-else>
      <header class="editor-shell__head">
        <div class="editor-shell__title-wrap">
          <h2 class="editor-shell__title">{{ currentDocument.title }}</h2>
          <span class="editor-shell__status" :class="{ 'editor-shell__status--dirty': dirty }">
            {{ saveStatusText }}
          </span>
        </div>
        <div class="editor-shell__meta kmono">
          <span>{{ currentDocument.contentType }}</span>
          <span class="editor-shell__dot">·</span>
          <span>updated {{ updatedText }}</span>
          <span v-if="currentDocument.sourceUrl" class="editor-shell__dot">·</span>
          <span v-if="currentDocument.sourceUrl">source ready</span>
        </div>
        <div class="editor-shell__actions">
          <button
            type="button"
            class="editor-shell__btn"
            :disabled="!dirty || saving"
            @click="editor.resetContent()"
          >
            Reset
          </button>
          <button
            type="button"
            class="editor-shell__btn editor-shell__btn--primary"
            :disabled="!dirty || saving"
            @click="onSave"
          >
            {{ saving ? 'Saving…' : 'Save' }}
          </button>
        </div>
      </header>
      <EditorToolbar />
      <div class="editor-shell__editor-wrap">
        <ExtractToolbar :visible="hasSelection" @create-extract="onCreateExtract" />
        <TiptapEditor
          :model-value="currentContent"
          placeholder="Draft cleaned article text, extracted note, or synthesis paragraph…"
          @update:model-value="editor.setContent($event)"
          @selection-change="editor.setSelection($event)"
        />
      </div>
      <ArticleExtractList
        v-if="showRelatedExtracts && currentDocument.contentType === 'article'"
        :article-id="currentDocument.id"
      />
    </template>
    <ExtractCreateDialog />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.editor-shell {
  display: flex;
  flex-direction: column;
  min-height: 0;
  height: 100%;
  background: $color-bg-pane;
}

.editor-shell__head {
  padding: $space-md;
  border-bottom: 1px solid $color-border-subtle;
}

.editor-shell__title-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $space-sm;
}

.editor-shell__title {
  margin: 0;
  font-size: $font-size-lg;
  font-weight: 600;
  line-height: $line-tight;
}

.editor-shell__status {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.editor-shell__status--dirty {
  color: #d4c89a;
}

.editor-shell__meta {
  margin-top: $space-xs;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  color: $color-text-secondary;
  font-size: 10px;
}

.editor-shell__dot {
  opacity: 0.6;
}

.editor-shell__actions {
  margin-top: $space-sm;
  display: flex;
  gap: $space-xs;
}

.editor-shell__btn {
  font-size: $font-size-xs;
  padding: 2px $space-sm;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text-secondary;
  background: $color-bg-pane-elevated;
}

.editor-shell__btn--primary {
  border-color: $color-active-border;
  color: $color-text;
}

.editor-shell__btn:disabled {
  opacity: 0.5;
}

.editor-shell__editor-wrap {
  position: relative;
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
  justify-content: center;
}

.editor-shell.pref-content-width-narrow .editor-shell__editor-wrap :deep(.tiptap-editor) {
  width: min(760px, 100%);
}

.editor-shell.pref-content-width-medium .editor-shell__editor-wrap :deep(.tiptap-editor) {
  width: min(900px, 100%);
}

.editor-shell.pref-content-width-wide .editor-shell__editor-wrap :deep(.tiptap-editor) {
  width: min(1080px, 100%);
}

.editor-shell.pref-editor-font-small .editor-shell__editor-wrap :deep(.tiptap-editor__textarea) {
  font-size: 12px;
}

.editor-shell.pref-editor-font-medium .editor-shell__editor-wrap :deep(.tiptap-editor__textarea) {
  font-size: 13px;
}

.editor-shell.pref-editor-font-large .editor-shell__editor-wrap :deep(.tiptap-editor__textarea) {
  font-size: 15px;
}
</style>
