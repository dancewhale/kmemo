<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { UI_Z_INDEX } from '@/shared/constants/ui'
import { clampMenuPosition } from '@/shared/utils/dom'
import { useContextMenuStore } from '../stores/context-menu.store'
import ContextMenuList from './ContextMenuList.vue'

const store = useContextMenuStore()
const { isOpen, x, y, menuItems } = storeToRefs(store)

const rootRef = ref<HTMLElement | null>(null)
const pos = ref({ x: 0, y: 0 })

function syncPosition() {
  const el = rootRef.value
  if (!el) {
    pos.value = { x: x.value, y: y.value }
    return
  }
  const rect = el.getBoundingClientRect()
  pos.value = clampMenuPosition(x.value, y.value, rect.width, rect.height)
}

function onDocPointerDown(ev: MouseEvent | PointerEvent) {
  const el = rootRef.value
  if (!el || !isOpen.value) {
    return
  }
  const t = ev.target as Node
  if (!el.contains(t)) {
    store.close()
  }
}

function onKeyDown(ev: KeyboardEvent) {
  if (ev.key === 'Escape' && isOpen.value) {
    ev.preventDefault()
    ev.stopPropagation()
    store.close()
  }
}

watch(
  () => [isOpen.value, x.value, y.value] as const,
  async ([open]) => {
    if (!open) {
      return
    }
    await nextTick()
    syncPosition()
    requestAnimationFrame(syncPosition)
  },
)

watch(isOpen, (open) => {
  if (open) {
    document.addEventListener('pointerdown', onDocPointerDown, true)
    document.addEventListener('keydown', onKeyDown, true)
  } else {
    document.removeEventListener('pointerdown', onDocPointerDown, true)
    document.removeEventListener('keydown', onKeyDown, true)
  }
})

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', onDocPointerDown, true)
  document.removeEventListener('keydown', onKeyDown, true)
})

function onSelect(id: string) {
  void store.execute(id)
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="isOpen && menuItems.length"
      ref="rootRef"
      class="context-menu"
      :style="{ left: `${pos.x}px`, top: `${pos.y}px`, zIndex: UI_Z_INDEX.contextMenu }"
      role="presentation"
    >
      <ContextMenuList :items="menuItems" @select="onSelect" />
    </div>
  </Teleport>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.context-menu {
  position: fixed;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  background: $color-bg-pane-elevated;
  box-shadow: 0 8px 28px rgba(0, 0, 0, 0.35);
}
</style>
