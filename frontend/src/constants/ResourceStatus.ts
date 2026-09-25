export const ResourceStatus = ["AVAILABLE", "BOOKED", "MAINTENANCE", "OFFLINE"] as const;
export type ResourceStatus = (typeof ResourceStatus)[number];
export const ResourceStatusText: Record<ResourceStatus, string> = {
  AVAILABLE: "可用",
  BOOKED: "已占用",
  MAINTENANCE: "维护中",
  OFFLINE: "停用"
};
