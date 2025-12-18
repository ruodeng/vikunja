// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package models

import (
	"code.vikunja.io/api/pkg/web"
	"xorm.io/xorm"
)

// CanDelete checks if the user can delete an task
func (t *Task) CanDelete(s *xorm.Session, a web.Auth) (bool, error) {
	// If the user is an assignee, they can delete the task
	isAssignee, err := t.IsAssignee(s, a)
	if err != nil {
		return false, err
	}
	if isAssignee {
		return true, nil
	}

	// If the user is the creator, they can delete the task
	dbTask, err := GetTaskByIDSimple(s, t.ID)
	if err != nil {
		return false, err
	}
	if dbTask.CreatedByID == a.GetID() {
		p := &Project{ID: dbTask.ProjectID}
		if can, _ := p.CanExecute(s, a); can {
			return true, nil
		}
	}

	return t.canDoTask(s, a)
}

func (t *Task) IsAssignee(s *xorm.Session, a web.Auth) (bool, error) {
	return s.Where("task_id = ? AND user_id = ?", t.ID, a.GetID()).Exist(&TaskAssginee{})
}

// CanUpdate determines if a user has the permission to update a project task
func (t *Task) CanUpdate(s *xorm.Session, a web.Auth) (bool, error) {
	// If the user is an assignee, they can update the task
	isAssignee, err := t.IsAssignee(s, a)
	if err != nil {
		return false, err
	}
	if isAssignee {
		return true, nil
	}

	// If the user is the creator, they can update the task
	dbTask, err := GetTaskByIDSimple(s, t.ID)
	if err != nil {
		return false, err
	}
	if dbTask.CreatedByID == a.GetID() {
		p := &Project{ID: dbTask.ProjectID}
		if can, _ := p.CanExecute(s, a); can {
			return true, nil
		}
	}

	return t.CanWrite(s, a)
}

// CanCreate determines if a user has the permission to create a project task
func (t *Task) CanCreate(s *xorm.Session, a web.Auth) (bool, error) {
	// A user can create a task if he has write access to its project
	l := &Project{ID: t.ProjectID}
	if can, _ := l.CanWrite(s, a); can {
		return true, nil
	}

	// If the user only has executor (=restricted) access, they can only create a subtask
	// and only if they are assigned to the parent task
	if can, _ := l.CanExecute(s, a); can {

		if parentTasks, ok := t.RelatedTasks[RelationKindSubtask]; ok && len(parentTasks) > 0 {
			for _, parent := range parentTasks {
				// We need to check if the user is an assignee of the parent task
				// or if they have update rights on it.
				pt := &Task{ID: parent.ID}
				if canUpdate, _ := pt.CanUpdate(s, a); canUpdate {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// CanRead determines if a user can read a task
func (t *Task) CanRead(s *xorm.Session, a web.Auth) (canRead bool, maxPermission int, err error) {
	expand := t.Expand
	// Get the task, error out if it doesn't exist
	*t, err = GetTaskByIDSimple(s, t.ID)
	if err != nil {
		return
	}

	t.Expand = expand

	// A user can read a task if it has access to the project
	l := &Project{ID: t.ProjectID}
	return l.CanRead(s, a)
}

// CanWrite checks if a user has write access to a task
func (t *Task) CanWrite(s *xorm.Session, a web.Auth) (canWrite bool, err error) {
	return t.canDoTask(s, a)
}

// CanExecute checks if a user has restricted access to a task
func (t *Task) CanExecute(s *xorm.Session, a web.Auth) (bool, error) {
	// Get the task
	ot, err := GetTaskByIDSimple(s, t.ID)
	if err != nil {
		return false, err
	}

	// Check if we're moving the task into a different project to check if the user has sufficient permissions for that on the new project
	if t.ProjectID != 0 && t.ProjectID != ot.ProjectID {
		newProject := &Project{ID: t.ProjectID}
		can, err := newProject.CanExecute(s, a)
		if err != nil {
			return false, err
		}
		if !can {
			return false, ErrGenericForbidden{}
		}
	}

	// A user can do a task if it has write acces to its project
	l := &Project{ID: ot.ProjectID}
	return l.CanExecute(s, a)
}

// Helper function to check if a user can do stuff on a project task
func (t *Task) canDoTask(s *xorm.Session, a web.Auth) (bool, error) {
	// Get the task
	ot, err := GetTaskByIDSimple(s, t.ID)
	if err != nil {
		return false, err
	}

	// Check if we're moving the task into a different project to check if the user has sufficient permissions for that on the new project
	if t.ProjectID != 0 && t.ProjectID != ot.ProjectID {
		newProject := &Project{ID: t.ProjectID}
		can, err := newProject.CanWrite(s, a)
		if err != nil {
			return false, err
		}
		if !can {
			return false, ErrGenericForbidden{}
		}
	}

	// A user can do a task if it has write acces to its project
	l := &Project{ID: ot.ProjectID}
	return l.CanWrite(s, a)
}
