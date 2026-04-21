<script setup lang="ts">
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useReaderStore } from '../stores/reader.store'
import ReaderToolbar from './ReaderToolbar.vue'
import ReaderListItem from './ReaderListItem.vue'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import AppLoading from '@/shared/components/AppLoading.vue'

const reader = useReaderStore()
const { filteredArticles, selectedArticleId, loading } = storeToRefs(reader)

onMounted(async () => {
  await reader.initialize()
})
</script>

<template>
  <div class="reader-list">
    <ReaderToolbar />
    <div class="reader-list__body">
      <div v-if="loading" class="reader-list__loading">
        <AppLoading label="Loading queue…" />
      </div>
      <AppEmpty v-else-if="!filteredArticles.length" message="No articles match the current filters" />
      <div v-else class="reader-list__items">
        <ReaderListItem
          v-for="a in filteredArticles"
          :key="a.id"
          :article="a"
          :selected="selectedArticleId === a.id"
          @select="reader.setSelectedArticle($event)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.reader-list {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: $color-bg-pane;
}

.reader-list__body {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  padding: $space-xs 0 $space-sm;
}

.reader-list__loading {
  display: flex;
  justify-content: center;
  padding: $space-xl 0;
}

.reader-list__items {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.reader-list__body :deep(.app-empty) {
  padding: $space-xl $space-md;
}
</style>
