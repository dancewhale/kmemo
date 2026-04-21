export function readLocalStorageJSON<T>(key: string): T | null {
  try {
    const raw = localStorage.getItem(key)
    if (raw == null) {
      return null
    }
    return JSON.parse(raw) as T
  } catch {
    return null
  }
}

export function writeLocalStorageJSON(key: string, value: unknown): void {
  try {
    localStorage.setItem(key, JSON.stringify(value))
  } catch {
    /* ignore quota / private mode */
  }
}
