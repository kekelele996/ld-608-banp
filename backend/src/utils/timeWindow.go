package utils

import "time"

// Overlaps 判断两个时间窗口 [start, end) 是否重叠，无法解析的时间视为不重叠。
func Overlaps(aStart, aEnd, bStart, bEnd string) bool {
	as, err1 := time.Parse(time.RFC3339, aStart)
	ae, err2 := time.Parse(time.RFC3339, aEnd)
	bs, err3 := time.Parse(time.RFC3339, bStart)
	be, err4 := time.Parse(time.RFC3339, bEnd)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return false
	}
	return as.Before(be) && bs.Before(ae)
}
