import { formatBookingStatus, formatResourceStatus, formatTaskStatus, formatTurnaroundStatus } from "../../utils/formatters";

// StatusBadge renders any domain status with its shared Chinese label and a
// CSS class derived from the raw enum value. The label lookup is centralized
// in utils/formatters so adding an enum value touches one map.
export function StatusBadge({ value }: { value: string }) {
  const label =
    formatTurnaroundStatus(value) !== value ? formatTurnaroundStatus(value)
    : formatTaskStatus(value) !== value ? formatTaskStatus(value)
    : formatResourceStatus(value) !== value ? formatResourceStatus(value)
    : formatBookingStatus(value) !== value ? formatBookingStatus(value)
    : value.replace(/_/g, " ");
  return <span className={"badge " + String(value).toLowerCase().replace(/_/g, "-")}>{label}</span>;
}
