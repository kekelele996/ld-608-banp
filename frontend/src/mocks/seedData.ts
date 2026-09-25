// 本地种子数据与后端 repositories/store.go 保持一致：
// 航班 1 覆盖三类放行阻碍（未完成任务、未关闭延误、预约冲突），
// 航班 2 仅有预约冲突（演示换绑后放行），航班 3 可直接放行。
export const mockData = {
  "flightTurnaround": [
    {
      "id": 1,
      "flight_no": "CA1234",
      "aircraft_reg": "B-1234",
      "stand_no": "101",
      "arrival_time": "2026-09-25T08:00:00Z",
      "departure_time": "2026-09-25T09:30:00Z",
      "turnaround_status": "IN_SERVICE",
      "delay_reason": "前段天气流控"
    },
    {
      "id": 2,
      "flight_no": "MU5678",
      "aircraft_reg": "B-5678",
      "stand_no": "102",
      "arrival_time": "2026-09-25T08:30:00Z",
      "departure_time": "2026-09-25T10:00:00Z",
      "turnaround_status": "IN_SERVICE",
      "delay_reason": ""
    },
    {
      "id": 3,
      "flight_no": "CZ9012",
      "aircraft_reg": "B-9012",
      "stand_no": "103",
      "arrival_time": "2026-09-25T09:00:00Z",
      "departure_time": "2026-09-25T10:30:00Z",
      "turnaround_status": "IN_SERVICE",
      "delay_reason": ""
    }
  ],
  "groundTask": [
    {
      "id": 1,
      "turnaround_id": 1,
      "task_type": "CLEANING",
      "team_id": 11,
      "planned_start": "2026-09-25T08:10:00Z",
      "deadline": "2026-09-25T08:50:00Z",
      "actual_finish": "",
      "status": "IN_PROGRESS",
      "blocker_note": "客舱遗留物清点中"
    },
    {
      "id": 2,
      "turnaround_id": 1,
      "task_type": "REFUEL",
      "team_id": 12,
      "planned_start": "2026-09-25T08:20:00Z",
      "deadline": "2026-09-25T08:55:00Z",
      "actual_finish": "2026-09-25T08:50:00Z",
      "status": "COMPLETED",
      "blocker_note": ""
    },
    {
      "id": 3,
      "turnaround_id": 1,
      "task_type": "BAGGAGE",
      "team_id": 13,
      "planned_start": "2026-09-25T09:00:00Z",
      "deadline": "2026-09-25T09:40:00Z",
      "actual_finish": "2026-09-25T09:35:00Z",
      "status": "COMPLETED",
      "blocker_note": ""
    },
    {
      "id": 4,
      "turnaround_id": 2,
      "task_type": "BAGGAGE",
      "team_id": 13,
      "planned_start": "2026-09-25T08:40:00Z",
      "deadline": "2026-09-25T09:20:00Z",
      "actual_finish": "2026-09-25T09:15:00Z",
      "status": "COMPLETED",
      "blocker_note": ""
    },
    {
      "id": 5,
      "turnaround_id": 2,
      "task_type": "CATERING",
      "team_id": 14,
      "planned_start": "2026-09-25T08:45:00Z",
      "deadline": "2026-09-25T09:15:00Z",
      "actual_finish": "2026-09-25T09:10:00Z",
      "status": "COMPLETED",
      "blocker_note": ""
    },
    {
      "id": 6,
      "turnaround_id": 3,
      "task_type": "CLEANING",
      "team_id": 11,
      "planned_start": "2026-09-25T09:05:00Z",
      "deadline": "2026-09-25T09:40:00Z",
      "actual_finish": "2026-09-25T09:35:00Z",
      "status": "COMPLETED",
      "blocker_note": ""
    },
    {
      "id": 7,
      "turnaround_id": 3,
      "task_type": "WATER_SERVICE",
      "team_id": 15,
      "planned_start": "2026-09-25T09:10:00Z",
      "deadline": "2026-09-25T09:35:00Z",
      "actual_finish": "2026-09-25T09:30:00Z",
      "status": "COMPLETED",
      "blocker_note": ""
    }
  ],
  "groundResource": [
    {
      "id": 1,
      "resource_code": "CLN-01",
      "resource_type": "CLEANING",
      "location": "T1 西区",
      "availability_status": "AVAILABLE",
      "maintenance_due_at": "2026-12-01T00:00:00Z",
      "owner_team": "保洁一组"
    },
    {
      "id": 2,
      "resource_code": "CAT-01",
      "resource_type": "CATERING",
      "location": "T1 配餐区",
      "availability_status": "AVAILABLE",
      "maintenance_due_at": "2026-11-15T00:00:00Z",
      "owner_team": "配餐组"
    },
    {
      "id": 3,
      "resource_code": "BAG-01",
      "resource_type": "BAGGAGE",
      "location": "T1 行李区",
      "availability_status": "AVAILABLE",
      "maintenance_due_at": "2026-10-20T00:00:00Z",
      "owner_team": "行李组"
    },
    {
      "id": 4,
      "resource_code": "BAG-02",
      "resource_type": "BAGGAGE",
      "location": "T1 行李区",
      "availability_status": "AVAILABLE",
      "maintenance_due_at": "2026-10-20T00:00:00Z",
      "owner_team": "行李组"
    },
    {
      "id": 5,
      "resource_code": "REF-01",
      "resource_type": "REFUEL",
      "location": "机坪加油点",
      "availability_status": "AVAILABLE",
      "maintenance_due_at": "2027-01-10T00:00:00Z",
      "owner_team": "油料组"
    },
    {
      "id": 6,
      "resource_code": "WAT-01",
      "resource_type": "WATER_SERVICE",
      "location": "T1 东区",
      "availability_status": "AVAILABLE",
      "maintenance_due_at": "2026-09-30T00:00:00Z",
      "owner_team": "勤务组"
    },
    {
      "id": 7,
      "resource_code": "PUSH-01",
      "resource_type": "PUSHBACK",
      "location": "机坪牵引区",
      "availability_status": "AVAILABLE",
      "maintenance_due_at": "2026-12-20T00:00:00Z",
      "owner_team": "牵引组"
    },
    {
      "id": 8,
      "resource_code": "BAG-03",
      "resource_type": "BAGGAGE",
      "location": "T1 行李区",
      "availability_status": "MAINTENANCE",
      "maintenance_due_at": "2026-09-26T08:00:00Z",
      "owner_team": "行李组"
    },
    {
      "id": 9,
      "resource_code": "CAT-02",
      "resource_type": "CATERING",
      "location": "T2 配餐区",
      "availability_status": "OFFLINE",
      "maintenance_due_at": "2026-10-01T00:00:00Z",
      "owner_team": "配餐组"
    }
  ],
  "resourceBooking": [
    {
      "id": 1,
      "resource_id": 1,
      "turnaround_id": 1,
      "task_id": 1,
      "start_time": "2026-09-25T08:10:00Z",
      "end_time": "2026-09-25T08:50:00Z",
      "booking_status": "CONFIRMED",
      "conflict_reason": ""
    },
    {
      "id": 2,
      "resource_id": 5,
      "turnaround_id": 1,
      "task_id": 2,
      "start_time": "2026-09-25T08:20:00Z",
      "end_time": "2026-09-25T08:55:00Z",
      "booking_status": "CONFIRMED",
      "conflict_reason": ""
    },
    {
      "id": 3,
      "resource_id": 4,
      "turnaround_id": 1,
      "task_id": 3,
      "start_time": "2026-09-25T09:00:00Z",
      "end_time": "2026-09-25T09:40:00Z",
      "booking_status": "CONFIRMED",
      "conflict_reason": ""
    },
    {
      "id": 4,
      "resource_id": 4,
      "turnaround_id": 2,
      "task_id": 4,
      "start_time": "2026-09-25T08:40:00Z",
      "end_time": "2026-09-25T09:20:00Z",
      "booking_status": "CONFIRMED",
      "conflict_reason": ""
    },
    {
      "id": 5,
      "resource_id": 2,
      "turnaround_id": 2,
      "task_id": 5,
      "start_time": "2026-09-25T08:45:00Z",
      "end_time": "2026-09-25T09:15:00Z",
      "booking_status": "CONFIRMED",
      "conflict_reason": ""
    },
    {
      "id": 6,
      "resource_id": 1,
      "turnaround_id": 3,
      "task_id": 6,
      "start_time": "2026-09-25T09:05:00Z",
      "end_time": "2026-09-25T09:40:00Z",
      "booking_status": "CONFIRMED",
      "conflict_reason": ""
    },
    {
      "id": 7,
      "resource_id": 6,
      "turnaround_id": 3,
      "task_id": 7,
      "start_time": "2026-09-25T09:10:00Z",
      "end_time": "2026-09-25T09:35:00Z",
      "booking_status": "CONFIRMED",
      "conflict_reason": ""
    }
  ],
  "delayEvent": [
    {
      "id": 1,
      "turnaround_id": 1,
      "delay_type": "WEATHER",
      "minutes": 25,
      "root_cause": "前段航路雷雨流控",
      "responsibility_team": "运控中心",
      "resolved_at": ""
    },
    {
      "id": 2,
      "turnaround_id": 1,
      "delay_type": "CLEANING",
      "minutes": 10,
      "root_cause": "客舱遗留物等待认领",
      "responsibility_team": "保洁一组",
      "resolved_at": "2026-09-25T08:45:00Z"
    },
    {
      "id": 3,
      "turnaround_id": 2,
      "delay_type": "BAGGAGE",
      "minutes": 10,
      "root_cause": "行李分拣系统短暂故障",
      "responsibility_team": "行李组",
      "resolved_at": "2026-09-25T09:05:00Z"
    }
  ]
} as const;
