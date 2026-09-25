package types

// RebindRequest 资源换绑请求：把冲突预约换到另一台可用资源。
type RebindRequest struct {
	ResourceID int `json:"resource_id"`
}
