<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowDown,
  Collection,
  Document,
  Notebook,
  Reading,
  Search,
} from '@element-plus/icons-vue'
import { ROUTE_PATHS } from '@/shared/constants/routes'

const route = useRoute()
const router = useRouter()

const items = computed(() => [
  { key: 'inbox', label: 'Inbox', to: ROUTE_PATHS.inbox, icon: Document },
  { key: 'reading', label: 'Reading', to: ROUTE_PATHS.reading, icon: Reading },
  { key: 'knowledge', label: 'Knowledge', to: ROUTE_PATHS.knowledge, icon: Collection },
  { key: 'review', label: 'Review', to: ROUTE_PATHS.review, icon: Notebook },
  { key: 'search', label: 'Search', to: ROUTE_PATHS.search, icon: Search },
])

const currentLabel = computed(() => {
  const hit = items.value.find((it) => route.path === it.to)
  return hit?.label ?? 'Workspace'
})

function isActive(path: string) {
  return route.path === path
}

function go(to: string) {
  void router.push(to)
}
</script>

<template>
  <el-dropdown trigger="click" class="workspace-nav-menu" popper-class="workspace-nav-menu__popper">
    <el-button type="default" class="workspace-nav-menu__trigger">
      <span class="workspace-nav-menu__trigger-label">{{ currentLabel }}</span>
      <el-icon class="workspace-nav-menu__trigger-chev"><ArrowDown /></el-icon>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="it in items"
          :key="it.key"
          :class="{ 'workspace-nav-menu__dropdown-item--active': isActive(it.to) }"
          @click="go(it.to)"
        >
          <span class="workspace-nav-menu__item">
            <el-icon class="workspace-nav-menu__ico"><component :is="it.icon" /></el-icon>
            <span>{{ it.label }}</span>
          </span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.workspace-nav-menu__trigger {
  display: inline-flex;
  align-items: center;
  gap: $space-xs;
  font-weight: 600;
}

.workspace-nav-menu__trigger-label {
  min-width: 4.5rem;
  text-align: left;
}

.workspace-nav-menu__trigger-chev {
  font-size: 12px;
  color: $color-text-secondary;
}

.workspace-nav-menu__item {
  display: inline-flex;
  align-items: center;
  gap: $space-sm;
}

.workspace-nav-menu__ico {
  font-size: 16px;
}
</style>

<style lang="scss">
@use '@/app/styles/variables.scss' as *;

.workspace-nav-menu__popper .workspace-nav-menu__dropdown-item--active {
  color: $color-text;
  background: $color-active;
  font-weight: 600;
}
</style>
