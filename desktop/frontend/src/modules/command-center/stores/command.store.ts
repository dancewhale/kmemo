import { defineStore } from 'pinia'
import { COMMAND_IDS } from '@/shared/constants/commands'
import { COMMAND_REGISTRY } from '../services/command.registry'
import { executeCommand, isCommandEnabled } from '../services/command.executor'
import { matchCommands } from '../services/command.matcher'
import type { CommandItem, CommandMatchResult, CommandState } from '../types'

const INITIAL_STATE: CommandState = {
  isOpen: false,
  query: '',
  activeIndex: 0,
}

function clampActive(index: number, length: number): number {
  if (length <= 0) {
    return 0
  }
  if (index < 0) {
    return 0
  }
  if (index > length - 1) {
    return length - 1
  }
  return index
}

export const useCommandStore = defineStore('command-center', {
  state: (): CommandState => ({ ...INITIAL_STATE }),
  getters: {
    commands(): CommandItem[] {
      return COMMAND_REGISTRY.map((command) => ({
        ...command,
        enabled: isCommandEnabled(command.id),
      }))
    },
    matchedCommands(): CommandMatchResult[] {
      return matchCommands(this.commands, this.query)
    },
    activeCommand(): CommandMatchResult | null {
      return this.matchedCommands[this.activeIndex] ?? null
    },
    hasResults(): boolean {
      return this.matchedCommands.length > 0
    },
  },
  actions: {
    open() {
      this.isOpen = true
      this.query = ''
      this.activeIndex = 0
    },
    close() {
      this.isOpen = false
      this.query = ''
      this.activeIndex = 0
    },
    toggle() {
      if (this.isOpen) {
        this.close()
      } else {
        this.open()
      }
    },
    setQuery(query: string) {
      this.query = query
      this.activeIndex = clampActive(this.activeIndex, this.matchedCommands.length)
    },
    moveUp() {
      if (!this.matchedCommands.length) {
        this.activeIndex = 0
        return
      }
      this.activeIndex =
        this.activeIndex <= 0 ? this.matchedCommands.length - 1 : this.activeIndex - 1
    },
    moveDown() {
      if (!this.matchedCommands.length) {
        this.activeIndex = 0
        return
      }
      this.activeIndex =
        this.activeIndex >= this.matchedCommands.length - 1 ? 0 : this.activeIndex + 1
    },
    setActiveIndex(index: number) {
      this.activeIndex = clampActive(index, this.matchedCommands.length)
    },
    async executeActiveCommand() {
      const active = this.activeCommand
      if (!active || active.enabled === false) {
        return false
      }
      const executed = await executeCommand(active.id)
      if (executed && active.id !== COMMAND_IDS.toggleCommandPalette) {
        this.close()
      }
      return executed
    },
    async executeById(commandId: CommandItem['id']) {
      if (commandId === COMMAND_IDS.toggleCommandPalette) {
        this.toggle()
        return true
      }
      if (commandId === COMMAND_IDS.closeCommandPalette) {
        this.close()
        return true
      }
      const executed = await executeCommand(commandId)
      if (executed) {
        this.close()
      }
      return executed
    },
    resetState() {
      this.isOpen = INITIAL_STATE.isOpen
      this.query = INITIAL_STATE.query
      this.activeIndex = INITIAL_STATE.activeIndex
    },
  },
})
