package models

import (
	"time"

	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/api/pkg/web"

	"xorm.io/xorm"
)

// ProjectAssignee represents an assignment of a user to a project
type ProjectAssignee struct {
	ID        int64     `xorm:"bigint autoincr not null unique pk" json:"-"`
	ProjectID int64     `xorm:"bigint INDEX not null" json:"-"`
	UserID    int64     `xorm:"bigint INDEX not null" json:"user_id"`
	Created   time.Time `xorm:"created not null"`

	web.CRUDable    `xorm:"-" json:"-"`
	web.Permissions `xorm:"-" json:"-"`
}

// TableName makes a pretty table name
func (*ProjectAssignee) TableName() string {
	return "project_assignees"
}

// ProjectAssigneeWithUser is a helper type to deal with user joins
type ProjectAssigneeWithUser struct {
	ProjectID int64
	user.User `xorm:"extends"`
}

func getRawProjectAssigneesForProjects(s *xorm.Session, projectIDs []int64) (projectAssignees []*ProjectAssigneeWithUser, err error) {
	projectAssignees = []*ProjectAssigneeWithUser{}
	err = s.Table("project_assignees").
		Select("project_id, users.*").
		In("project_id", projectIDs).
		Join("INNER", "users", "project_assignees.user_id = users.id").
		Find(&projectAssignees)
	return
}

// Create or update a bunch of project assignees
func (p *Project) updateProjectAssignees(s *xorm.Session, assignees []*user.User, doer web.Auth) (err error) {

	// Load the current assignees
	currentAssignees, err := getRawProjectAssigneesForProjects(s, []int64{p.ID})
	if err != nil {
		return err
	}

	p.Assignees = make([]*user.User, 0, len(currentAssignees))
	for i := range currentAssignees {
		p.Assignees = append(p.Assignees, &currentAssignees[i].User)
	}

	// If we don't have any new assignees, delete everything right away.
	if len(assignees) == 0 && len(p.Assignees) > 0 {
		_, err = s.Where("project_id = ?", p.ID).
			Delete(&ProjectAssignee{})
		p.setProjectAssignees(assignees)
		return err
	}

	// If we didn't change anything (from 0 to zero) don't do anything.
	if len(assignees) == 0 && len(p.Assignees) == 0 {
		return nil
	}

	// Make a hashmap of the new assignees for easier comparison
	newAssignees := make(map[int64]*user.User, len(assignees))
	for _, newAssignee := range assignees {
		newAssignees[newAssignee.ID] = newAssignee
	}

	// Get old assignees to delete
	var found bool
	var assigneesToDelete []int64
	oldAssignees := make(map[int64]*user.User, len(p.Assignees))
	for _, oldAssignee := range p.Assignees {
		found = false
		if newAssignees[oldAssignee.ID] != nil {
			found = true // If a new assignee is already in the project with old assignees
		}

		// Put all assignees which are only on the old project to the trash
		if !found {
			assigneesToDelete = append(assigneesToDelete, oldAssignee.ID)
		}

		oldAssignees[oldAssignee.ID] = oldAssignee
	}

	// Delete all assignees not passed
	if len(assigneesToDelete) > 0 {
		_, err = s.In("user_id", assigneesToDelete).
			And("project_id = ?", p.ID).
			Delete(&ProjectAssignee{})
		if err != nil {
			return err
		}
	}

	// Loop through our users and add them
	for _, u := range assignees {
		// Check if the user is already assigned and assign him only if not
		if oldAssignees[u.ID] != nil {
			// continue outer loop
			continue
		}

		// Add the new assignee
		_, err = s.Insert(&ProjectAssignee{
			ProjectID: p.ID,
			UserID:    u.ID,
		})
		if err != nil {
			return err
		}
	}

	p.setProjectAssignees(assignees)

	return
}

// Small helper functions to set the new assignees in various places
func (p *Project) setProjectAssignees(assignees []*user.User) {
	if len(assignees) == 0 {
		p.Assignees = nil
		return
	}
	p.Assignees = assignees
}
