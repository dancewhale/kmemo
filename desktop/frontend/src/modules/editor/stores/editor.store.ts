import { defineStore } from 'pinia'
import type { EditorDocument, EditorSelection } from '../types'

interface EditorState {
  currentDocument: EditorDocument | null
  originalContent: string
  currentContent: string
  selection: EditorSelection | null
  dirty: boolean
  saving: boolean
  lastSavedAt: string | null
}

function formatTime(iso: string | null): string {
  if (!iso) {
    return ''
  }
  try {
    return new Date(iso).toLocaleTimeString(undefined, {
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return iso
  }
}

export const useEditorStore = defineStore('editor', {
  state: (): EditorState => ({
    currentDocument: null,
    originalContent: '',
    currentContent: '',
    selection: null,
    dirty: false,
    saving: false,
    lastSavedAt: null,
  }),

  getters: {
    hasDocument(): boolean {
      return this.currentDocument != null
    },

    documentTitle(): string {
      return this.currentDocument?.title ?? 'Untitled'
    },

    documentType(): string {
      return this.currentDocument?.contentType ?? 'article'
    },

    saveStatusText(): string {
      if (this.saving) {
        return 'Saving…'
      }
      if (!this.currentDocument) {
        return 'Idle'
      }
      if (this.dirty) {
        return 'Edited · Unsaved'
      }
      if (this.lastSavedAt) {
        return `Saved at ${formatTime(this.lastSavedAt)}`
      }
      return 'Idle'
    },
  },

  actions: {
    openDocument(doc: EditorDocument) {
      if (this.currentDocument?.id === doc.id && this.dirty) {
        return
      }
      this.currentDocument = { ...doc }
      this.originalContent = doc.content
      this.currentContent = doc.content
      this.selection = null
      this.dirty = false
      this.saving = false
      this.lastSavedAt = doc.updatedAt
    },

    setContent(content: string) {
      this.currentContent = content
      this.dirty = this.currentContent !== this.originalContent
    },

    setSelection(selection: EditorSelection | null) {
      this.selection = selection
    },

    clearSelection() {
      this.selection = null
    },

    async markSaved() {
      if (!this.currentDocument || this.saving) {
        return
      }
      this.saving = true
      await new Promise((resolve) => window.setTimeout(resolve, 360))
      this.originalContent = this.currentContent
      this.dirty = false
      this.saving = false
      const now = new Date().toISOString()
      this.lastSavedAt = now
      this.currentDocument = {
        ...this.currentDocument,
        content: this.currentContent,
        updatedAt: now,
      }
    },

    resetContent() {
      this.currentContent = this.originalContent
      this.dirty = false
    },

    clearDocument() {
      this.currentDocument = null
      this.originalContent = ''
      this.currentContent = ''
      this.selection = null
      this.dirty = false
      this.saving = false
      this.lastSavedAt = null
    },
  },
})
