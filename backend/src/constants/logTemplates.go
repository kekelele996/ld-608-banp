package constants

// LogTemplates is consumed by the audit-log middleware/service. Every write
// action records an entry; placeholders are filled with fmt.Sprintf.
var LogTemplates = map[string]string{
	"TURNAROUND_REGISTERED": "航班过站登记：航班 %s，机位 %s",
	"TURNAROUND_RELEASED":   "航班放行通过：航班 %s，状态 %s -> %s",
	"RELEASE_REJECTED":      "航班放行驳回：航班 %s，阻塞项 %d（任务 %d / 延误 %d / 冲突 %d）",

	"TASK_DISPATCHED": "任务派发：#%d %s -> 班组 %d",
	"TASK_COMPLETED":  "任务签收完成：#%d %s",
	"TASK_BLOCKED":    "任务阻塞：#%d，原因 %s",

	"RESOURCE_BOOKED":   "资源预约：预约 #%d，资源 %s，任务 #%d",
	"RESOURCE_RELEASED": "资源释放：预约 #%d，资源 %s",
	"RESOURCE_REBOUND":  "资源换绑：预约 #%d -> 新预约 #%d，资源 %s -> %s，任务 #%d",
	"CONFLICT_DETECTED": "预约冲突：#%d 与 #%d 在资源 %s 重叠",

	"DELAY_REPORTED": "延误登记：航班 %s，#%d %s，%d 分钟",
	"DELAY_RESOLVED": "延误关闭：#%d，航班 %s",
}
