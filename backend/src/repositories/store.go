package repositories

import (
	"sync"

	"groundTurn/src/constants"
	"groundTurn/src/models"
)

// store 内存数据仓库：种子数据覆盖放行演示的三类场景
// （未完成任务、未关闭延误、预约时间冲突），换绑/放行在互斥锁内原子完成。
type store struct {
	mu          sync.Mutex
	turnarounds map[int]models.FlightTurnaround
	tasks       map[int]models.GroundTask
	resources   map[int]models.GroundResource
	bookings    map[int]models.ResourceBooking
	delays      map[int]models.DelayEvent
	nextBooking int
}

var db = seed()

func seed() *store {
	s := &store{
		turnarounds: map[int]models.FlightTurnaround{},
		tasks:       map[int]models.GroundTask{},
		resources:   map[int]models.GroundResource{},
		bookings:    map[int]models.ResourceBooking{},
		delays:      map[int]models.DelayEvent{},
	}
	for _, t := range []models.FlightTurnaround{
		{ID: 1, FlightNo: "CA1234", AircraftReg: "B-1234", StandNo: "101", ArrivalTime: "2026-09-25T08:00:00Z", DepartureTime: "2026-09-25T09:30:00Z", TurnaroundStatus: "IN_SERVICE", DelayReason: "前段天气流控"},
		{ID: 2, FlightNo: "MU5678", AircraftReg: "B-5678", StandNo: "102", ArrivalTime: "2026-09-25T08:30:00Z", DepartureTime: "2026-09-25T10:00:00Z", TurnaroundStatus: "IN_SERVICE", DelayReason: ""},
		{ID: 3, FlightNo: "CZ9012", AircraftReg: "B-9012", StandNo: "103", ArrivalTime: "2026-09-25T09:00:00Z", DepartureTime: "2026-09-25T10:30:00Z", TurnaroundStatus: "IN_SERVICE", DelayReason: ""},
	} {
		s.turnarounds[t.ID] = t
	}
	for _, t := range []models.GroundTask{
		{ID: 1, TurnaroundID: 1, TaskType: "CLEANING", TeamID: 11, PlannedStart: "2026-09-25T08:10:00Z", Deadline: "2026-09-25T08:50:00Z", ActualFinish: "", Status: constants.GroundTaskStatusInProgress, BlockerNote: "客舱遗留物清点中"},
		{ID: 2, TurnaroundID: 1, TaskType: "REFUEL", TeamID: 12, PlannedStart: "2026-09-25T08:20:00Z", Deadline: "2026-09-25T08:55:00Z", ActualFinish: "2026-09-25T08:50:00Z", Status: constants.GroundTaskStatusCompleted, BlockerNote: ""},
		{ID: 3, TurnaroundID: 1, TaskType: "BAGGAGE", TeamID: 13, PlannedStart: "2026-09-25T09:00:00Z", Deadline: "2026-09-25T09:40:00Z", ActualFinish: "2026-09-25T09:35:00Z", Status: constants.GroundTaskStatusCompleted, BlockerNote: ""},
		{ID: 4, TurnaroundID: 2, TaskType: "BAGGAGE", TeamID: 13, PlannedStart: "2026-09-25T08:40:00Z", Deadline: "2026-09-25T09:20:00Z", ActualFinish: "2026-09-25T09:15:00Z", Status: constants.GroundTaskStatusCompleted, BlockerNote: ""},
		{ID: 5, TurnaroundID: 2, TaskType: "CATERING", TeamID: 14, PlannedStart: "2026-09-25T08:45:00Z", Deadline: "2026-09-25T09:15:00Z", ActualFinish: "2026-09-25T09:10:00Z", Status: constants.GroundTaskStatusCompleted, BlockerNote: ""},
		{ID: 6, TurnaroundID: 3, TaskType: "CLEANING", TeamID: 11, PlannedStart: "2026-09-25T09:05:00Z", Deadline: "2026-09-25T09:40:00Z", ActualFinish: "2026-09-25T09:35:00Z", Status: constants.GroundTaskStatusCompleted, BlockerNote: ""},
		{ID: 7, TurnaroundID: 3, TaskType: "WATER_SERVICE", TeamID: 15, PlannedStart: "2026-09-25T09:10:00Z", Deadline: "2026-09-25T09:35:00Z", ActualFinish: "2026-09-25T09:30:00Z", Status: constants.GroundTaskStatusCompleted, BlockerNote: ""},
	} {
		s.tasks[t.ID] = t
	}
	for _, r := range []models.GroundResource{
		{ID: 1, ResourceCode: "CLN-01", ResourceType: "CLEANING", Location: "T1 西区", AvailabilityStatus: "AVAILABLE", MaintenanceDueAt: "2026-12-01T00:00:00Z", OwnerTeam: "保洁一组"},
		{ID: 2, ResourceCode: "CAT-01", ResourceType: "CATERING", Location: "T1 配餐区", AvailabilityStatus: "AVAILABLE", MaintenanceDueAt: "2026-11-15T00:00:00Z", OwnerTeam: "配餐组"},
		{ID: 3, ResourceCode: "BAG-01", ResourceType: "BAGGAGE", Location: "T1 行李区", AvailabilityStatus: "AVAILABLE", MaintenanceDueAt: "2026-10-20T00:00:00Z", OwnerTeam: "行李组"},
		{ID: 4, ResourceCode: "BAG-02", ResourceType: "BAGGAGE", Location: "T1 行李区", AvailabilityStatus: "AVAILABLE", MaintenanceDueAt: "2026-10-20T00:00:00Z", OwnerTeam: "行李组"},
		{ID: 5, ResourceCode: "REF-01", ResourceType: "REFUEL", Location: "机坪加油点", AvailabilityStatus: "AVAILABLE", MaintenanceDueAt: "2027-01-10T00:00:00Z", OwnerTeam: "油料组"},
		{ID: 6, ResourceCode: "WAT-01", ResourceType: "WATER_SERVICE", Location: "T1 东区", AvailabilityStatus: "AVAILABLE", MaintenanceDueAt: "2026-09-30T00:00:00Z", OwnerTeam: "勤务组"},
		{ID: 7, ResourceCode: "PUSH-01", ResourceType: "PUSHBACK", Location: "机坪牵引区", AvailabilityStatus: "AVAILABLE", MaintenanceDueAt: "2026-12-20T00:00:00Z", OwnerTeam: "牵引组"},
		{ID: 8, ResourceCode: "BAG-03", ResourceType: "BAGGAGE", Location: "T1 行李区", AvailabilityStatus: "MAINTENANCE", MaintenanceDueAt: "2026-09-26T08:00:00Z", OwnerTeam: "行李组"},
		{ID: 9, ResourceCode: "CAT-02", ResourceType: "CATERING", Location: "T2 配餐区", AvailabilityStatus: "OFFLINE", MaintenanceDueAt: "2026-10-01T00:00:00Z", OwnerTeam: "配餐组"},
	} {
		s.resources[r.ID] = r
	}
	for _, b := range []models.ResourceBooking{
		{ID: 1, ResourceID: 1, TurnaroundID: 1, TaskID: 1, StartTime: "2026-09-25T08:10:00Z", EndTime: "2026-09-25T08:50:00Z", BookingStatus: constants.BookingStatusConfirmed, ConflictReason: ""},
		{ID: 2, ResourceID: 5, TurnaroundID: 1, TaskID: 2, StartTime: "2026-09-25T08:20:00Z", EndTime: "2026-09-25T08:55:00Z", BookingStatus: constants.BookingStatusConfirmed, ConflictReason: ""},
		{ID: 3, ResourceID: 4, TurnaroundID: 1, TaskID: 3, StartTime: "2026-09-25T09:00:00Z", EndTime: "2026-09-25T09:40:00Z", BookingStatus: constants.BookingStatusConfirmed, ConflictReason: ""},
		{ID: 4, ResourceID: 4, TurnaroundID: 2, TaskID: 4, StartTime: "2026-09-25T08:40:00Z", EndTime: "2026-09-25T09:20:00Z", BookingStatus: constants.BookingStatusConfirmed, ConflictReason: ""},
		{ID: 5, ResourceID: 2, TurnaroundID: 2, TaskID: 5, StartTime: "2026-09-25T08:45:00Z", EndTime: "2026-09-25T09:15:00Z", BookingStatus: constants.BookingStatusConfirmed, ConflictReason: ""},
		{ID: 6, ResourceID: 1, TurnaroundID: 3, TaskID: 6, StartTime: "2026-09-25T09:05:00Z", EndTime: "2026-09-25T09:40:00Z", BookingStatus: constants.BookingStatusConfirmed, ConflictReason: ""},
		{ID: 7, ResourceID: 6, TurnaroundID: 3, TaskID: 7, StartTime: "2026-09-25T09:10:00Z", EndTime: "2026-09-25T09:35:00Z", BookingStatus: constants.BookingStatusConfirmed, ConflictReason: ""},
	} {
		s.bookings[b.ID] = b
	}
	s.nextBooking = 8
	for _, d := range []models.DelayEvent{
		{ID: 1, TurnaroundID: 1, DelayType: "WEATHER", Minutes: 25, RootCause: "前段航路雷雨流控", ResponsibilityTeam: "运控中心", ResolvedAt: ""},
		{ID: 2, TurnaroundID: 1, DelayType: "CLEANING", Minutes: 10, RootCause: "客舱遗留物等待认领", ResponsibilityTeam: "保洁一组", ResolvedAt: "2026-09-25T08:45:00Z"},
		{ID: 3, TurnaroundID: 2, DelayType: "BAGGAGE", Minutes: 10, RootCause: "行李分拣系统短暂故障", ResponsibilityTeam: "行李组", ResolvedAt: "2026-09-25T09:05:00Z"},
	} {
		s.delays[d.ID] = d
	}
	return s
}
