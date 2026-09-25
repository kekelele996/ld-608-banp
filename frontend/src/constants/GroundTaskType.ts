export const GroundTaskType = ["CLEANING", "CATERING", "BAGGAGE", "REFUEL", "WATER_SERVICE", "PUSHBACK"] as const;
export type GroundTaskType = (typeof GroundTaskType)[number];
export const GroundTaskTypeText: Record<GroundTaskType, string> = {
  CLEANING: "客舱清洁",
  CATERING: "航食配餐",
  BAGGAGE: "行李装卸",
  REFUEL: "航油加注",
  WATER_SERVICE: "清水保障",
  PUSHBACK: "牵引车推出"
};
