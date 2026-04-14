import { callApp } from "./appBridge";
import type { CardDTO, CardDetailDTO, CardFilters, CreateCardRequest, ListCardsResult } from "../types/dto";

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
