import { callApp } from "./appBridge";
import type { KnowledgeTreeNode } from "../types/dto";

export function getKnowledgeTree(rootId?: string | null) {
  return callApp<KnowledgeTreeNode[]>("GetKnowledgeTree", rootId ?? null);
}
