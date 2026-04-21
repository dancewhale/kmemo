<script setup lang="ts">
import { storeToRefs } from 'pinia'
import AppActionBar from '@/shared/components/AppActionBar.vue'
import AppSectionTitle from '@/shared/components/AppSectionTitle.vue'
import CaptureQuickEntry from '@/modules/inbox-capture/components/CaptureQuickEntry.vue'
import { useReaderStore } from '../stores/reader.store'
import type { ReaderStatusFilter } from '../types'

const reader = useReaderStore()
const { filter } = storeToRefs(reader)

const statusOptions: { value: ReaderStatusFilter; label: string }[] = [
  { value: 'all', label: 'All' },
  { value: 'inbox', label: 'Inbox' },
  { value: 'reading', label: 'Reading' },
  { value: 'processed', label: 'Processed' },
]
</script>

<template>
  <AppActionBar class="reader-toolbar">
    <div class="reader-toolbar__left">
      <AppSectionTitle title="Reading Queue" />
      <CaptureQuickEntry :compact="true" class="reader-toolbar__capture" />
    </div>
    <div class="reader-toolbar__row">
      <el-input
        :model-value="filter.keyword"
        size="small"
        clearable
        placeholder="Search title or summary…"
        class="reader-toolbar__search"
        @update:model-value="reader.setKeyword(String($event ?? ''))"
      />
      <el-select
        :model-value="filter.status"
        size="small"
        class="reader-toolbar__status"
        @update:model-value="reader.setStatusFilter($event as ReaderStatusFilter)"
      >
        <el-option v-for="o in statusOptions" :key="o.value" :label="o.label" :value="o.value" />
      </el-select>
      <button type="button" class="reader-toolbar__clear" @click="reader.clearFilters()">Clear</button>
    </div>
  </AppActionBar>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.reader-toolbar {
  justify-content: space-between;
}

.reader-toolbar__left {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: $space-sm;
}

.reader-toolbar__capture {
  flex-shrink: 0;
}

.reader-toolbar__row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: $space-sm;
}

.reader-toolbar__search {
  flex: 1 1 160px;
  min-width: 120px;
}

.reader-toolbar__status {
  width: 112px;
  flex-shrink: 0;
}

.reader-toolbar__clear {
  flex-shrink: 0;
  font-size: $font-size-xs;
  padding: 4px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text-secondary;
  background: $color-bg-pane-elevated;
}

.reader-toolbar__clear:hover {
  border-color: $color-active-border;
  color: $color-text;
}
</style>
