import { createCard, deleteCard, getCardDetail, moveCard, reorderCardChildren, updateCard } from "../../../bridge/cardService";
import type { CardDTO } from "../../../types/dto";
import type { CardTreeDropPayload } from "../types/cardTree";

interface UseCardTreeCommandsOptions {
  getKnowledgeId: () => string | null;
  onReload: () => Promise<void>;
  onSelect: (id: string) => Promise<void>;
}

function toErrorMessage(err: unknown, fallback: string): string {
  return err instanceof Error ? err.message : fallback;
}

export function useCardTreeCommands(options: UseCardTreeCommandsOptions) {
  async function onNodeCreateChild(parent: CardDTO): Promise<void> {
    const knowledgeId = options.getKnowledgeId();
    if (!knowledgeId) return;
    const title = window.prompt("子卡片标题", `New child of ${parent.title}`);
    if (!title || !title.trim()) return;
    try {
      const id = await createCard({
        knowledgeId,
        parentId: parent.id,
        title: title.trim(),
        cardType: parent.cardType || "article",
        htmlContent: "",
        tagIds: [],
      });
      await options.onReload();
      await options.onSelect(id);
    } catch (err) {
      window.alert(toErrorMessage(err, "创建子卡片失败"));
    }
  }

  async function onNodeRename(node: CardDTO): Promise<void> {
    const nextTitle = window.prompt("重命名卡片", node.title);
    if (!nextTitle || !nextTitle.trim() || nextTitle.trim() === node.title) return;
    try {
      const detail = await getCardDetail(node.id);
      await updateCard(node.id, {
        title: nextTitle.trim(),
        htmlContent: detail.htmlContent ?? "",
        status: detail.status ?? node.status ?? "draft",
      });
      await options.onReload();
      await options.onSelect(node.id);
    } catch (err) {
      window.alert(toErrorMessage(err, "重命名失败"));
    }
  }

  async function onNodeDelete(node: CardDTO): Promise<void> {
    const ok = window.confirm(`删除卡片「${node.title}」？该操作不可撤销。`);
    if (!ok) return;
    try {
      await deleteCard(node.id);
      await options.onReload();
    } catch (err) {
      window.alert(toErrorMessage(err, "删除失败"));
    }
  }

  async function onNodeDrop(payload: CardTreeDropPayload): Promise<void> {
    const knowledgeId = options.getKnowledgeId();
    if (!knowledgeId) return;
    try {
      await moveCard({
        cardId: payload.cardId,
        targetParentId: payload.targetParentId,
        targetIndex: payload.targetIndex,
      });
      await reorderCardChildren({
        knowledgeId,
        parentId: payload.targetParentId,
        orderedChildIds: payload.siblingIds,
      });
      await options.onReload();
    } catch (err) {
      throw new Error(toErrorMessage(err, "拖拽移动失败"));
    }
  }

  return {
    onNodeCreateChild,
    onNodeRename,
    onNodeDelete,
    onNodeDrop,
  };
}
