import { request } from "./client";
import type { GroundTask } from "../types/GroundTask";

export function listGroundTasks(turnaroundId?: number): Promise<GroundTask[]> {
  const qs = turnaroundId ? `?turnaround_id=${turnaroundId}` : "";
  return request<GroundTask[]>(`/api/ground-tasks${qs}`);
}

export function updateTaskStatus(id: number, status: string, blockerNote = ""): Promise<GroundTask> {
  return request<GroundTask>(`/api/ground-tasks/${id}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status, blocker_note: blockerNote, actor: "dispatcher" })
  });
}
