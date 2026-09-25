import { GroundTaskTypeText } from "../constants/GroundTaskType";
import { GroundTaskStatusText } from "../constants/GroundTaskStatus";
import { TurnaroundStatusText } from "../constants/TurnaroundStatus";
import { ResourceStatusText } from "../constants/ResourceStatus";
import { BookingStatusText } from "../constants/BookingStatus";

export const formatDate = (value: string) => new Date(value).toLocaleString("zh-CN");
export const formatStatus = (value: string) => value.replace(/_/g, " ");
export const formatNumber = (value: number) => new Intl.NumberFormat("zh-CN").format(value);
export const formatRisk = (value: string) => ({ LOW: "低", MEDIUM: "中", HIGH: "高", CRITICAL: "严重", EXTREME: "极高" }[value] ?? value);

// HH:mm time window used by booking rows and the conflict panel.
export const formatTime = (value?: string | null) =>
  value ? new Date(value).toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit", hour12: false }) : "—";

export const formatDateTime = (value?: string | null) =>
  value
    ? new Date(value).toLocaleString("zh-CN", { month: "2-digit", day: "2-digit", hour: "2-digit", minute: "2-digit", hour12: false })
    : "—";

export const formatTaskType = (value: string) =>
  (GroundTaskTypeText as Record<string, string>)[value] ?? value;

export const formatTaskStatus = (value: string) =>
  (GroundTaskStatusText as Record<string, string>)[value] ?? value;

export const formatTurnaroundStatus = (value: string) =>
  (TurnaroundStatusText as Record<string, string>)[value] ?? value;

export const formatResourceStatus = (value: string) =>
  (ResourceStatusText as Record<string, string>)[value] ?? value;

export const formatBookingStatus = (value: string) =>
  (BookingStatusText as Record<string, string>)[value] ?? value;
