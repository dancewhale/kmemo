<template>
  <section class="panel-block">
    <div class="card-list-header">
      <div>
        <h2 class="panel-title">Cards</h2>
        <p class="panel-description">
          {{ subtitle }}
        </p>
      </div>
      <span class="card-total">{{ total }}</span>
    </div>

    <div class="filter-row">
      <input
        :value="keyword"
        class="search-input"
        placeholder="Search cards by title"
        type="text"
        @input="$emit('update:keyword', ($event.target as HTMLInputElement).value)"
      />
    </div>

    <p v-if="loading" class="panel-description">Loading cards...</p>
    <div v-else-if="error" class="list-feedback-block">
      <p class="panel-description">{{ error }}</p>
      <button class="list-retry-button" type="button" @click="$emit('retry')">Retry</button>
    </div>
    <p v-else-if="items.length === 0" class="panel-description">No cards in the current scope.</p>

    <ul v-else class="card-list">
      <li v-for="item in items" :key="item.id">
        <button class="card-item" :class="{ active: item.id === selectedId }" type="button" @click="$emit('select', item.id)">
          <div class="card-item-title">{{ item.title }}</div>
          <div class="card-item-meta">
            <span>{{ item.cardType }}</span>
            <span>{{ item.knowledgeName || item.knowledgeId }}</span>
          </div>
        </button>
      </li>
    </ul>
  </section>
</template>

<script setup lang="ts">
import type { CardDTO } from "../../types/dto";

defineProps<{
  items: CardDTO[];
  total: number;
  selectedId: string | null;
  keyword: string;
  subtitle: string;
  loading: boolean;
  error: string | null;
}>();

defineEmits<{
  (event: "retry"): void;
  (event: "select", id: string): void;
  (event: "update:keyword", value: string): void;
}>();
</script>

<style scoped>
.card-list-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.card-total {
  min-width: 36px;
  height: 36px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: var(--accent-soft);
  color: var(--text-primary);
  font-weight: 600;
}

.filter-row {
  display: flex;
}

.search-input {
  width: 100%;
  border: 1px solid var(--border-color);
  background: rgba(15, 23, 42, 0.65);
  color: var(--text-primary);
  border-radius: 10px;
  padding: 10px 12px;
}

.list-feedback-block {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
}

.list-retry-button {
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: rgba(15, 23, 42, 0.55);
  color: var(--text-primary);
  padding: 8px 12px;
  cursor: pointer;
}

.card-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-item {
  width: 100%;
  text-align: left;
  border: 1px solid transparent;
  border-radius: 12px;
  background: rgba(15, 23, 42, 0.5);
  color: var(--text-primary);
  padding: 12px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-item.active {
  border-color: var(--accent);
  background: var(--accent-soft);
}

.card-item-title {
  font-weight: 700;
}

.card-item-meta {
  display: flex;
  gap: 10px;
  color: var(--text-secondary);
  font-size: 12px;
}
</style>
