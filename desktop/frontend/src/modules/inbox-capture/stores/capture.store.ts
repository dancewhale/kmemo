import { defineStore } from 'pinia'
import { router } from '@/app/router'
import { ROUTE_NAMES } from '@/shared/constants/routes'
import { useToast } from '@/shared/composables/useToast'
import { useEditorStore } from '@/modules/editor/stores/editor.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useSearchStore } from '@/modules/search/stores/search.store'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import {
  buildArticleFromBlankInput,
  buildArticleFromPasteTextInput,
  buildArticleFromUrlInput,
} from '../services/capture.mapper'
import { validateBlankInput, validatePasteTextInput, validateUrlInput } from '../services/capture.validator'
import type { ReaderArticle } from '@/modules/reader/types'
import type {
  BlankArticleInput,
  CaptureDraftState,
  CaptureType,
  PasteTextCaptureInput,
  UrlCaptureInput,
} from '../types'

interface CaptureState {
  isDialogOpen: boolean
  activeType: CaptureType
  submitting: boolean
  draft: CaptureDraftState
}

function newArticleId(): string {
  return `art-cap-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 7)}`
}

function emptyDraft(): CaptureDraftState {
  return {
    type: 'blank',
    blank: { title: '', content: '', tags: [] },
    pasteText: { title: '', content: '', tags: [] },
    url: { url: '', title: '', note: '', tags: [] },
  }
}

function parseTagsLine(line: string): string[] {
  return line
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
}

export const useCaptureStore = defineStore('inbox-capture', {
  state: (): CaptureState => ({
    isDialogOpen: false,
    activeType: 'blank',
    submitting: false,
    draft: emptyDraft(),
  }),

  getters: {
    activeDraft(): BlankArticleInput | PasteTextCaptureInput | UrlCaptureInput {
      switch (this.activeType) {
        case 'blank':
          return this.draft.blank
        case 'paste-text':
          return this.draft.pasteText
        case 'url':
          return this.draft.url
        default:
          return this.draft.blank
      }
    },

    canSubmit(): boolean {
      switch (this.activeType) {
        case 'blank':
          return validateBlankInput(this.draft.blank).valid
        case 'paste-text':
          return validatePasteTextInput(this.draft.pasteText).valid
        case 'url':
          return validateUrlInput(this.draft.url).valid
        default:
          return false
      }
    },

    dialogTitle(): string {
      return 'Add to Inbox'
    },
  },

  actions: {
    openDialog(type: CaptureType = 'blank') {
      this.activeType = type
      this.draft.type = type
      this.isDialogOpen = true
    },

    closeDialog() {
      this.isDialogOpen = false
    },

    setActiveType(type: CaptureType) {
      this.activeType = type
      this.draft.type = type
    },

    updateBlankDraft(payload: Partial<BlankArticleInput>) {
      this.draft.blank = { ...this.draft.blank, ...payload }
    },

    updatePasteTextDraft(payload: Partial<PasteTextCaptureInput>) {
      this.draft.pasteText = { ...this.draft.pasteText, ...payload }
    },

    updateUrlDraft(payload: Partial<UrlCaptureInput>) {
      this.draft.url = { ...this.draft.url, ...payload }
    },

    resetDraft(type?: CaptureType) {
      if (type === 'blank') {
        this.draft.blank = { title: '', content: '', tags: [] }
      } else if (type === 'paste-text') {
        this.draft.pasteText = { title: '', content: '', tags: [] }
      } else if (type === 'url') {
        this.draft.url = { url: '', title: '', note: '', tags: [] }
      } else {
        this.draft = emptyDraft()
        this.draft.type = this.activeType
      }
    },

    async submit() {
      if (this.activeType === 'blank') {
        await this.submitBlank()
      } else if (this.activeType === 'paste-text') {
        await this.submitPasteText()
      } else {
        await this.submitUrl()
      }
    },

    async finalizeArticle(article: ReaderArticle, successMessage: string) {
      const toast = useToast()
      const reader = useReaderStore()
      const editor = useEditorStore()
      const workspace = useWorkspaceStore()
      const search = useSearchStore()

      await reader.initialize()
      reader.addArticle(article)
      reader.clearFilters()
      reader.openArticleById(article.id)
      editor.openDocument({
        id: article.id,
        title: article.title,
        content: article.content,
        contentType: 'article',
        updatedAt: article.updatedAt,
        sourceUrl: article.sourceUrl,
        tags: article.tags,
      })
      workspace.setContext('reading')
      await router.push({ name: ROUTE_NAMES.reading })
      search.refreshIndex()
      if (search.initialized) {
        search.runSearch()
      }
      toast.success(successMessage)
      this.closeDialog()
      this.resetDraft()
    },

    async submitBlank() {
      const v = validateBlankInput(this.draft.blank)
      if (!v.valid) {
        useToast().warning(v.message ?? 'Invalid input')
        return
      }
      if (this.submitting) {
        return
      }
      this.submitting = true
      try {
        const id = newArticleId()
        const article = buildArticleFromBlankInput(this.draft.blank, id)
        await this.finalizeArticle(article, 'Blank article created')
      } finally {
        this.submitting = false
      }
    },

    async submitPasteText() {
      const v = validatePasteTextInput(this.draft.pasteText)
      if (!v.valid) {
        useToast().warning(v.message ?? 'Invalid input')
        return
      }
      if (this.submitting) {
        return
      }
      this.submitting = true
      try {
        const id = newArticleId()
        const article = buildArticleFromPasteTextInput(this.draft.pasteText, id)
        await this.finalizeArticle(article, 'Text capture added to inbox')
      } finally {
        this.submitting = false
      }
    },

    async submitUrl() {
      const v = validateUrlInput(this.draft.url)
      if (!v.valid) {
        useToast().warning(v.message ?? 'Invalid input')
        return
      }
      if (this.submitting) {
        return
      }
      this.submitting = true
      try {
        const id = newArticleId()
        const article = buildArticleFromUrlInput(this.draft.url, id)
        await this.finalizeArticle(article, 'URL capture placeholder created')
      } finally {
        this.submitting = false
      }
    },

    /** Tags from a single comma-separated field (phase 1). */
    applyTagsLine(
      target: CaptureType,
      line: string,
    ) {
      const tags = parseTagsLine(line)
      if (target === 'blank') {
        this.updateBlankDraft({ tags })
      } else if (target === 'paste-text') {
        this.updatePasteTextDraft({ tags })
      } else {
        this.updateUrlDraft({ tags })
      }
    },
  },
})
