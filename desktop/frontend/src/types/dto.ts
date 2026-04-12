export interface KnowledgeTreeNode {
  id: string;
  name: string;
  description: string;
  parentId?: string | null;
  cardCount: number;
  dueCount: number;
  children: KnowledgeTreeNode[];
}

export interface TagDTO {
  id: string;
  name: string;
  slug: string;
  color: string;
  icon: string;
  description: string;
  cardCount: number;
  createdAt: string;
}

export interface SRSDTO {
  cardId: string;
  fsrsState: string;
  dueAt?: string | null;
  lastReviewAt?: string | null;
  stability?: number | null;
  difficulty?: number | null;
  reps: number;
  lapses: number;
}

export interface CardDTO {
  id: string;
  knowledgeId: string;
  knowledgeName: string;
  sourceDocumentId?: string | null;
  parentId?: string | null;
  title: string;
  cardType: string;
  htmlPath: string;
  status: string;
  tags?: TagDTO[];
  srs?: SRSDTO | null;
  htmlContent?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CardSummaryDTO {
  id: string;
  parentId?: string | null;
  title: string;
  cardType: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface CardDetailDTO extends CardDTO {
  parent?: CardSummaryDTO | null;
  children: CardSummaryDTO[];
}

export interface CardFilters {
  knowledgeId?: string | null;
  cardType?: string;
  status?: string;
  tagIds?: string[];
  keyword?: string;
  parentId?: string | null;
  isRoot?: boolean;
  orderBy?: string;
  orderDesc?: boolean;
  limit?: number;
  offset?: number;
}

export interface CreateCardRequest {
  knowledgeId: string;
  sourceDocumentId?: string | null;
  parentId?: string | null;
  title: string;
  cardType: string;
  htmlContent: string;
  tagIds: string[];
}
