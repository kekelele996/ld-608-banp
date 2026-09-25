package constants

// Error message templates. Services fill the placeholders via fmt.Sprintf.
// They intentionally live separately from error codes so wording changes do
// not touch branching logic.
const (
	AuthRequiredMessage = "missing token"

	TurnaroundNotFoundMessage      = "航班 %v 不存在"
	TurnaroundNotReleasableMessage = "航班 %s 当前状态 %s，不允许重复放行"
	TaskUnfinishedMessage          = "任务 #%d（%s）尚未完成，状态为 %s"
	TaskNotFoundMessage            = "任务 #%d 不存在"
	DelayNotClosedMessage          = "航班 %s 存在未关闭延误 #%d（%s，%d 分钟）"
	BookingConflictMessage         = "预约 #%d 与预约 #%d 在资源 %s 上时间重叠"
	BookingNotFoundMessage         = "预约 #%d 不存在"
	BookingNotRebindableMessage    = "预约 #%d 状态为 %s，仅未释放预约允许换绑"
	ResourceUnavailableMessage     = "资源 %s 当前状态为 %s，不可换用"
	ResourceNotFoundMessage        = "资源 #%d 不存在"
	ResourceTypeMismatchMessage    = "资源 %s 类型为 %s，与任务 %s 不匹配"
	ResourceBusyMessage            = "资源 %s 在 %s ~ %s 已被预约 #%d 占用"

	ReleaseRejectedMessage = "放行核对未通过：存在 %d 项阻塞条件"
)
