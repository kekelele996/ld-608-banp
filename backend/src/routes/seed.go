package routes

import (
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"

	"gorm.io/gorm"
)

// seedIfEmpty inserts a deterministic local dataset used for offline review.
// The data intentionally contains all three blocker kinds on flight CA1852 so
// the release coordination flow can be exercised end to end:
//   - an unfinished BLOCKED refuel task
//   - an open (unresolved) delay
//   - a catering booking overlapping another turnaround on the same truck
//
// Flight CA2477 is fully ready and should pass release immediately.
func seedIfEmpty(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.FlightTurnaround{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	base := time.Date(2026, 9, 25, 0, 0, 0, 0, time.Local)
	at := func(h, m int) time.Time { return base.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute) }

	flights := []models.FlightTurnaround{
		{FlightNo: "CA1852", AircraftReg: "B-2031", StandNo: "212", ArrivalTime: at(9, 0), DepartureTime: at(10, 5), TurnaroundStatus: constants.TurnaroundStatusInService, DelayReason: "等待加油车"},
		{FlightNo: "CA2477", AircraftReg: "B-6712", StandNo: "108", ArrivalTime: at(9, 30), DepartureTime: at(10, 40), TurnaroundStatus: constants.TurnaroundStatusInService, DelayReason: ""},
		{FlightNo: "MU5180", AircraftReg: "B-8866", StandNo: "305", ArrivalTime: at(10, 10), DepartureTime: at(11, 20), TurnaroundStatus: constants.TurnaroundStatusOnStand, DelayReason: ""},
	}
	if err := db.Create(&flights).Error; err != nil {
		return err
	}
	ca, ca2, mu := flights[0].ID, flights[1].ID, flights[2].ID

	resources := []models.GroundResource{
		{ResourceCode: "CAT-01", ResourceType: constants.GroundTaskTypeCatering, Location: "北区餐食库", AvailabilityStatus: constants.ResourceStatusBooked, OwnerTeam: "配餐一组"},
		{ResourceCode: "CAT-02", ResourceType: constants.GroundTaskTypeCatering, Location: "南区餐食库", AvailabilityStatus: constants.ResourceStatusAvailable, OwnerTeam: "配餐二组"},
		{ResourceCode: "FUEL-01", ResourceType: constants.GroundTaskTypeRefuel, Location: "油库A", AvailabilityStatus: constants.ResourceStatusBooked, OwnerTeam: "加油一组"},
		{ResourceCode: "FUEL-02", ResourceType: constants.GroundTaskTypeRefuel, Location: "油库B", AvailabilityStatus: constants.ResourceStatusAvailable, OwnerTeam: "加油二组"},
		{ResourceCode: "BAG-01", ResourceType: constants.GroundTaskTypeBaggage, Location: "行李转盘3", AvailabilityStatus: constants.ResourceStatusBooked, OwnerTeam: "装卸一组"},
		{ResourceCode: "CLN-01", ResourceType: constants.GroundTaskTypeCleaning, Location: "保洁站", AvailabilityStatus: constants.ResourceStatusBooked, OwnerTeam: "保洁一组"},
		{ResourceCode: "WTR-01", ResourceType: constants.GroundTaskTypeWaterService, Location: "净水站", AvailabilityStatus: constants.ResourceStatusMaintenance, OwnerTeam: "水务一组"},
		{ResourceCode: "PUSH-01", ResourceType: constants.GroundTaskTypePushback, Location: "牵引车坪", AvailabilityStatus: constants.ResourceStatusBooked, OwnerTeam: "牵引车组"},
	}
	if err := db.Create(&resources).Error; err != nil {
		return err
	}
	cat01, cat02, fuel01, bag01, cln01, push01 := resources[0].ID, resources[1].ID, resources[2].ID, resources[4].ID, resources[5].ID, resources[7].ID

	tasks := []models.GroundTask{
		// CA1852: one blocked refuel task -> unfinished blocker
		{TurnaroundID: ca, TaskType: constants.GroundTaskTypeCatering, TeamID: 1, PlannedStart: at(9, 10), Deadline: at(9, 40), Status: constants.GroundTaskStatusCompleted, ActualFinish: ptrTime(at(9, 35))},
		{TurnaroundID: ca, TaskType: constants.GroundTaskTypeRefuel, TeamID: 3, PlannedStart: at(9, 20), Deadline: at(9, 55), Status: constants.GroundTaskStatusBlocked, BlockerNote: "加油车被 CA2477 占用"},
		{TurnaroundID: ca, TaskType: constants.GroundTaskTypeCleaning, TeamID: 2, PlannedStart: at(9, 15), Deadline: at(9, 45), Status: constants.GroundTaskStatusCompleted, ActualFinish: ptrTime(at(9, 42))},
		{TurnaroundID: ca, TaskType: constants.GroundTaskTypePushback, TeamID: 5, PlannedStart: at(9, 55), Deadline: at(10, 0), Status: constants.GroundTaskStatusPending},

		// CA2477: all tasks complete -> ready to release
		{TurnaroundID: ca2, TaskType: constants.GroundTaskTypeCatering, TeamID: 1, PlannedStart: at(10, 5), Deadline: at(10, 35), Status: constants.GroundTaskStatusCompleted, ActualFinish: ptrTime(at(10, 33))},
		{TurnaroundID: ca2, TaskType: constants.GroundTaskTypeRefuel, TeamID: 3, PlannedStart: at(9, 45), Deadline: at(10, 15), Status: constants.GroundTaskStatusCompleted, ActualFinish: ptrTime(at(10, 10))},
		{TurnaroundID: ca2, TaskType: constants.GroundTaskTypeBaggage, TeamID: 4, PlannedStart: at(9, 35), Deadline: at(10, 20), Status: constants.GroundTaskStatusCompleted, ActualFinish: ptrTime(at(10, 18))},

		// MU5180: pending tasks, open delay, booking conflict
		{TurnaroundID: mu, TaskType: constants.GroundTaskTypeBaggage, TeamID: 4, PlannedStart: at(10, 20), Deadline: at(10, 50), Status: constants.GroundTaskStatusInProgress},
		{TurnaroundID: mu, TaskType: constants.GroundTaskTypeCatering, TeamID: 1, PlannedStart: at(9, 25), Deadline: at(9, 55), Status: constants.GroundTaskStatusPending},
	}
	if err := db.Create(&tasks).Error; err != nil {
		return err
	}
	caCateringTask, caFuelTask := tasks[0].ID, tasks[1].ID
	ca2CateringTask := tasks[4].ID
	muCateringTask := tasks[8].ID

	bookings := []models.ResourceBooking{
		// CA1852 catering uses CAT-01 09:10-09:40
		{ResourceID: cat01, TurnaroundID: ca, TaskID: caCateringTask, StartTime: at(9, 10), EndTime: at(9, 40), BookingStatus: constants.BookingStatusConfirmed},
		// CA1852 refuel booking on FUEL-01 09:20-09:55 (resource busy looking)
		{ResourceID: fuel01, TurnaroundID: ca, TaskID: caFuelTask, StartTime: at(9, 20), EndTime: at(9, 55), BookingStatus: constants.BookingStatusConfirmed},

		// MU5180 catering on CAT-01 09:25-09:55 -> overlaps CA1852's CAT-01
		// booking (conflict pair); the matching task is pending so dispatchers
		// can either swap it to CAT-02 or finish the other flight's service.
		{ResourceID: cat01, TurnaroundID: mu, TaskID: muCateringTask, StartTime: at(9, 25), EndTime: at(9, 55), BookingStatus: constants.BookingStatusConflict, ConflictReason: "与 CA1852 配餐预约重叠"},

		// CA2477 completed catering uses CAT-02 (a different, free truck;
		// window 10:10-10:35 leaves CAT-02 free for the 09:25-09:55 swap)
		{ResourceID: cat02, TurnaroundID: ca2, TaskID: ca2CateringTask, StartTime: at(10, 10), EndTime: at(10, 35), BookingStatus: constants.BookingStatusConfirmed},
		{ResourceID: bag01, TurnaroundID: ca2, TaskID: tasks[6].ID, StartTime: at(9, 35), EndTime: at(10, 20), BookingStatus: constants.BookingStatusConfirmed},
		{ResourceID: cln01, TurnaroundID: ca, TaskID: tasks[2].ID, StartTime: at(9, 15), EndTime: at(9, 45), BookingStatus: constants.BookingStatusConfirmed},
		{ResourceID: push01, TurnaroundID: ca, TaskID: tasks[3].ID, StartTime: at(9, 55), EndTime: at(10, 0), BookingStatus: constants.BookingStatusConfirmed},
	}
	if err := db.Create(&bookings).Error; err != nil {
		return err
	}

	delays := []models.DelayEvent{
		// Open delay on CA1852 -> release blocker
		{TurnaroundID: ca, DelayType: constants.GroundTaskTypeRefuel, Minutes: 25, RootCause: "加油车调配冲突", ResponsibilityTeam: "加油一组"},
		// Closed delay on CA2477 -> not a blocker
		{TurnaroundID: ca2, DelayType: constants.GroundTaskTypeBaggage, Minutes: 10, RootCause: "行李分拣延迟", ResponsibilityTeam: "装卸一组", ResolvedAt: ptrTime(at(10, 25))},
	}
	return db.Create(&delays).Error
}

func ptrTime(t time.Time) *time.Time { return &t }
