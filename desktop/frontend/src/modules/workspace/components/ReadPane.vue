<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useWorkspaceStore } from '../stores/workspace.store'
import AppPane from '@/shared/components/AppPane.vue'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useEditorStore } from '@/modules/editor/stores/editor.store'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import EditorShell from '@/modules/editor/components/EditorShell.vue'
import ExtractDetailPanel from '@/modules/extract/components/ExtractDetailPanel.vue'
import CardDetailPanel from '@/modules/card/components/CardDetailPanel.vue'
import { useCardStore } from '@/modules/card/stores/card.store'
import ReviewCard from '@/modules/review/components/ReviewCard.vue'

const store = useWorkspaceStore()
const tree = useTreeStore()
const reader = useReaderStore()
const editor = useEditorStore()
const extract = useExtractStore()
const card = useCardStore()

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
const knChildCount = computed(() => (store.currentContext === 'knowledge' ? tree.selectedNodeChildCount : 0))
const selectedExtract = computed(() => {
  const extractId = kn.value?.sourceExtractId
  if (!extractId) {
    return null
  }
  return extract.items.find((item) => item.id === extractId) ?? null
})
const showExtractDetail = computed(() => {
  return store.currentContext === 'knowledge' && tree.isSelectedExtractNode
})

const showCardDetail = computed(() => {
  return store.currentContext === 'knowledge' && kn.value?.type === 'card'
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
  activeEditorArticle,
  (doc) => {
    if (!doc) {
      editor.clearDocument()
      return
    }
    editor.openDocument({
      id: doc.id,
      title: doc.title,
      content: doc.content,
      contentType: 'article',
      updatedAt: doc.updatedAt,
      sourceUrl: doc.sourceUrl,
      tags: doc.tags,
    })
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
      return 'Knowledge detail'
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
    <template v-if="store.currentContext === 'reading' || store.currentContext === 'inbox'">
      <EditorShell />
    </template>

    <template v-else-if="store.currentContext === 'knowledge'">
      <AppEmpty v-if="!kn" message="Select a node in the tree" />
      <ExtractDetailPanel v-else-if="showExtractDetail" />
      <CardDetailPanel v-else-if="showCardDetail" />
      <div v-else class="read-pane__block">
        <h2 class="read-pane__h">{{ kn.title }}</h2>
        <template v-if="kn.description && (kn.type === 'topic' || kn.type === 'card')">
          <p class="read-pane__label">Note</p>
          <p class="read-pane__p">{{ kn.description }}</p>
        </template>
        <p v-if="selectedExtract" class="read-pane__label">Extract quote</p>
        <p v-if="selectedExtract" class="read-pane__p read-pane__p--quote">“{{ selectedExtract.quote }}”</p>
        <p v-if="selectedExtract?.note" class="read-pane__label">Note</p>
        <p v-if="selectedExtract?.note" class="read-pane__p">{{ selectedExtract.note }}</p>
        <dl class="read-pane__dl">
          <div><dt>Id</dt><dd>{{ kn.id }}</dd></div>
          <div><dt>Type</dt><dd>{{ kn.type }}</dd></div>
          <div><dt>Parent</dt><dd>{{ kn.parentId ?? '—' }}</dd></div>
          <div><dt>Child count</dt><dd>{{ knChildCount }}</dd></div>
          <div v-if="selectedExtract"><dt>Source article</dt><dd>{{ selectedExtract.sourceArticleTitle }}</dd></div>
          <div v-if="kn.createdAt"><dt>Created</dt><dd>{{ kn.createdAt }}</dd></div>
          <div v-if="kn.updatedAt"><dt>Updated</dt><dd>{{ kn.updatedAt }}</dd></div>
        </dl>
      </div>
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

.read-pane__label {
  margin: 0 0 $space-xs;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.read-pane__p {
  margin: 0 0 $space-md;
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: $line-normal;
}

.read-pane__p--quote {
  color: $color-text;
  font-style: italic;
}

.read-pane__dl {
  margin: 0;
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.read-pane__dl > div {
  display: grid;
  grid-template-columns: 72px 1fr;
  gap: $space-sm;
  padding: $space-xs 0;
  border-top: 1px solid $color-border-subtle;
}

.read-pane__dl dt {
  margin: 0;
  color: $color-text-secondary;
}

.read-pane__dl dd {
  margin: 0;
  color: $color-text;
  font-family: $font-mono;
  word-break: break-all;
}
</style>
