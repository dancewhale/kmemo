export type CaptureType = 'blank' | 'paste-text' | 'url'

export interface BlankArticleInput {
  title: string
  content: string
  tags?: string[]
}

export interface PasteTextCaptureInput {
  title: string
  content: string
  tags?: string[]
}

export interface UrlCaptureInput {
  url: string
  title?: string
  note?: string
  tags?: string[]
}

export interface CaptureDraftState {
  type: CaptureType
  blank: BlankArticleInput
  pasteText: PasteTextCaptureInput
  url: UrlCaptureInput
}

export interface CaptureSubmissionResult {
  articleId: string
}

export interface CaptureValidationResult {
  valid: boolean
  message?: string
}
