//go:build sqlite

package services_test

import (
	"testing"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/services"
	"groundTurn/src/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.FlightTurnaround{}, &models.GroundTask{}, &models.GroundResource{},
		&models.ResourceBooking{}, &models.DelayEvent{}, &models.AuditLog{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM resource_booking")
		db.Exec("DELETE FROM ground_task")
		db.Exec("DELETE FROM delay_event")
		db.Exec("DELETE FROM ground_resource")
		db.Exec("DELETE FROM flight_turnaround")
	})
	return db
}

// TestReleaseRejectedLeavesEverythingUntouched is the core guarantee: when
// any of the three conditions fails, flight status and original bookings are
// byte-for-byte unchanged after the rejected submission.
func TestReleaseRejectedLeavesEverythingUntouched(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.New(db)
	conflict := services.NewConflictService(repo)
	release := services.NewReleaseService(repo, conflict)

	start := time.Date(2026, 9, 25, 9, 0, 0, 0, time.UTC)
	flight := models.FlightTurnaround{FlightNo: "X1", TurnaroundStatus: constants.TurnaroundStatusInService, ArrivalTime: start, DepartureTime: start.Add(time.Hour)}
	if err := db.Create(&flight).Error; err != nil {
		t.Fatal(err)
	}
	task := models.GroundTask{TurnaroundID: flight.ID, TaskType: constants.GroundTaskTypeCatering, Status: constants.GroundTaskStatusPending, PlannedStart: start, Deadline: start.Add(time.Hour)}
	db.Create(&task)
	db.Create(&models.DelayEvent{TurnaroundID: flight.ID, DelayType: "WEATHER", Minutes: 15})

	res := models.GroundResource{ResourceCode: "R1", ResourceType: constants.GroundTaskTypeCatering, AvailabilityStatus: constants.ResourceStatusAvailable}
	db.Create(&res)
	booking := models.ResourceBooking{ResourceID: res.ID, TurnaroundID: flight.ID, TaskID: task.ID, StartTime: start, EndTime: start.Add(time.Hour), BookingStatus: constants.BookingStatusConfirmed}
	db.Create(&booking)

	statusBefore := flight.TurnaroundStatus
	bookingBefore := booking.BookingStatus

	if _, err := release.SubmitRelease(flight.ID, "tester"); err == nil {
		t.Fatal("expected release to be rejected, got success")
	} else {
		bizErr := err.(*services.BusinessError)
		if bizErr.Code != constants.ReleaseRejected {
			t.Fatalf("code = %s, want RELEASE_REJECTED", bizErr.Code)
		}
		if len(bizErr.Violations) != 2 { // 1 unfinished task + 1 open delay
			t.Fatalf("violations = %d, want 2", len(bizErr.Violations))
		}
	}

	var freshFlight models.FlightTurnaround
	db.First(&freshFlight, flight.ID)
	if freshFlight.TurnaroundStatus != statusBefore {
		t.Fatalf("flight status changed on rejection: %s -> %s", statusBefore, freshFlight.TurnaroundStatus)
	}
	var freshBooking models.ResourceBooking
	db.First(&freshBooking, booking.ID)
	if freshBooking.BookingStatus != bookingBefore {
		t.Fatalf("booking status changed on rejection: %s -> %s", bookingBefore, freshBooking.BookingStatus)
	}
}

// TestRebookLinksOldAndNew verifies the before/after history chain.
func TestRebookLinksOldAndNew(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.New(db)
	conflict := services.NewConflictService(repo)
	rebook := services.NewRebookService(repo, conflict)

	start := time.Date(2026, 9, 25, 9, 0, 0, 0, time.UTC)
	flight := models.FlightTurnaround{FlightNo: "X2", TurnaroundStatus: constants.TurnaroundStatusInService, ArrivalTime: start, DepartureTime: start.Add(2 * time.Hour)}
	db.Create(&flight)
	task := models.GroundTask{TurnaroundID: flight.ID, TaskType: constants.GroundTaskTypeRefuel, Status: constants.GroundTaskStatusInProgress, PlannedStart: start, Deadline: start.Add(time.Hour)}
	db.Create(&task)
	old := models.GroundResource{ResourceCode: "OLD", ResourceType: constants.GroundTaskTypeRefuel, AvailabilityStatus: constants.ResourceStatusBooked}
	target := models.GroundResource{ResourceCode: "NEW", ResourceType: constants.GroundTaskTypeRefuel, AvailabilityStatus: constants.ResourceStatusAvailable}
	db.Create(&old)
	db.Create(&target)
	booking := models.ResourceBooking{ResourceID: old.ID, TurnaroundID: flight.ID, TaskID: task.ID, StartTime: start, EndTime: start.Add(time.Hour), BookingStatus: constants.BookingStatusConfirmed}
	db.Create(&booking)

	result, err := rebook.Rebook(booking.ID, types.RebookRequest{NewResourceID: target.ID, Actor: "tester"})
	if err != nil {
		t.Fatalf("rebook failed: %v", err)
	}
	if result.ReleasedBooking.BookingStatus != constants.BookingStatusReleased {
		t.Fatalf("old status = %s, want RELEASED", result.ReleasedBooking.BookingStatus)
	}
	if result.NewBooking.BookingStatus != constants.BookingStatusConfirmed || result.NewBooking.ResourceID != target.ID {
		t.Fatalf("new booking wrong: %+v", result.NewBooking)
	}
	if result.ReleasedBooking.ReplacedByID == nil || *result.ReleasedBooking.ReplacedByID != result.NewBooking.ID {
		t.Fatal("old.replaced_by_id must point to new booking")
	}
	if result.NewBooking.ReplacedBookingID == nil || *result.NewBooking.ReplacedBookingID != booking.ID {
		t.Fatal("new.replaced_booking_id must point to old booking")
	}

	history, err := rebook.History(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 2 {
		t.Fatalf("history length = %d, want 2 (old + new)", len(history))
	}

	// Rebooking the already released row must fail.
	if _, err := rebook.Rebook(booking.ID, types.RebookRequest{NewResourceID: target.ID}); err == nil {
		t.Fatal("rebook of RELEASED booking must be rejected")
	}
}

// TestWindowsOverlap pins the conflict semantics: touching edges are allowed.
func TestWindowsOverlap(t *testing.T) {
	base := time.Date(2026, 9, 25, 9, 0, 0, 0, time.UTC)
	if services.WindowsOverlap(base, base.Add(30*time.Minute), base.Add(30*time.Minute), base.Add(time.Hour)) {
		t.Fatal("touching windows must not conflict")
	}
	if !services.WindowsOverlap(base, base.Add(time.Hour), base.Add(30*time.Minute), base.Add(2*time.Hour)) {
		t.Fatal("overlapping windows must conflict")
	}
}
