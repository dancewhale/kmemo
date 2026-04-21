import type { ExtractItem } from '@/modules/extract/types'

export const mockExtractItems: ExtractItem[] = [
  {
    id: 'ext-seed-01',
    sourceArticleId: 'art-ir-basics',
    sourceArticleTitle: 'Incremental Reading Basics',
    title: 'Extract: Outline-first pass',
    quote:
      'A practical flow starts with triage: skim the structure, identify one useful insight, and defer the rest.',
    note: 'Use this as the default first pass rule for long articles.',
    parentNodeId: 'kn-reading-art-basics',
    treeNodeId: 'kn-reading-extract-outline',
    reviewItemId: 'rev-03',
    createdAt: '2026-04-01T16:20:00Z',
    updatedAt: '2026-04-16T19:40:00Z',
  },
  {
    id: 'ext-seed-02',
    sourceArticleId: 'art-active-recall',
    sourceArticleTitle: 'Active Recall vs Passive Review',
    title: 'Extract: Active recall note',
    quote: 'Active recall feels harder because it exposes uncertainty. That discomfort is useful signal.',
    note: '',
    parentNodeId: 'kn-mem-topic',
    treeNodeId: 'kn-mem-extract',
    createdAt: '2026-03-28T15:45:00Z',
    updatedAt: '2026-04-14T09:30:00Z',
  },
]
