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

const { startVertical, startHorizontal } = usePaneResize({
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
    ]"
    @pointerdown="onPointerDown"
  />
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-splitter {
  flex-shrink: 0;
  z-index: 2;
  background: transparent;
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
.app-splitter:active {
  background: rgba(107, 155, 209, 0.25);
}
</style>
