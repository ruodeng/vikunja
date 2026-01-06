package models

import (
	"time"

	"code.vikunja.io/api/pkg/log"
	"code.vikunja.io/api/pkg/user"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type UserDashboardStat struct {
	User           *user.User `json:"user"`
	ActiveTasks    int64      `json:"active_tasks"`
	OverdueTasks   int64      `json:"overdue_tasks"`
	CompletedToday int64      `json:"completed_today"`
}

type DashboardStats struct {
	TotalTasks   int64               `json:"total_tasks"`
	OverdueTasks int64               `json:"overdue_tasks"`
	TotalUsers   int64               `json:"total_users"`
	UserStats    []UserDashboardStat `json:"user_stats"`
}

// GetDashboardStats returns statistics for the admin dashboard
func GetDashboardStats(s *xorm.Session) (*DashboardStats, error) {
	stats := &DashboardStats{}
	var err error

	// Global Stats
	// Use Where("id > 0") to avoid xorm using struct zero-values as filters
	stats.TotalTasks, err = s.Where("id > 0").Count(&Task{})
	if err != nil {
		log.Errorf("Error counting total tasks: %v", err)
		return nil, err
	}

	stats.OverdueTasks, err = s.Where("due_date < ? AND done = ?", time.Now(), false).Count(&Task{})
	if err != nil {
		log.Errorf("Error counting overdue tasks: %v", err)
		return nil, err
	}

	stats.TotalUsers, err = s.Where("id > 0").Count(&user.User{})
	if err != nil {
		log.Errorf("Error counting total users: %v", err)
		return nil, err
	}

	log.Infof("Dashboard stats: Tasks=%d, Overdue=%d, Users=%d", stats.TotalTasks, stats.OverdueTasks, stats.TotalUsers)

	// DEBUG: Verify raw count
	var rawUserCount int64
	_, err = s.SQL("SELECT count(*) FROM users").Get(&rawUserCount)
	if err != nil {
		log.Errorf("Raw SQL error: %v", err)
	} else {
		log.Infof("Raw SQL User Count: %d", rawUserCount)
	}

	// User Stats
	// Fetch all users to iterate
	var users []*user.User
	err = s.Where("id > 0").Find(&users)
	if err != nil {
		log.Errorf("Error fetching users: %v", err)
		return nil, err
	}

	log.Debugf("Found %d users for dashboard", len(users))

	stats.UserStats = make([]UserDashboardStat, 0, len(users))
	startOfDay := time.Now().Truncate(24 * time.Hour)

	for _, u := range users {
		stat := UserDashboardStat{User: u}

		// Active Tasks
		// Logic: (Assigned to User) OR (Created by User AND Not Assigned to Anyone Else)
		// Note regarding "Not Assigned to Anyone Else": If created by me and assigned to me, it's covered by "Assigned to User".
		// If created by me and assigned to Bob, it's "Assigned to Others" -> Exclude.
		// If created by me and Unassigned -> Include.
		// So we want: Match entry in task_assignees for user OR (created_by = user AND ID NOT IN task_assignees)

		stat.ActiveTasks, err = s.Where("done = ?", false).
			And(builder.Or(
				builder.In("id", builder.Select("task_id").From("task_assignees").Where(builder.Eq{"user_id": u.ID})),
				builder.And(
					builder.Eq{"created_by_id": u.ID},
					builder.NotIn("id", builder.Select("task_id").From("task_assignees")),
				),
			)).
			Count(&Task{})
		if err != nil {
			return nil, err
		}

		// Overdue Tasks
		// Same ownership logic, but checked against due_date
		stat.OverdueTasks, err = s.Where("done = ? AND due_date < ?", false, time.Now()).
			And(builder.Or(
				builder.In("id", builder.Select("task_id").From("task_assignees").Where(builder.Eq{"user_id": u.ID})),
				builder.And(
					builder.Eq{"created_by_id": u.ID},
					builder.NotIn("id", builder.Select("task_id").From("task_assignees")),
				),
			)).
			Count(&Task{})
		if err != nil {
			return nil, err
		}

		// Completed Today (Assigned, done, and done_at >= today)
		stat.CompletedToday, err = s.Table("tasks").
			Join("INNER", "task_assignees", "task_assignees.task_id = tasks.id").
			Where("task_assignees.user_id = ? AND tasks.done = ? AND tasks.done_at >= ?", u.ID, true, startOfDay).
			Count(&Task{})
		if err != nil {
			return nil, err
		}

		stats.UserStats = append(stats.UserStats, stat)
	}

	return stats, nil
}

type AdminUserDetail struct {
	User         *user.User `json:"user"`
	TotalTasks   int64      `json:"total_tasks"`
	ActiveTasks  int64      `json:"active_tasks"`
	OverdueTasks int64      `json:"overdue_tasks"`
	Tasks        []*Task    `json:"tasks"`
}

// GetAdminUserDetails returns detailed stats and tasks for a specific user
func GetAdminUserDetails(s *xorm.Session, userID int64) (*AdminUserDetail, error) {
	u, err := user.GetUserByID(s, userID)
	if err != nil {
		return nil, err
	}

	stats := &AdminUserDetail{User: u}

	// Ownership logic: (Assigned to User) OR (Created by User AND Not Assigned to Anyone Else)
	ownershipCond := builder.Or(
		builder.In("id", builder.Select("task_id").From("task_assignees").Where(builder.Eq{"user_id": userID})),
		builder.And(
			builder.Eq{"created_by_id": userID},
			builder.NotIn("id", builder.Select("task_id").From("task_assignees")),
		),
	)

	// Total Tasks (All time, done or undone)
	stats.TotalTasks, err = s.Where(ownershipCond).Count(&Task{})
	if err != nil {
		return nil, err
	}

	// Active Tasks (Undone)
	stats.ActiveTasks, err = s.Where("done = ?", false).And(ownershipCond).Count(&Task{})
	if err != nil {
		return nil, err
	}

	// Overdue Tasks (Undone and Due < Now)
	stats.OverdueTasks, err = s.Where("done = ? AND due_date < ?", false, time.Now()).
		And(ownershipCond).
		Count(&Task{})
	if err != nil {
		return nil, err
	}

	// Fetch Tasks (Limit 50, Sort by Due Date ASC)
	// We want to see overdues/urgent first
	stats.Tasks = make([]*Task, 0)
	err = s.Where(ownershipCond).
		Asc("done").     // Unfinished first
		Asc("due_date"). // Earliest due date first
		Desc("created"). // Then newest created
		Limit(50).
		Find(&stats.Tasks)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
