import { StatusBadge } from "./StatusBadge";

export function DelayTag({ title = "延误", value = "DELAYED" }: { title?: string; value?: string }) {
  return (
    <span className="delay-tag" title={title}>
      <StatusBadge value={value} />
    </span>
  );
}
