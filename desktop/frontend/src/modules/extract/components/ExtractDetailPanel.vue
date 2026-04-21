<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import { useContextMenuStore } from '@/modules/context-actions/stores/context-menu.store'
import { useReviewStore } from '@/modules/review/stores/review.store'
import { useExtractStore } from '../stores/extract.store'
import ExtractMetaSection from './ExtractMetaSection.vue'
import ExtractReviewSection from './ExtractReviewSection.vue'
import ExtractSourceSection from './ExtractSourceSection.vue'

const contextMenu = useContextMenuStore()
const extract = useExtractStore()
const { selectedExtract, saving } = storeToRefs(extract)

const titleDraft = ref('')
const noteDraft = ref('')

onMounted(async () => {
  const review = useReviewStore()
  await Promise.all([extract.initialize(), review.initialize()])
})

watch(
  selectedExtract,
  (item) => {
    titleDraft.value = item?.title ?? ''
    noteDraft.value = item?.note ?? ''
  },
  { immediate: true },
)

const dirty = computed(() => {
  if (!selectedExtract.value) {
    return false
  }
  return (
    titleDraft.value.trim() !== selectedExtract.value.title.trim() ||
    noteDraft.value.trim() !== selectedExtract.value.note.trim()
  )
})

async function save() {
  if (!selectedExtract.value || !dirty.value) {
    return
  }
  await extract.updateExtract(selectedExtract.value.id, {
    title: titleDraft.value.trim(),
    note: noteDraft.value.trim(),
  })
}

function onPanelContextMenu(e: MouseEvent) {
  if (!selectedExtract.value) {
    return
  }
  e.preventDefault()
  contextMenu.open({
    entityType: 'extract',
    entityId: selectedExtract.value.id,
    extractId: selectedExtract.value.id,
    articleId: selectedExtract.value.sourceArticleId,
    x: e.clientX,
    y: e.clientY,
  })
}
</script>

<template>
  <div class="extract-detail" @contextmenu="onPanelContextMenu">
    <AppEmpty v-if="!selectedExtract" message="Select an extract node to view details" />
    <template v-else>
      <header class="extract-detail__head">
        <h2 class="extract-detail__title">Extract Detail</h2>
        <span class="extract-detail__status" :class="{ 'extract-detail__status--dirty': dirty }">
          {{ saving ? 'Saving…' : dirty ? 'Edited · Unsaved' : 'Saved' }}
        </span>
      </header>

      <section class="extract-detail__section">
        <label class="extract-detail__label">Title</label>
        <el-input v-model="titleDraft" size="small" maxlength="96" />
      </section>

      <section class="extract-detail__section">
        <label class="extract-detail__label">Quote</label>
        <div class="extract-detail__quote">“{{ selectedExtract.quote }}”</div>
      </section>

      <section class="extract-detail__section">
        <label class="extract-detail__label">Note</label>
        <el-input
          v-model="noteDraft"
          type="textarea"
          :autosize="{ minRows: 3, maxRows: 6 }"
          placeholder="Add interpretation, context, or follow-up idea…"
        />
      </section>

      <section class="extract-detail__section">
        <ExtractSourceSection
          :source-article-id="selectedExtract.sourceArticleId"
          :source-article-title="selectedExtract.sourceArticleTitle"
        />
      </section>

      <section class="extract-detail__section">
        <ExtractMetaSection :extract="selectedExtract" />
      </section>

      <ExtractReviewSection />

      <footer class="extract-detail__footer">
        <button type="button" class="extract-detail__btn" :disabled="!dirty || saving" @click="save">
          Save
        </button>
      </footer>
    </template>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.extract-detail {
  display: flex;
  flex-direction: column;
  min-height: 0;
  height: 100%;
  padding: $space-md;
  overflow: auto;
  background: $color-bg-pane;
}

.extract-detail__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $space-sm;
  margin-bottom: $space-sm;
}

.extract-detail__title {
  margin: 0;
  font-size: $font-size-lg;
  font-weight: 600;
}

.extract-detail__status {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.extract-detail__status--dirty {
  color: #d4c89a;
}

.extract-detail__section {
  margin-top: $space-md;
}

.extract-detail__label {
  display: block;
  margin-bottom: $space-xs;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.extract-detail__quote {
  padding: $space-sm;
  border: 1px solid $color-border-subtle;
  border-radius: $radius-sm;
  background: $color-bg-pane-elevated;
  font-size: $font-size-sm;
  line-height: $line-normal;
  color: $color-text;
  font-style: italic;
}

.extract-detail__footer {
  margin-top: $space-lg;
  display: flex;
  justify-content: flex-end;
}

.extract-detail__btn {
  font-size: $font-size-xs;
  padding: 3px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-active-border;
  color: $color-text;
  background: $color-bg-pane-elevated;
}

.extract-detail__btn:disabled {
  opacity: 0.5;
}
</style>
