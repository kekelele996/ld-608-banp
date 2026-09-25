export const BookingStatus = ["PENDING","CONFIRMED","RELEASED","CANCELLED"] as const;
export type BookingStatus = (typeof BookingStatus)[number];
export const BookingStatusText: Record<BookingStatus, string> = {
  PENDING: "待确认",
  CONFIRMED: "已确认",
  RELEASED: "已释放",
  CANCELLED: "已取消"
};
