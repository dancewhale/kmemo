<script setup lang="ts">
const props = defineProps<{
  hasQuery: boolean
}>()
</script>

<template>
  <div class="search-empty">
    <p class="search-empty__title">
      {{ props.hasQuery ? 'No matching objects' : 'Start searching your knowledge workspace' }}
    </p>
    <p class="search-empty__desc">
      {{
        props.hasQuery
          ? 'Try another keyword or switch the filter type.'
          : 'Search articles, extracts, knowledge nodes, and review items.'
      }}
    </p>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.search-empty {
  height: 100%;
  min-height: 180px;
  display: grid;
  place-content: center;
  gap: $space-xs;
  text-align: center;
  color: $color-text-secondary;
}

.search-empty__title {
  margin: 0;
  color: $color-text;
  font-size: $font-size-sm;
  font-weight: 600;
}

.search-empty__desc {
  margin: 0;
  font-size: $font-size-xs;
  line-height: $line-normal;
}
</style>
