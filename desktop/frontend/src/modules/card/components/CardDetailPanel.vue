<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import { useCardStore } from '../stores/card.store'
import { useReviewStore } from '@/modules/review/stores/review.store'
import CardEditorForm from './CardEditorForm.vue'
import CardMetaSection from './CardMetaSection.vue'
import CardReviewSection from './CardReviewSection.vue'

const card = useCardStore()
const review = useReviewStore()
const { selectedCard, saving } = storeToRefs(card)

const titleDraft = ref('')
const promptDraft = ref('')
const answerDraft = ref('')

watch(
  selectedCard,
  (c) => {
    titleDraft.value = c?.title ?? ''
    promptDraft.value = c?.prompt ?? ''
    answerDraft.value = c?.answer ?? ''
  },
  { immediate: true },
)

const dirty = computed(() => {
  if (!selectedCard.value) {
    return false
  }
  return (
    titleDraft.value.trim() !== selectedCard.value.title.trim() ||
    promptDraft.value.trim() !== selectedCard.value.prompt.trim() ||
    answerDraft.value.trim() !== selectedCard.value.answer.trim()
  )
})

const hasReview = computed(() => {
  const c = selectedCard.value
  if (!c) {
    return false
  }
  if (c.reviewItemId) {
    return true
  }
  return Boolean(review.getReviewItemByCardId(c.id, c.nodeId))
})

async function save() {
  if (!selectedCard.value || !dirty.value) {
    return
  }
  await card.updateCard(selectedCard.value.id, {
    title: titleDraft.value.trim(),
    prompt: promptDraft.value.trim(),
    answer: answerDraft.value.trim(),
  })
}
</script>

<template>
  <div class="card-detail-panel">
    <AppEmpty v-if="!selectedCard" message="Select a card node to edit" />
    <template v-else>
      <header class="card-detail-panel__head">
        <h2 class="card-detail-panel__title">Card</h2>
        <span class="card-detail-panel__status" :class="{ 'card-detail-panel__status--dirty': dirty }">
          {{ saving ? 'Saving…' : dirty ? 'Edited · Unsaved' : 'Saved' }}
        </span>
      </header>

      <CardEditorForm
        :card="{
          ...selectedCard,
          title: titleDraft,
          prompt: promptDraft,
          answer: answerDraft,
        }"
        :disabled="saving"
        @update:title="titleDraft = $event"
        @update:prompt="promptDraft = $event"
        @update:answer="answerDraft = $event"
      />

      <CardMetaSection :card="selectedCard" :has-review="hasReview" />

      <CardReviewSection :has-review="hasReview" />

      <footer class="card-detail-panel__footer">
        <button type="button" class="card-detail-panel__btn" :disabled="!dirty || saving" @click="save">Save</button>
      </footer>
    </template>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.card-detail-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  height: 100%;
  padding: $space-md;
  overflow: auto;
  background: $color-bg-pane;
}

.card-detail-panel__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $space-sm;
  margin-bottom: $space-sm;
}

.card-detail-panel__title {
  margin: 0;
  font-size: $font-size-lg;
  font-weight: 600;
}

.card-detail-panel__status {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.card-detail-panel__status--dirty {
  color: #d4c89a;
}

.card-detail-panel__footer {
  margin-top: $space-lg;
  display: flex;
  justify-content: flex-end;
}

.card-detail-panel__btn {
  font-size: $font-size-xs;
  padding: 3px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-active-border;
  color: $color-text;
  background: $color-bg-pane-elevated;
}

.card-detail-panel__btn:disabled {
  opacity: 0.5;
}
</style>
