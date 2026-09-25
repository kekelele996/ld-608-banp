export const ReleaseConclusion = ["RELEASABLE","BLOCKED"] as const;
export type ReleaseConclusion = (typeof ReleaseConclusion)[number];
export const ReleaseConclusionText: Record<ReleaseConclusion, string> = {
  RELEASABLE: "可放行",
  BLOCKED: "暂不可放行"
};
export const ReleaseBlockerKindText: Record<string, string> = {
  TASK: "未完成任务",
  DELAY: "未关闭延误",
  BOOKING: "预约冲突"
};
