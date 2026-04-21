import type {
  CreateCardInput,
  CreateManualExtractInput,
  CreateNodeInput,
  CreateObjectValidationResult,
} from '../types'

export function validateNodeInput(_input: CreateNodeInput): CreateObjectValidationResult {
  return { valid: true }
}

export function validateCardInput(input: CreateCardInput): CreateObjectValidationResult {
  if (!input.prompt.trim()) {
    return { valid: false, message: 'Please enter a prompt' }
  }
  if (!input.answer.trim()) {
    return { valid: false, message: 'Please enter an answer' }
  }
  return { valid: true }
}

export function validateManualExtractInput(input: CreateManualExtractInput): CreateObjectValidationResult {
  if (!input.quote.trim()) {
    return { valid: false, message: 'Please enter a quote' }
  }
  return { valid: true }
}
