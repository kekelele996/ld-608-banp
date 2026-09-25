export const ERROR_MESSAGES = {
  AUTH_REQUIRED: "请先登录后再继续操作",
  RBAC_DENIED: "当前角色没有执行该动作的权限",
  VALIDATION_FAILED: "表单字段缺失或格式错误",
  RATE_LIMITED: "请求过于频繁，请稍后再试",
  TURNAROUND_NOT_FOUND: "航班过站不存在",
  BOOKING_NOT_FOUND: "资源预约不存在或已释放",
  RESOURCE_NOT_FOUND: "保障资源不存在",
  RELEASE_CHECK_FAILED: "放行复核未通过，航班状态与预约保持不变",
  RESOURCE_UNAVAILABLE: "目标资源不可用",
  REBIND_CONFLICT: "目标资源在该时段已被占用，换绑失败"
};
