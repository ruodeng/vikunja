# Executor Role Implementation

## Overview
This document summarizes the changes made to implement the "Executor" permission role in Vikunja. This role is designed to give users limited ability to interact with tasks (specifically updating assigned tasks and creating subtasks) without full write access to the project.

## Permission Model
The `PermissionExecutor` (value: 1) sits between `PermissionRead` (0) and `PermissionWrite` (2).

### Privileges
- **Read Access**: Full read access to projects and tasks.
- **Update Access**: Can update tasks *only if* they are assigned to the task.
- **Create Access**: 
    - Cannot create standalone tasks.
    - Can create subtasks *only if* they are assigned to the parent task.
- **Delete Access**: Can delete tasks *only if* they are assigned to the task (or created it).

## Implementation Details

### Backend (`pkg/models`)
- **`permissions.go`**: Added `PermissionExecutor` constant.
- **`tasks_permissions.go`**: 
    - Updated `CanCreate` to restrict Executors to only creating subtasks for tasks they are assigned to.
    - (`CanUpdate` relies on `IsAssignee` check which allows update if assigned).
- **`executor_permissions_test.go`**: Added unit tests covering:
    - `Executor cannot create standalone task`
    - `Executor cannot update unassigned task`
    - `Executor updated assigned task`
    - `Executor can create subtask for assigned task`
    - `Executor cannot create subtask for unassigned task`

### Frontend (`frontend/src`)
- **`constants/permissions.ts`**: Added `EXECUTOR` constant.
- **`components/tasks/partials/RelatedTasks.vue`**: Updated to show only "subtask" relation option for users with Executor permissions.

## Verification
Unit tests have been added and passed (`go test -v ./pkg/models -run TestExecutorPermissions`).
