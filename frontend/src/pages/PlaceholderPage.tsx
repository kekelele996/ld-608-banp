import { StatusBadge } from "../components/common/StatusBadge";

// PlaceholderPage keeps the non-feature pages reachable in the navigation
// while the release coordination flow owns the implemented workbench.
export function PlaceholderPage({ name }: { name: string }) {
  return (
    <main className="page">
      <section className="page-head">
        <div>
          <p className="eyebrow">ground-turn</p>
          <h1>{name}</h1>
        </div>
        <StatusBadge value="ON_STAND" />
      </section>
      <section className="workbench">
        <div className="panel wide">
          <h2>{name}模块</h2>
          <p className="dim">
            本页为台账浏览模块。本次上线的协同流程位于「过站放行」：未完成任务、未关闭延误与预约冲突汇总、放行核对与资源换绑。
          </p>
        </div>
      </section>
    </main>
  );
}
