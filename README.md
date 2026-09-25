# 航空地勤周转保障平台

面向机场地勤团队的航班过站保障、资源调度、异常延误和任务签收系统。

## 快速启动

```bash
cp .env.example .env && docker compose up -d
```

## 访问地址或 CLI 示例

前端：<http://localhost:20108>

后端健康检查：<http://localhost:21108/health>


## 本地开发方式

- 前端：`cd frontend && npm install && npm run dev`（开发服务器已把 `/api` 代理到本地后端）
- 后端：`cd backend && go run main.go`（默认监听 `:3000`，可用 `PORT` 环境变量覆盖），接口统一挂在 `/api`。

## 过站放行协同流程

调度员在「过站放行」页（`/release`）逐项核对航班放行条件，避免放走条件不足的航班：

1. **航班详情汇总**：`GET /api/flight-turnaround/:id/release-summary` 汇总未完成任务、未关闭延误、预约时间冲突三类条件，并给出放行结论（`RELEASABLE` / `BLOCKED`）与逐项阻碍说明。
2. **提交放行复核**：`POST /api/flight-turnaround/:id/release` 由后端重新核对三类条件；任一不满足即返回 `409 RELEASE_CHECK_FAILED`，并在 `blockers` 中指出具体航班、任务或资源，航班状态与原预约保持不变；全部通过才把航班置为 `READY`。
3. **资源换绑**：`GET /api/resource-booking/:id/rebind-options` 列出可换用的同类可用资源（`AVAILABLE` 且时段无冲突），`POST /api/resource-booking/:id/rebind` 完成换绑——原预约转为 `RELEASED` 并生成 `CONFIRMED` 新预约，两步在同一锁内原子完成；`GET /api/resource-booking?turnaround_id=` 返回含已释放在内的全部预约，历史记录可查前后两次预约。

种子数据内置三种演示场景：航班 CA1234 同时命中三类阻碍；航班 MU5678 仅有预约冲突（换绑 BAG-02 → BAG-01 后即可放行）；航班 CZ9012 可直接放行。

## 技术栈


## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | React 18 + TypeScript + Vite + Ant Design + Redux Toolkit |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 部署 | Docker Compose |

## 项目目录结构

```text
frontend/src/api, stores, types, constants, constructors, components/common, hooks, pages, router, utils, mocks
backend/src/routes, controllers, services, models, repositories, middlewares, constants, constructors, utils, types, config
```

## 环境变量说明

- `COMPOSE_PROJECT_NAME`: Compose 项目名，默认 `ground-turn`
- `FRONTEND_PORT`: 前端端口，默认 `20108`
- `BACKEND_PORT`: 后端端口，默认 `21108`
- `DB_PORT`: 数据库宿主机端口
- `DB_USER/DB_PASSWORD/DB_NAME`: 本地数据库凭据

## Docker 部署说明

- 根 Compose 文件不写 `version`，顶层 `name: ground-turn`。
- 容器名均使用 `${COMPOSE_PROJECT_NAME:-ground-turn}` 前缀。
- 数据库使用命名卷，避免绑定中文路径。
- 常见问题：端口占用时修改 `.env` 中端口后重启；需要重置数据时执行 `docker compose down -v`。

## 枚举/常量出现位置清单

- GroundTaskType: constants/GroundTaskType、types/GroundTaskType、constructors、logTemplates、errorMessages、筛选器、展示组件/控制器均有引用。
- TurnaroundStatus: constants/TurnaroundStatus、types/TurnaroundStatus、constructors、logTemplates、errorMessages、筛选器、展示组件/控制器均有引用。
- ResourceStatus: constants/ResourceStatus、types/ResourceStatus、constructors、logTemplates、errorMessages、筛选器、展示组件/控制器均有引用。
- BookingStatus（PENDING/CONFIRMED/RELEASED/CANCELLED）: 后端 `constants/BookingStatus.go`（含 `IsActiveBookingStatus`）、`services/ResourceBooking.go`（冲突检测与换绑）、`repositories/ResourceBooking.go`；前端 `constants/BookingStatus.ts`、`types/BookingStatus.ts`、`constants/statusText.ts`、`pages/ReleasePage.tsx`（预约历史）、`constructors/TurnaroundReleaseConstructor.ts`。
- GroundTaskStatus（PENDING/IN_PROGRESS/BLOCKED/COMPLETED）: 后端 `constants/GroundTaskStatus.go`、`services/FlightTurnaround.go`（放行复核）、仓储种子；前端 `constants/GroundTaskStatus.ts`、`types/GroundTaskStatus.ts`、`pages/ReleasePage.tsx`（未完成任务列表）、`constructors/TurnaroundReleaseConstructor.ts`。
- ReleaseConclusion（RELEASABLE/BLOCKED）: 后端 `constants/ReleaseConclusion.go`、`constructors/TurnaroundRelease.go`、`services/FlightTurnaround.go`；前端 `constants/ReleaseConclusion.ts`、`types/TurnaroundRelease.ts`、`pages/ReleasePage.tsx`（放行结论横幅）。

## 为什么会牵一发动全身

实体字段、枚举、日志模板、错误消息、构造器、筛选器和展示组件被刻意拆散到多个目录；修改一个状态值通常需要同步类型、构造器、服务、控制器、store、页面、README 与数据库种子。

## License

MIT
