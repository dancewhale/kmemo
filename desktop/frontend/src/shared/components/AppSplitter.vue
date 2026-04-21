<script setup lang="ts">
import { usePaneResize } from '@/shared/composables/usePaneResize'

const props = withDefaults(
  defineProps<{
    orientation: 'vertical' | 'horizontal'
  }>(),
  {},
)

const emit = defineEmits<{
  drag: [delta: number]
}>()

const { startVertical, startHorizontal, isDragging } = usePaneResize({
  onMove(delta) {
    emit('drag', delta)
  },
})

function onPointerDown(e: PointerEvent) {
  if (props.orientation === 'vertical') {
    startVertical(e)
  } else {
    startHorizontal(e)
  }
}
</script>

<template>
  <div
    class="app-splitter"
    :class="[
      orientation === 'vertical' ? 'app-splitter--vertical' : 'app-splitter--horizontal',
      { 'app-splitter--dragging': isDragging },
    ]"
    @pointerdown="onPointerDown"
  >
    <span class="app-splitter__line" />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-splitter {
  flex-shrink: 0;
  z-index: 2;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
}

.app-splitter--vertical {
  width: $splitter-hit;
  cursor: col-resize;
}

.app-splitter--horizontal {
  height: $splitter-hit;
  cursor: row-resize;
}

.app-splitter:hover,
.app-splitter--dragging {
  background: $color-hover;
}

.app-splitter__line {
  opacity: 0;
  transition: opacity 0.1s ease, background-color 0.1s ease;
}

.app-splitter--vertical .app-splitter__line {
  width: 1px;
  height: 100%;
}

.app-splitter--horizontal .app-splitter__line {
  width: 100%;
  height: 1px;
}

.app-splitter:hover .app-splitter__line,
.app-splitter--dragging .app-splitter__line {
  opacity: 1;
  background: $color-active-border;
}
</style>
