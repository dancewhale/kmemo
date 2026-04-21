<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import {
  Collection,
  Document,
  Notebook,
  Reading,
  Search,
  Setting,
} from '@element-plus/icons-vue'
import { useWorkspaceStore } from '../stores/workspace.store'
import AppIconButton from '@/shared/components/AppIconButton.vue'
import { ROUTE_PATHS } from '@/shared/constants/routes'

const route = useRoute()
const router = useRouter()
const store = useWorkspaceStore()

const items = computed(() => [
  { key: 'inbox', label: 'Inbox', to: ROUTE_PATHS.inbox, icon: Document },
  { key: 'reading', label: 'Reading', to: ROUTE_PATHS.reading, icon: Reading },
  { key: 'knowledge', label: 'Knowledge', to: ROUTE_PATHS.knowledge, icon: Collection },
  { key: 'review', label: 'Review', to: ROUTE_PATHS.review, icon: Notebook },
  { key: 'search', label: 'Search', to: ROUTE_PATHS.search, icon: Search },
])

function isActive(path: string) {
  return route.path === path
}

function goSettings() {
  void router.push(ROUTE_PATHS.settings)
}
</script>

<template>
  <aside class="left-sidebar">
    <div class="left-sidebar__brand">kmemo</div>
    <nav class="left-sidebar__nav">
      <RouterLink
        v-for="it in items"
        :key="it.key"
        :to="it.to"
        class="left-sidebar__link"
        :class="{ 'left-sidebar__link--active': isActive(it.to) }"
      >
        <el-icon class="left-sidebar__ico">
          <component :is="it.icon" />
        </el-icon>
        <span v-if="!store.isLeftCollapsed" class="left-sidebar__text">{{ it.label }}</span>
      </RouterLink>
    </nav>
    <div class="left-sidebar__footer">
      <AppIconButton
        :label="store.isLeftCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"
        @click="store.setLeftCollapsed(!store.isLeftCollapsed)"
      >
        <span class="left-sidebar__chev">{{ store.isLeftCollapsed ? '»' : '«' }}</span>
      </AppIconButton>
      <AppIconButton v-if="!store.isLeftCollapsed" label="Settings" @click="goSettings">
        <el-icon><Setting /></el-icon>
      </AppIconButton>
      <button
        v-else
        type="button"
        class="left-sidebar__icon-only"
        title="Settings"
        @click="goSettings"
      >
        <el-icon><Setting /></el-icon>
      </button>
    </div>
  </aside>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.left-sidebar {
  display: flex;
  flex-direction: column;
  border-right: 1px solid $color-border-subtle;
  background: $color-bg-pane;
  min-height: 0;
}

.left-sidebar__brand {
  flex: 0 0 auto;
  padding: $space-md $space-md $space-sm;
  font-size: $font-size-xs;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: $color-text-secondary;
  border-bottom: 1px solid $color-border-subtle;
}

.left-sidebar__nav {
  flex: 1 1 auto;
  overflow: auto;
  padding: $space-sm 0;
}

.left-sidebar__link {
  display: flex;
  align-items: center;
  gap: $space-sm;
  padding: $space-sm $space-md;
  margin: 0 $space-xs;
  border-radius: $radius-sm;
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: $line-tight;
}

.left-sidebar__link:hover {
  background: $color-hover;
  color: $color-text;
}

.left-sidebar__link--active {
  background: $color-active;
  color: $color-text;
  border: 1px solid $color-active-border;
}

.left-sidebar__ico {
  font-size: 16px;
  flex-shrink: 0;
}

.left-sidebar__text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.left-sidebar__footer {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  gap: $space-xs;
  padding: $space-sm;
  border-top: 1px solid $color-border-subtle;
}

.left-sidebar__chev {
  font-size: $font-size-md;
  line-height: 1;
}

.left-sidebar__icon-only {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: $radius-sm;
  color: $color-text-secondary;
}

.left-sidebar__icon-only:hover {
  background: $color-hover;
  color: $color-text;
}
</style>
