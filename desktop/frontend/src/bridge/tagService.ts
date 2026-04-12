import { callApp } from "./appBridge";
import type { TagDTO } from "../types/dto";

export function listTags() {
  return callApp<TagDTO[]>("ListTags");
}
