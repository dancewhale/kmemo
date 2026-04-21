import type { ReviewStatus } from '@/shared/types/common'

export interface MockReviewItem {
  id: string
  title: string
  prompt: string
  dueAt: string
  status: ReviewStatus
}

export const mockReviewItems: MockReviewItem[] = [
  {
    id: 'rev-01',
    title: 'FSRS: retention meaning',
    prompt: 'What does the retention parameter control in daily workload?',
    dueAt: '2026-04-21T08:00:00Z',
    status: 'due',
  },
  {
    id: 'rev-02',
    title: 'Proto versioning',
    prompt: 'Why keep algorithm output contracts versioned across Go/Python?',
    dueAt: '2026-04-21T12:30:00Z',
    status: 'due',
  },
  {
    id: 'rev-03',
    title: 'Incremental reading rule',
    prompt: 'When should a passage be deferred instead of finished immediately?',
    dueAt: '2026-04-22T09:00:00Z',
    status: 'learned',
  },
  {
    id: 'rev-04',
    title: 'SQLite locking',
    prompt: 'What is WAL mode useful for in a desktop reader app?',
    dueAt: '2026-04-23T07:15:00Z',
    status: 'due',
  },
  {
    id: 'rev-05',
    title: 'HTML pipeline',
    prompt: 'Name two risks if import HTML is not sanitized before display.',
    dueAt: '2026-04-24T18:00:00Z',
    status: 'suspended',
  },
]
