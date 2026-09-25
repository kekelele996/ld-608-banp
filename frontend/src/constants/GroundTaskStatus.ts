export const GroundTaskStatus = ["PENDING","IN_PROGRESS","BLOCKED","COMPLETED"] as const;
export type GroundTaskStatus = (typeof GroundTaskStatus)[number];
export const GroundTaskStatusText: Record<GroundTaskStatus, string> = {
  PENDING: "待开始",
  IN_PROGRESS: "进行中",
  BLOCKED: "已阻塞",
  COMPLETED: "已完成"
};
