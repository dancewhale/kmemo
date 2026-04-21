import { isValidHttpUrl } from '@/shared/utils/url'
import type {
  BlankArticleInput,
  CaptureValidationResult,
  PasteTextCaptureInput,
  UrlCaptureInput,
} from '../types'

export function validateBlankInput(_input: BlankArticleInput): CaptureValidationResult {
  return { valid: true }
}

export function validatePasteTextInput(input: PasteTextCaptureInput): CaptureValidationResult {
  if (!input.content.trim()) {
    return { valid: false, message: 'Please enter some text' }
  }
  return { valid: true }
}

export function validateUrlInput(input: UrlCaptureInput): CaptureValidationResult {
  const url = input.url.trim()
  if (!url) {
    return { valid: false, message: 'Please enter a valid URL' }
  }
  if (!isValidHttpUrl(url)) {
    return { valid: false, message: 'Please enter a valid URL' }
  }
  return { valid: true }
}
