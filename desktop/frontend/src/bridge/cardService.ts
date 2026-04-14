import { callApp } from "./appBridge";
import type {
  CardDTO,
  CardDetailDTO,
  CardFilters,
  CreateCardRequest,
  ListCardsResult,
  MoveCardRequest,
  ReorderCardChildrenRequest,
  UpdateCardRequest,
} from "../types/dto";

export async function listCards(filters: CardFilters) {
  return callApp<ListCardsResult>("ListCards", filters);
}

export function getCard(id: string) {
  return callApp<CardDTO>("GetCard", id);
}

export function getCardDetail(id: string) {
  return callApp<CardDetailDTO>("GetCardDetail", id);
}

export function getCardChildren(parentId: string) {
  return callApp<CardDTO[]>("GetCardChildren", parentId);
}

export function createCard(payload: CreateCardRequest) {
  return callApp<string>("CreateCard", payload);
}

export function updateCard(id: string, payload: UpdateCardRequest) {
  return callApp<void>("UpdateCard", id, payload);
}

export function deleteCard(id: string) {
  return callApp<void>("DeleteCard", id);
}

export function moveCard(payload: MoveCardRequest) {
  return callApp<void>("MoveCard", payload);
}

export function reorderCardChildren(payload: ReorderCardChildrenRequest) {
  return callApp<void>("ReorderCardChildren", payload);
}
