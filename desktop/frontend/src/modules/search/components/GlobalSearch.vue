<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import SearchFilterTabs from './SearchFilterTabs.vue'
import SearchResultList from './SearchResultList.vue'
import SearchToolbar from './SearchToolbar.vue'
import { useSearchStore } from '../stores/search.store'
import type { SearchFilterType, SearchResultItem } from '../types'

const search = useSearchStore()
const { query, filter, results, selectedResultId, hasQuery, loading, resultCount, focusToken } =
  storeToRefs(search)

const toolbarRef = ref<InstanceType<typeof SearchToolbar> | null>(null)

onMounted(async () => {
  await search.initialize()
  toolbarRef.value?.focusInput()
})

function onQueryChange(value: string) {
  search.setQuery(value)
}

function onFilterChange(value: SearchFilterType) {
  search.setFilter(value)
}

function onSelect(id: string) {
  search.setSelectedResult(id)
}

async function onOpen(item: SearchResultItem) {
  await search.openResult(item)
}

async function onSubmit() {
  if (!search.selectedResult) {
    return
  }
  await search.openResult(search.selectedResult)
}

function onClear() {
  search.clearSearch()
}
</script>

<template>
  <section class="global-search">
    <header class="global-search__header">
      <SearchToolbar
        ref="toolbarRef"
        :model-value="query"
        :focus-token="focusToken"
        @update:model-value="onQueryChange"
        @submit="onSubmit"
        @clear="onClear"
      />
      <div class="global-search__meta">
        <SearchFilterTabs :model-value="filter" @update:model-value="onFilterChange" />
        <span class="global-search__count">
          {{ loading ? 'Indexing…' : `${resultCount} result${resultCount === 1 ? '' : 's'}` }}
        </span>
      </div>
    </header>
    <SearchResultList
      :items="results"
      :selected-id="selectedResultId"
      :has-query="hasQuery"
      @select="onSelect"
      @open="onOpen"
    />
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.global-search {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.global-search__header {
  padding: $space-sm;
  display: grid;
  gap: $space-sm;
}

.global-search__meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $space-sm;
}

.global-search__count {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}
</style>
