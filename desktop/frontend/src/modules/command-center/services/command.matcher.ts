import type { CommandItem, CommandMatchResult } from '../types'

function include(haystack: string | undefined, needle: string): boolean {
  if (!haystack) {
    return false
  }
  return haystack.toLowerCase().includes(needle)
}

export function matchCommands(commands: CommandItem[], query: string): CommandMatchResult[] {
  const q = query.trim().toLowerCase()
  if (!q) {
    return commands.map((c) => ({ ...c, score: 0 }))
  }

  const out: CommandMatchResult[] = []
  for (const c of commands) {
    let score = -1
    if (include(c.title, q)) {
      score = Math.max(score, 3)
    }
    if (include(c.subtitle, q)) {
      score = Math.max(score, 2)
    }
    if ((c.keywords ?? []).some((kw) => kw.toLowerCase().includes(q))) {
      score = Math.max(score, 1)
    }
    if (score >= 0) {
      out.push({ ...c, score })
    }
  }

  return out.sort((a, b) => b.score - a.score || a.title.localeCompare(b.title))
}
