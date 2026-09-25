import { request } from "./client";
import type { GroundResource } from "../types/GroundResource";

export function listGroundResources(): Promise<GroundResource[]> {
  return request<GroundResource[]>("/api/ground-resources");
}
