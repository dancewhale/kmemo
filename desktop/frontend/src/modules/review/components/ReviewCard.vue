<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useEditorStore } from '@/modules/editor/stores/editor.store'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useCardStore } from '@/modules/card/stores/card.store'
import { useReviewStore } from '../stores/review.store'
import type { ReviewGrade } from '../types'
import ReviewActions from './ReviewActions.vue'
import ReviewMetaSection from './ReviewMetaSection.vue'
import AppSectionTitle from '@/shared/components/AppSectionTitle.vue'
import AppStatusTag from '@/shared/components/AppStatusTag.vue'

const review = useReviewStore()
const workspace = useWorkspaceStore()
const reader = useReaderStore()
const editor = useEditorStore()
const extract = useExtractStore()
const tree = useTreeStore()
const card = useCardStore()

const { selectedItem, isQueueEmpty, submitting } = storeToRefs(review)
const showAnswer = ref(false)

watch(
  () => selectedItem.value?.id,
  () => {
    showAnswer.value = false
  },
)

const canOpenSourceArticle = computed(() => Boolean(selectedItem.value?.sourceArticleId))
const canOpenExtract = computed(() => Boolean(selectedItem.value?.extractId))
const canOpenCard = computed(
  () =>
    selectedItem.value?.type === 'card' &&
    Boolean(selectedItem.value.cardId || selectedItem.value.nodeId),
)
const canRevealTree = computed(() => Boolean(selectedItem.value?.nodeId))

async function onGrade(grade: ReviewGrade) {
  await review.submitGrade(grade)
  showAnswer.value = false
}

async function openSourceArticle() {
  const sourceId = selectedItem.value?.sourceArticleId
  if (!sourceId) {
    return
  }
  await reader.initialize()
  reader.openArticleById(sourceId)
  workspace.setContext('reading')
  const article = reader.getArticleById(sourceId)
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
}

async function openCard() {
  const item = selectedItem.value
  if (!item || item.type !== 'card') {
    return
  }
  await card.initialize()
  await tree.initialize()
  if (item.nodeId) {
    card.openCardByNodeId(item.nodeId)
  } else if (item.cardId) {
    card.setSelectedCard(item.cardId)
    const c = card.getCardById(item.cardId)
    if (c?.nodeId) {
      tree.setSelectedNode(c.nodeId)
    }
  }
  workspace.setContext('knowledge')
}

async function openExtract() {
  const extractId = selectedItem.value?.extractId
  if (!extractId) {
    return
  }
  await extract.initialize()
  extract.setSelectedExtract(extractId)
  const item = extract.getExtractById(extractId)
  if (item?.treeNodeId) {
    tree.setSelectedNode(item.treeNodeId)
  } else if (selectedItem.value?.nodeId) {
    tree.setSelectedNode(selectedItem.value.nodeId)
  }
  workspace.setContext('knowledge')
}

function revealInTree() {
  if (!selectedItem.value?.nodeId) {
    return
  }
  tree.setSelectedNode(selectedItem.value.nodeId)
  workspace.setContext('knowledge')
}
</script>

<template>
  <div class="review-card">
    <AppEmpty v-if="isQueueEmpty" message="All reviews completed for now." />
    <AppEmpty v-else-if="!selectedItem" message="Select a review item from the queue." />
    <template v-else>
      <header class="review-card__head">
        <h2 class="review-card__title">{{ selectedItem.title }}</h2>
        <AppStatusTag kind="info">{{ selectedItem.type }} · {{ selectedItem.status }}</AppStatusTag>
      </header>

      <section class="review-card__section">
        <AppSectionTitle title="Prompt" />
        <p class="review-card__prompt">{{ selectedItem.prompt }}</p>
      </section>

      <section class="review-card__section">
        <button
          v-if="!showAnswer"
          type="button"
          class="review-card__show-btn"
          @click="showAnswer = true"
        >
          Show Answer
        </button>
        <template v-else>
          <AppSectionTitle title="Answer" />
          <p class="review-card__answer">{{ selectedItem.answer }}</p>
        </template>
      </section>

      <section class="review-card__section">
        <AppSectionTitle title="Source" />
        <div class="review-card__source kmono">
          <span>{{ selectedItem.sourceArticleTitle ?? '—' }}</span>
          <span v-if="selectedItem.sourceArticleId">({{ selectedItem.sourceArticleId }})</span>
        </div>
        <div class="review-card__source-actions">
          <button
            type="button"
            class="review-card__link-btn"
            :disabled="!canOpenSourceArticle"
            @click="openSourceArticle"
          >
            Open Source Article
          </button>
          <button
            type="button"
            class="review-card__link-btn"
            :disabled="!canOpenExtract"
            @click="openExtract"
          >
            Open Extract
          </button>
          <button
            type="button"
            class="review-card__link-btn"
            :disabled="!canOpenCard"
            @click="openCard"
          >
            Open Card
          </button>
          <button
            type="button"
            class="review-card__link-btn"
            :disabled="!canRevealTree"
            @click="revealInTree"
          >
            Reveal in Tree
          </button>
        </div>
      </section>

      <section class="review-card__section">
        <ReviewMetaSection :item="selectedItem" />
      </section>

      <section class="review-card__section">
        <ReviewActions :enabled="showAnswer" :submitting="submitting" @grade="onGrade" />
      </section>
    </template>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.review-card {
  display: flex;
  flex-direction: column;
  min-height: 0;
  height: 100%;
  overflow: auto;
  padding: $space-md;
  background: $color-bg-pane;
}

.review-card__head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: $space-sm;
}

.review-card__title {
  margin: 0;
  font-size: $font-size-lg;
  font-weight: 600;
}

.review-card__section {
  margin-top: $space-lg;
}

.review-card__prompt,
.review-card__answer {
  margin: 0;
  padding: $space-sm;
  border: 1px solid $color-border-subtle;
  border-radius: $radius-sm;
  background: $color-bg-pane-elevated;
  font-size: $font-size-sm;
  line-height: $line-normal;
}

.review-card__show-btn {
  font-size: $font-size-xs;
  padding: 3px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-active-border;
  color: $color-text;
  background: $color-bg-pane-elevated;
}

.review-card__source {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.review-card__source-actions {
  margin-top: $space-sm;
  display: flex;
  flex-wrap: wrap;
  gap: $space-xs;
}

.review-card__link-btn {
  font-size: $font-size-xs;
  padding: 2px $space-sm;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text-secondary;
  background: $color-bg-pane-elevated;
}

.review-card__link-btn:disabled {
  opacity: 0.45;
}
</style>
