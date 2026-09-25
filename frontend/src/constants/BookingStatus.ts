export const BookingStatus = ["CONFIRMED", "CONFLICT", "RELEASED"] as const;
export type BookingStatus = (typeof BookingStatus)[number];
export const BookingStatusText: Record<BookingStatus, string> = {
  CONFIRMED: "已确认",
  CONFLICT: "冲突",
  RELEASED: "已释放"
};
