import { StatusBadge } from "./StatusBadge";

export function ConflictBadge({ title = "预约冲突", value = "CONFLICT" }: { title?: string; value?: string }) {
  return (
    <span className="conflict-badge" title={title}>
      <StatusBadge value={value} />
    </span>
  );
}
