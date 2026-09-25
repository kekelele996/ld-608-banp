export interface GroundTask {
  id: number;
  turnaround_id: number;
  task_type: string;
  team_id: number;
  planned_start: string;
  deadline: string;
  actual_finish: string | null;
  status: string;
  blocker_note: string;
}
