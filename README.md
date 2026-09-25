# 航空地勤周转保障平台

面向机场地勤团队的航班过站保障、资源调度、异常延误和任务签收系统。

## 过站放行协同流程（核心）

调度员在「过站放行」页选中航班后，系统汇总三类放行条件并给出结论：

1. **未完成任务**：状态非 `COMPLETED` 的地勤任务（可当场签收完成）；
2. **未关闭延误**：`resolved_at` 为空的延误事件（可当场关闭）；
3. **预约时间冲突**：同一资源上与其他航班预约的时间窗重叠（可给冲突任务**换用可用资源**）。

关键保证：

- 提交放行时**后端在一个数据库事务内重新核对**三类条件，任一项不满足即返回 `409 RELEASE_REJECTED`，并在 `violations` 中逐条指出具体航班、任务（`task_id`）、延误（`delay_id`）或资源/预约（`resource_id`、`booking_id`）及原因；**航班状态与原预约保持不变**。
- 资源换绑（`POST /api/bookings/:id/rebook`）同样在事务内完成：原预约置为 `RELEASED`（`replaced_by_id` 指向新预约），生成新的 `CONFIRMED` 预约（`replaced_booking_id` 回指原预约）。换绑前校验资源类型匹配、`AVAILABLE` 状态与时间窗空闲，任一不满足则整笔回滚。前后两次预约可经 `GET /api/bookings/task/:taskId/history` 双向追溯。

| 接口 | 说明 |
|---|---|
| `GET /api/turnarounds/:id/release-check` | 三类条件汇总 + `can_release` 结论 + 每个冲突的可换用资源 |
| `POST /api/turnarounds/:id/release` | 提交放行（后端事务内复核，失败 409 且不变更） |
| `POST /api/bookings/:bookingId/rebook` | 冲突预约换绑资源（旧 RELEASED / 新 CONFIRMED，双向链接） |
| `GET /api/bookings/:bookingId/candidates` | 可换用资源候选（同类型、可用、时间窗空闲） |
| `GET /api/bookings/task/:taskId/history` | 某任务的前后预约历史链 |
| `PATCH /api/ground-tasks/:id/status` | 任务派发/签收/阻塞 |
| `POST /api/delay-events/:id/resolve` | 关闭延误 |

## 快速启动

```bash
cp .env.example .env && docker compose up -d
```

## 访问地址或 CLI 示例

前端：<http://localhost:20108>

后端健康检查：<http://localhost:21108/health>


## 本地开发方式

- 前端：`cd frontend && npm install && npm run dev`（Vite 将 `/api` 代理到 `http://localhost:21108`，见 `vite.config.ts`）
- 后端：进入 `backend` 后按技术栈运行开发命令，接口统一挂在 `/api`。
  - 连本地 MySQL（与容器一致）：`go run .`，通过 `DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME` 配置；
  - 无 MySQL 时可用 SQLite 快速冒烟（需要 CGO）：`DB_DRIVER=sqlite DB_SQLITE_PATH=ground-turn.db go run -tags sqlite .`，首次启动自动建表并写入种子数据。
  - 服务层测试：`go test -tags sqlite ./src/services/`（覆盖“驳回放行状态/预约不变”与“换绑前后双向链接”）。


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
- GroundTaskStatus（PENDING/IN_PROGRESS/BLOCKED/COMPLETED）：`backend/src/constants/GroundTaskStatus.go`、`frontend/src/constants/GroundTaskStatus.ts`、放行核对服务（`services/ReleaseService.go`）、任务状态更新（`services/TaskService.go`）、前端 `utils/formatters.ts`、`StatusBadge`、未完成任务区块。
- BookingStatus（CONFIRMED/CONFLICT/RELEASED）：`backend/src/constants/BookingStatus.go`、`frontend/src/constants/BookingStatus.ts`、冲突检测（`services/ConflictService.go`）、换绑（`services/RebookService.go`）、仓储预约链查询、前端 `utils/formatters.ts`、`StatusBadge`、换绑弹窗与预约历史。
- 放行错误码：`backend/src/constants/errorCodes.go`（`TASK_UNFINISHED`/`DELAY_NOT_CLOSED`/`BOOKING_CONFLICT`/`RELEASE_REJECTED` 等）与 `errorMessages.go` 配套，由 service 包装、controller 转为 HTTP 响应、前端 `api/client.ts` 的 `ApiError.violations` 消费。

## 为什么会牵一发动全身

实体字段、枚举、日志模板、错误消息、构造器、筛选器和展示组件被刻意拆散到多个目录；修改一个状态值通常需要同步类型、构造器、服务、控制器、store、页面、README 与数据库种子。

## License

MIT
