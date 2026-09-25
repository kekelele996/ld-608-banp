export const TurnaroundStatus = ["ARRIVING", "ON_STAND", "IN_SERVICE", "READY", "DEPARTED", "DELAYED"] as const;
export type TurnaroundStatus = (typeof TurnaroundStatus)[number];
export const TurnaroundStatusText: Record<TurnaroundStatus, string> = {
  ARRIVING: "即将到港",
  ON_STAND: "靠桥在站",
  IN_SERVICE: "保障中",
  READY: "已放行",
  DEPARTED: "已离港",
  DELAYED: "延误"
};
