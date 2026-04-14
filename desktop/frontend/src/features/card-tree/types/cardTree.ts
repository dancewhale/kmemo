import type { CardDTO } from "../../../types/dto";

export interface CardTreeNodeVM {
  id: string;
  parentId: string | null;
  open: boolean;
  title: string;
  cardType: string;
  status: string;
  hasChildrenHint: boolean;
  depth: number;
  draggable: boolean;
  droppable: boolean;
  source: CardDTO;
  children: CardTreeNodeVM[];
}

export interface CardTreeDropPayload {
  cardId: string;
  sourceParentId: string | null;
  targetParentId: string | null;
  targetIndex: number;
  siblingIds: string[];
}
