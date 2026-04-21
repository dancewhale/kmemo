import type { ReaderArticle } from '@/modules/reader/types'

export type { ReaderArticle }

export const mockArticles: ReaderArticle[] = [
  {
    id: 'art-ir-basics',
    title: 'Incremental Reading Basics',
    summary:
      'Defer details, revisit on schedule, and promote durable items into your knowledge tree instead of finishing everything in one pass.',
    content:
      `Incremental reading treats every article as a queue item, not a one-shot obligation. The goal is to preserve momentum while still moving important ideas into long-term memory.\n\nA practical flow starts with triage: skim the structure, identify one useful insight, and defer the rest. On later passes, you reduce friction by turning dense paragraphs into shorter notes, prompts, or explicit questions.\n\nIn a desktop workflow, this means each article should carry metadata, editable content, and links to derived knowledge nodes. The editor shell in this phase is intentionally simple, but it already supports that future loop: open, refine, and save.`,
    sourceType: 'web',
    sourceUrl: 'https://example.com/docs/incremental-reading',
    status: 'reading',
    updatedAt: '2026-04-20T09:30:00Z',
    tags: ['workflow', 'sm'],
  },
  {
    id: 'art-pks',
    title: 'How to Build a Personal Knowledge System',
    summary:
      'Capture, refine, and link notes so retrieval stays cheap: atomic notes, clear titles, and periodic consolidation sessions.',
    content:
      `A personal knowledge system is less about collecting everything and more about reducing retrieval cost. If a note cannot be found, understood, and reused quickly, it is expensive to keep.\n\nStructure follows use. Keep broad topics as stable anchors, attach article extracts under them, and convert recurring ideas into concise cards. This hierarchy keeps both exploration and review practical.\n\nThe most effective systems are iterative: capture today, refine tomorrow, and connect during weekly synthesis. The editor should support this rhythm by making small edits cheap and visible.`,
    sourceType: 'web',
    sourceUrl: 'https://example.com/pks',
    status: 'inbox',
    updatedAt: '2026-04-19T14:12:00Z',
    tags: ['pkm', 'structure'],
  },
  {
    id: 'art-active-recall',
    title: 'Active Recall vs Passive Review',
    summary:
      'Testing strengthens memory more than re-reading; pair prompts with spaced schedules to keep workload predictable.',
    content:
      `Passive review feels productive because it is fluent. Active recall feels harder because it exposes uncertainty. That discomfort is useful signal.\n\nA robust workflow converts claims into prompts. Instead of highlighting “retention improves with spacing,” ask: “Why does spacing reduce relearning cost?” The answer forces retrieval and strengthens future access.\n\nWhen merged with spaced scheduling, recall prompts prevent binge-study cycles. Small, regular checks outperform occasional deep rereads for most durable knowledge work.`,
    sourceType: 'text',
    status: 'processed',
    updatedAt: '2026-04-18T11:05:00Z',
    tags: ['memory', 'fsrs'],
  },
  {
    id: 'art-learning-public',
    title: 'Learning in Public',
    summary:
      'Writing while you learn tightens feedback loops: small public artifacts reduce fear and document decisions over time.',
    content:
      `Learning in public is not performance; it is instrumentation. Publishing small artifacts creates external checkpoints that reveal whether your explanation is coherent.\n\nShort notes, diagrams, and revision logs make progress visible. They also leave breadcrumbs for future you, who often forgets why a decision once felt obvious.\n\nFor a knowledge tool, this suggests a simple requirement: every draft should be easy to revise, timestamp, and export later. The editor shell provides the minimum surface to start that habit.`,
    sourceType: 'web',
    status: 'inbox',
    updatedAt: '2026-04-17T16:40:00Z',
    tags: ['habits'],
  },
  {
    id: 'art-reading-workflows',
    title: 'Designing Reading Workflows',
    summary:
      'Split intake, triage, and deep reading; keep a single queue discipline so the inbox does not become a graveyard.',
    content:
      `Reading workflows break when all material is treated equally. Intake should be cheap, triage should be fast, and deep processing should be reserved for high-value pieces.\n\nA single queue with explicit states keeps decisions reversible. Inbox means “not yet processed,” reading means “actively refining,” processed means “stable enough for downstream use.”\n\nThis stage uses mock statuses, but the pattern is production-ready: list state in the center, editable document in the right pane, and a clear save signal to avoid losing context.`,
    sourceType: 'pdf',
    sourceUrl: 'file:///mock/reading-workflows.pdf',
    status: 'reading',
    updatedAt: '2026-04-16T08:22:00Z',
    tags: ['workflow', 'pdf'],
  },
  {
    id: 'art-atomic-notes',
    title: 'Why Notes Should Be Atomic',
    summary:
      'One idea per note improves reuse and linking; large blobs resist review and hide the smallest reusable insight.',
    content:
      `Atomic notes reduce ambiguity. If a note contains one claim, it can be cited, challenged, and connected without rereading a page of context.\n\nLarge mixed notes often encode multiple assumptions and make spaced review painful. Splitting them does not lose information; it exposes structure that was previously hidden.\n\nIn practice, extraction is iterative: identify a long paragraph, isolate one argument, rewrite in your own words, and keep a pointer to source context.`,
    sourceType: 'text',
    status: 'processed',
    updatedAt: '2026-04-15T19:18:00Z',
    tags: ['notes'],
  },
  {
    id: 'art-prog-summ',
    title: 'Progressive Summarization for Dense Material',
    summary:
      'Layer highlights into bolded phrases, then into your own words, so each pass adds compression without losing meaning.',
    content:
      `Progressive summarization is a compression strategy. Early passes mark candidate sections; later passes reduce them into claims and operational rules.\n\nThe key is that each pass should be cheap. If summarization requires perfect understanding in one session, it collapses under real workload.\n\nA good editor can preserve this layering by keeping raw text, refined summary, and extracted prompts in the same document lifecycle.`,
    sourceType: 'web',
    status: 'reading',
    updatedAt: '2026-04-14T13:50:00Z',
    tags: ['reading', 'compression'],
  },
  {
    id: 'art-wal-sqlite',
    title: 'SQLite WAL for Desktop Readers',
    summary:
      'Write-ahead logging improves concurrent reads during imports; pair with bounded transactions for smooth UI.',
    content:
      `Desktop tools benefit from local persistence, but concurrency still matters. WAL mode allows readers and writers to coexist with fewer visible stalls.\n\nFor user-facing workflows, transaction boundaries are more important than raw throughput. Short writes and predictable checkpoints reduce UI jitter during import or save bursts.\n\nEven in this mock phase, treating edits as explicit save events mirrors the eventual persistence path and keeps state transitions easy to reason about.`,
    sourceType: 'text',
    status: 'inbox',
    updatedAt: '2026-04-13T10:08:00Z',
    tags: ['engineering', 'sqlite'],
  },
  {
    id: 'art-html-pipeline',
    title: 'HTML Cleaning Before Reading Mode',
    summary:
      'Strip chrome, normalize headings, and keep semantic blocks so the reader surface stays calm and typographically consistent.',
    content:
      `Imported HTML usually contains navigation, ads, and layout wrappers that pollute reading context. Cleaning should preserve meaning while removing interaction noise.\n\nA reliable pipeline keeps heading structure, paragraph boundaries, lists, and code blocks. Everything else is optional and can often be deferred.\n\nWhen users edit extracted content later, semantic structure pays off: paragraphs remain coherent and references are easier to preserve.`,
    sourceType: 'web',
    sourceUrl: 'https://example.com/html-pipeline',
    status: 'processed',
    updatedAt: '2026-04-12T07:44:00Z',
    tags: ['import', 'safety'],
  },
  {
    id: 'art-grpc-boundaries',
    title: 'gRPC Worker Boundaries in Hybrid Apps',
    summary:
      'Keep algorithms in a sidecar, version protos, and persist canonical results in the host so the UI stays offline-friendly.',
    content:
      `Hybrid desktop stacks often combine a host runtime with language-specific workers. Clean boundaries matter more than transport details.\n\nVersioned protobuf contracts prevent hidden drift between algorithm outputs and UI expectations. The host should own canonical storage; workers should remain replaceable.\n\nThis architectural split aligns with editor boundaries too: UI edits happen locally, while heavy processing can remain in sidecar services without polluting component logic.`,
    sourceType: 'pdf',
    status: 'inbox',
    updatedAt: '2026-04-11T21:30:00Z',
    tags: ['architecture'],
  },
  {
    id: 'art-density-ui',
    title: 'Desktop UI Density Without Admin Tables',
    summary:
      'Prefer split panes, compact lists, and status rails over wide grids when the product is a personal knowledge cockpit.',
    content:
      `Knowledge work UIs optimize for switching cost, not report generation. Dense panes with stable affordances beat large table layouts for iterative reading and editing.\n\nUsers need quick movement between queue, detail, and action state. A split layout provides that with minimal navigation overhead.\n\nThe editor shell in this phase extends the right pane from passive detail into active workspace while preserving the same dense visual grammar.`,
    sourceType: 'text',
    status: 'reading',
    updatedAt: '2026-04-10T15:02:00Z',
    tags: ['ui'],
  },
  {
    id: 'art-import-prep',
    title: 'Import Material Preparation',
    summary:
      'Normalize sources, chunk long pages, and attach provenance so later review always knows where an excerpt originated.',
    content:
      `Import preparation is a quality gate. If source text enters the system unstructured, every downstream step becomes more expensive.\n\nChunking should respect semantic boundaries, not arbitrary character counts. Each chunk needs provenance metadata so excerpts can always trace back to origin.\n\nIn later phases, this provenance will connect editor extracts to review cards. For now, keeping content editable and save-aware is enough to validate the shell.`,
    sourceType: 'web',
    status: 'processed',
    updatedAt: '2026-04-09T09:15:00Z',
    tags: ['import'],
  },
]
