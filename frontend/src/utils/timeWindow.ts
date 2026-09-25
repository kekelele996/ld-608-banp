// overlaps 判断两个时间窗口 [start, end) 是否重叠，与后端 utils/timeWindow.go 同规则。
export function overlaps(aStart: string, aEnd: string, bStart: string, bEnd: string): boolean {
  const as = Date.parse(aStart), ae = Date.parse(aEnd), bs = Date.parse(bStart), be = Date.parse(bEnd);
  if ([as, ae, bs, be].some((t) => Number.isNaN(t))) return false;
  return as < be && bs < ae;
}
