package agilecrm

import (
	"fmt"
	"net/http"
)

type TaskType string

const (
	TaskCall      TaskType = "CALL"
	TaskEmail     TaskType = "EMAIL"
	TaskFollowUp  TaskType = "FOLLOW_UP"
	TaskMeeting   TaskType = "MEETING"
	TaskMilestone TaskType = "MILESTONE"
	TaskSend      TaskType = "SEND"
	TaskTweet     TaskType = "TWEET"
	TaskOther     TaskType = "OTHER"
)

type TaskPriority string

const (
	TaskPriorityHigh TaskPriority = "HIGH"
	TaskPriorityNorm TaskPriority = "NORMAL"
	TaskPriorityLow  TaskPriority = "LOW"
)

type TaskStatus string

const (
	TaskStatusUnstarted TaskStatus = "YET_TO_START"
	TaskStatusStarted   TaskStatus = "IN_PROGRESS"
	TaskStatusDone      TaskStatus = "COMPLETED"
)

type Task struct {
	ID           int    `json:"id,omitempty"`
	Type         string `json:"type,omitempty"`
	PriorityType string `json:"priority_type,omitempty"`
	Due          int    `json:"due,omitempty"`
	CreatedTime  int    `json:"created_time,omitempty"`
	IsComplete   bool   `json:"is_complete,omitempty"`
	Subject      string `json:"subject,omitempty"`
	Progress     int    `json:"progress,omitempty"`
	Status       string `json:"status,omitempty"`
	Owner        string `json:"owner,omitempty"`

	Contacts ContactList `json:"contacts,omitempty"`
	Notes    NoteList    `json:"notes,omitempty"`

	TaskOwner  *TaskOwner `json:"task_owner,omitempty"`
	EntityType string     `json:"entity_type,omitempty"`
}

type TaskCreate struct {
	ID             *int64       `json:"id,omitempty"`
	Progress       int          `json:"progress"`
	IsComplete     bool         `json:"is_complete"`
	Subject        string       `json:"subject"`
	Type           TaskType     `json:"type"`
	Due            int64        `json:"due"`
	TaskEndingTime string       `json:"task_ending_time"`
	OwnerID        int64        `json:"owner_id"`
	PriorityType   TaskPriority `json:"priority_type"`
	Status         TaskStatus   `json:"status"`
	Description    string       `json:"description"`

	ContactIDs []string `json:"contacts"`
	NoteIDs    []string `json:"notes"`
	DealIDs    []string `json:"deal_ids"`
}

type TaskList []Task

type TaskOwner struct {
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Email          string `json:"email,omitempty"`
	Domain         string `json:"domain,omitempty"`
	IsAdmin        bool   `json:"is_admin,omitempty"`
	IsAccountOwner bool   `json:"is_account_owner,omitempty"`
	IsDisabled     bool   `json:"is_disabled,omitempty"`
}

// ListTasks ...
func (c *Client) ListTasks() (TaskList, error) {
	out := TaskList{}
	_, err := c.get("GET", "api/tasks", nil, nil, &out)
	return out, err
}

// GetContactTasks ...
func (c *Client) GetContactTasks(id int) (TaskList, error) {
	route := fmt.Sprintf("api/contacts/%v/tasks", id)
	out := TaskList{}
	st, err := c.get("GET", route, nil, nil, &out)
	if st == http.StatusUnauthorized {
		return TaskList{}, ErrUnauthorized
	}
	return out, err
}

// GetPendingTasks ...
func (c *Client) GetPendingTasks(numDays int) (TaskList, error) {
	if numDays < 1 {
		numDays = 1
	}
	r := fmt.Sprintf("api/tasks/pending/%v", numDays)
	out := TaskList{}
	_, err := c.get("GET", r, nil, nil, &out)
	return out, err
}

// TODO: get tasks based on filters

// GetTaskByID ...
func (c *Client) GetTaskByID(id int64) (*Task, error) {
	r := fmt.Sprintf("api/tasks/%v", id)
	out := &Task{}
	err := c.findByID(r, out)
	return out, err
}

// CreateTask ...
func (c *Client) CreateTask(in TaskCreate) (*Task, error) {
	in.ID = nil
	out := &Task{}
	_, err := c.send("POST", "api/tasks", nil, in, out)
	return out, err
}

// UpdateTask ...
func (c *Client) UpdateTask(id int64, in TaskCreate) (*Task, error) {
	in.ID = &id
	out := &Task{}
	_, err := c.send("PUT", "api/tasks/partial-update", nil, in, out)
	return out, err
}

// DeleteTask ...
func (c *Client) DeleteTask(id int64) error {
	r := fmt.Sprintf("api/tasks/%v", id)
	return c.delete(r)
}
