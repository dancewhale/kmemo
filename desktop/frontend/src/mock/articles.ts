import type { ArticleSourceType, ArticleStatus } from '@/shared/types/common'

export interface MockArticle {
  id: string
  title: string
  summary: string
  sourceType: ArticleSourceType
  status: ArticleStatus
  updatedAt: string
}

export const mockArticles: MockArticle[] = [
  {
    id: 'art-01',
    title: 'FSRS scheduling notes',
    summary: 'Retention targets, review intervals, and workload smoothing for incremental reading.',
    sourceType: 'web',
    status: 'active',
    updatedAt: '2026-04-18T09:12:00Z',
  },
  {
    id: 'art-02',
    title: 'Progressive reading workflow',
    summary: 'Extract highlights, defer details, and promote stable knowledge into the tree.',
    sourceType: 'rss',
    status: 'pending',
    updatedAt: '2026-04-17T14:40:00Z',
  },
  {
    id: 'art-03',
    title: 'SQLite + GORM patterns',
    summary: 'Migrations, DAO generation, and local-first storage constraints for desktop tools.',
    sourceType: 'file',
    status: 'active',
    updatedAt: '2026-04-16T11:05:00Z',
  },
  {
    id: 'art-04',
    title: 'gRPC worker boundaries',
    summary: 'Keep algorithms in Python, persist structured results in Go, version proto contracts.',
    sourceType: 'manual',
    status: 'pending',
    updatedAt: '2026-04-15T20:22:00Z',
  },
  {
    id: 'art-05',
    title: 'HTML cleaning pipeline',
    summary: 'Sanitize imports, strip boilerplate, preserve semantic blocks for reading mode.',
    sourceType: 'web',
    status: 'active',
    updatedAt: '2026-04-14T08:55:00Z',
  },
  {
    id: 'art-06',
    title: 'Desktop UI density',
    summary: 'Compact lists, split panes, and status surfaces without admin-table aesthetics.',
    sourceType: 'manual',
    status: 'archived',
    updatedAt: '2026-04-12T17:30:00Z',
  },
  {
    id: 'art-07',
    title: 'Knowledge graph vs tree',
    summary: 'When to keep strict hierarchy vs loose linking for recall-oriented study.',
    sourceType: 'rss',
    status: 'pending',
    updatedAt: '2026-04-11T10:18:00Z',
  },
  {
    id: 'art-08',
    title: 'Import material preparation',
    summary: 'Normalize sources, chunk long pages, attach provenance for later review.',
    sourceType: 'file',
    status: 'active',
    updatedAt: '2026-04-10T07:44:00Z',
  },
]
