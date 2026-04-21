/**
 * Basic HTTP(S) URL check for capture forms (not a full WHATWG validator).
 */
export function isValidHttpUrl(raw: string): boolean {
  const t = raw.trim()
  if (!t) {
    return false
  }
  const candidate = /^https?:\/\//i.test(t) ? t : `https://${t}`
  try {
    const u = new URL(candidate)
    return u.protocol === 'http:' || u.protocol === 'https:'
  } catch {
    return false
  }
}

/** Normalize user URL for storage (add https if scheme missing). */
export function normalizeCaptureUrl(raw: string): string {
  const t = raw.trim()
  if (!t) {
    return t
  }
  return /^https?:\/\//i.test(t) ? t : `https://${t}`
}
