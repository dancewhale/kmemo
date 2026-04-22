<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useWorkspaceStore } from '../stores/workspace.store'
import AppPane from '@/shared/components/AppPane.vue'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useEditorStore } from '@/modules/editor/stores/editor.store'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import EditorShell from '@/modules/editor/components/EditorShell.vue'
import { useCardStore } from '@/modules/card/stores/card.store'
import ReviewCard from '@/modules/review/components/ReviewCard.vue'
import type { EditorDocument } from '@/modules/editor/types'

const store = useWorkspaceStore()
const tree = useTreeStore()
const reader = useReaderStore()
const editor = useEditorStore()
const extract = useExtractStore()
const card = useCardStore()

function htmlEscape(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

onMounted(() => {
  void reader.initialize()
})

const inboxArticle = computed(() => {
  if (store.currentContext !== 'inbox') {
    return null
  }
  const id = store.selectedArticleId
  if (!id) {
    return null
  }
  return reader.getArticleById(id)
})
const readingArticle = computed(() => (store.currentContext === 'reading' ? reader.selectedArticle : null))
const activeEditorArticle = computed(() => {
  if (store.currentContext === 'reading') {
    return readingArticle.value
  }
  if (store.currentContext === 'inbox') {
    return inboxArticle.value
  }
  return null
})
const kn = computed(() => (store.currentContext === 'knowledge' ? tree.selectedNode : null))
const showExtractDetail = computed(() => {
  return store.currentContext === 'knowledge' && tree.isSelectedExtractNode
})

const documentForEditor = computed((): EditorDocument | null => {
  const ctx = store.currentContext
  if (ctx === 'reading' || ctx === 'inbox') {
    const doc = activeEditorArticle.value
    if (!doc) {
      return null
    }
    return {
      id: doc.id,
      title: doc.title,
      content: doc.content,
      contentType: 'article',
      updatedAt: doc.updatedAt,
      sourceUrl: doc.sourceUrl,
      tags: doc.tags,
    }
  }
  if (ctx !== 'knowledge') {
    return null
  }
  const node = kn.value
  if (!node) {
    return null
  }
  if (showExtractDetail.value && extract.selectedExtract) {
    const e = extract.selectedExtract
    const noteBlock = e.note.trim()
      ? `<h2>Note</h2><p>${htmlEscape(e.note)}</p>`
      : ''
    const content = `<blockquote><p>${htmlEscape(e.quote)}</p></blockquote>${noteBlock}<p class="kmono">${htmlEscape(
      e.sourceArticleTitle,
    )}</p>`
    return {
      id: e.id,
      title: e.title || 'Extract',
      content,
      contentType: 'extract',
      updatedAt: e.updatedAt,
    }
  }
  if (node.type === 'card' && card.selectedCard) {
    const c = card.selectedCard
    const content = `<h2>Prompt</h2><p>${htmlEscape(c.prompt)}</p><h2>Answer</h2><p>${htmlEscape(c.answer)}</p>`
    return {
      id: c.id,
      title: c.title,
      content,
      contentType: 'node',
      updatedAt: c.updatedAt,
    }
  }
  const desc = node.description?.trim() ?? ''
  const content = desc ? `<p>${htmlEscape(desc)}</p>` : '<p><em>No description</em></p>'
  return {
    id: node.id,
    title: node.title,
    content,
    contentType: 'node',
    updatedAt: node.updatedAt ?? new Date().toISOString(),
  }
})

watch(
  () => [store.currentContext, tree.selectedNodeId] as const,
  async ([ctx, nodeId]) => {
    if (ctx !== 'knowledge' || !nodeId) {
      return
    }
    const node = tree.findNodeById(nodeId)
    if (node?.type === 'card') {
      await card.initialize()
      card.openCardByNodeId(node.id)
    }
  },
  { immediate: true },
)

watch(
  documentForEditor,
  (doc) => {
    if (!doc) {
      editor.clearDocument()
      return
    }
    editor.openDocument(doc)
  },
  { immediate: true },
)

watch(
  () => tree.selectedNodeId,
  async (nodeId) => {
    await extract.initialize()
    const node = tree.findNodeById(nodeId)
    if (node?.type === 'extract' && node.sourceExtractId) {
      extract.openExtract(node.sourceExtractId)
    } else if (store.currentContext === 'knowledge') {
      extract.setSelectedExtract(null)
    }
  },
  { immediate: true },
)

const title = computed(() => {
  switch (store.currentContext) {
    case 'inbox':
      return 'Inbox editor'
    case 'reading':
      return 'Editor'
    case 'knowledge':
      return 'Editor'
    case 'review':
      return 'Review'
    case 'search':
      return 'Result detail'
    default:
      return 'Detail'
  }
})
</script>

<template>
  <AppPane :title="title" class="read-pane" :scrollable="false" :padded="'none'">
    <template
      v-if="
        store.currentContext === 'reading' ||
        store.currentContext === 'inbox' ||
        store.currentContext === 'knowledge'
      "
    >
      <EditorShell />
    </template>

    <template v-else-if="store.currentContext === 'review'">
      <ReviewCard />
    </template>

    <template v-else-if="store.currentContext === 'search'">
      <div class="read-pane__block">
        <h2 class="read-pane__h">Search workspace</h2>
        <p class="read-pane__p">
          Use the tree pane as the global search hub. Selecting a result jumps directly to the
          target workflow (reading, knowledge, or review) so editing and review stay in their
          native modules.
        </p>
      </div>
    </template>
  </AppPane>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.read-pane {
  height: 100%;
}

.read-pane__block {
  padding: $space-md;
}

.read-pane__h {
  margin: 0 0 $space-md;
  font-size: $font-size-lg;
  font-weight: 600;
  line-height: $line-tight;
}

.read-pane__p {
  margin: 0 0 $space-md;
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: $line-normal;
}
</style>
