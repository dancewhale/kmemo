<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useSettingsStore } from '@/modules/settings/stores/settings.store'
import ShortcutHint from './ShortcutHint.vue'
import CommandList from './CommandList.vue'
import { useCommandStore } from '../stores/command.store'

const commandStore = useCommandStore()
const settings = useSettingsStore()
const { isOpen, query, matchedCommands, activeIndex } = storeToRefs(commandStore)
const inputRef = ref<HTMLInputElement | null>(null)

watch(isOpen, async (open) => {
  if (!open) {
    return
  }
  await nextTick()
  inputRef.value?.focus()
  inputRef.value?.select()
})

function onInput(event: Event) {
  const target = event.target as HTMLInputElement
  commandStore.setQuery(target.value)
}

function onOverlayClick(event: MouseEvent) {
  if (event.target === event.currentTarget) {
    commandStore.close()
  }
}

function onSelect(index: number) {
  commandStore.setActiveIndex(index)
}

async function onExecute(index: number) {
  commandStore.setActiveIndex(index)
  await commandStore.executeActiveCommand()
}
</script>

<template>
  <Teleport to="body">
    <div v-if="isOpen" class="command-palette__overlay" @mousedown="onOverlayClick">
      <section class="command-palette" role="dialog" aria-modal="true" aria-label="Command palette">
        <div class="command-palette__input-wrap">
          <input
            ref="inputRef"
            :value="query"
            type="text"
            class="command-palette__input"
            placeholder="Type a command..."
            @input="onInput"
          />
        </div>
        <CommandList
          :items="matchedCommands"
          :active-index="activeIndex"
          @select="onSelect"
          @execute="onExecute"
        />
        <footer v-if="settings.workspacePreferences.showShortcutHints" class="command-palette__footer">
          <span class="command-palette__hint">
            <ShortcutHint :keys="['Enter']" />
            Run
          </span>
          <span class="command-palette__hint">
            <ShortcutHint :keys="['↑', '↓']" />
            Navigate
          </span>
          <span class="command-palette__hint">
            <ShortcutHint :keys="['Esc']" />
            Close
          </span>
        </footer>
      </section>
    </div>
  </Teleport>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.command-palette__overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 12vh;
  background: rgb(0 0 0 / 40%);
  z-index: 1400;
}

.command-palette {
  width: min(760px, 92vw);
  border-radius: $radius-lg;
  border: 1px solid $color-border;
  background: $color-bg-pane-elevated;
  box-shadow: $shadow-overlay;
  overflow: hidden;
}

.command-palette__input-wrap {
  padding: $space-sm;
  border-bottom: 1px solid $color-border-subtle;
}

.command-palette__input {
  width: 100%;
  border: 1px solid $color-border;
  border-radius: $radius-sm;
  background: $color-bg-pane;
  color: $color-text;
  padding: 9px 12px;
  font-size: $font-size-sm;
  line-height: $line-normal;
  outline: none;
}

.command-palette__input:focus {
  border-color: $color-active-border;
}

.command-palette__footer {
  display: flex;
  align-items: center;
  gap: $space-md;
  padding: $space-sm $space-md;
  border-top: 1px solid $color-border-subtle;
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.command-palette__hint {
  display: inline-flex;
  align-items: center;
  gap: $space-xs;
}
</style>
