<script setup lang="ts">
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useEditorStore } from '@/modules/editor/stores/editor.store'
import { useExtractStore } from '../stores/extract.store'

const props = defineProps<{
  sourceArticleId: string
  sourceArticleTitle: string
}>()

const workspace = useWorkspaceStore()
const reader = useReaderStore()
const editor = useEditorStore()
const extract = useExtractStore()

async function openSourceArticle() {
  await reader.initialize()
  reader.openArticleById(props.sourceArticleId)
  workspace.setContext('reading')
  const article = reader.getArticleById(props.sourceArticleId)
  if (article) {
    editor.openDocument({
      id: article.id,
      title: article.title,
      content: article.content,
      contentType: 'article',
      updatedAt: article.updatedAt,
      sourceUrl: article.sourceUrl,
      tags: article.tags,
    })
  }
  extract.setSelectedExtract(null)
}
</script>

<template>
  <section class="extract-source">
    <h3 class="extract-source__title">Source</h3>
    <div class="extract-source__row">
      <span class="extract-source__name">{{ sourceArticleTitle }}</span>
      <span class="extract-source__id kmono">{{ sourceArticleId }}</span>
    </div>
    <button type="button" class="extract-source__btn" @click="openSourceArticle">
      Open Source Article
    </button>
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.extract-source__title {
  margin: 0 0 $space-sm;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.extract-source__row {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-bottom: $space-sm;
}

.extract-source__name {
  font-size: $font-size-sm;
  color: $color-text;
}

.extract-source__id {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.extract-source__btn {
  font-size: $font-size-xs;
  padding: 2px $space-sm;
  border-radius: $radius-sm;
  border: 1px solid $color-active-border;
  color: $color-text;
  background: $color-bg-pane-elevated;
}
</style>
