export const GroundTaskStatus = ["PENDING", "IN_PROGRESS", "BLOCKED", "COMPLETED"] as const;
export type GroundTaskStatus = (typeof GroundTaskStatus)[number];
export const GroundTaskStatusText: Record<GroundTaskStatus, string> = {
  PENDING: "待执行",
  IN_PROGRESS: "执行中",
  BLOCKED: "阻塞",
  COMPLETED: "已完成"
};
