package models

// GroundResource 保障资源：被 ResourceBooking 预约。
type GroundResource struct {
	ID                 int    `json:"id"`
	ResourceCode       string `json:"resource_code"`
	ResourceType       string `json:"resource_type"`
	Location           string `json:"location"`
	AvailabilityStatus string `json:"availability_status"`
	MaintenanceDueAt   string `json:"maintenance_due_at"`
	OwnerTeam          string `json:"owner_team"`
}
