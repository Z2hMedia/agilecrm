package agilecrm

import (
	"fmt"
	"time"
)

type Event struct {
	ID             int64  `json:"id,omitempty"`
	CreatedTime    int    `json:"created_time,omitempty"`
	AllDay         bool   `json:"all_day,omitempty"`
	Title          string `json:"title,omitempty"`
	Color          string `json:"color,omitempty"`
	Start          int    `json:"start,omitempty"`
	End            int    `json:"end,omitempty"`
	IsEventStarred bool   `json:"is_event_starred,omitempty"`

	Contacts ContactList `json:"contacts,omitempty"`
}

type EventUpsert struct {
	ID             *int64   `json:"id,omitempty"`
	CreatedTime    int      `json:"created_time,omitempty"`
	AllDay         bool     `json:"all_day,omitempty"`
	Title          string   `json:"title,omitempty"`
	Color          string   `json:"color,omitempty"`
	Start          int      `json:"start,omitempty"`
	End            int      `json:"end,omitempty"`
	IsEventStarred bool     `json:"is_event_starred,omitempty"`
	Contacts       []string `json:"contacts,omitempty"`
}

type EventList []Event

// ListEvents ...
func (c *Client) ListEvents(start, end time.Time) (EventList, error) {
	s := start.Unix()
	e := end.Unix()

	params := map[string]string{
		"start": fmt.Sprintf("%v", s),
		"end":   fmt.Sprintf("%v", e),
	}

	out := EventList{}
	_, err := c.get("GET", "api/events", nil, params, &out)

	return out, err
}

// GetContactEvents ...
func (c *Client) GetContactEvents(id int64) (EventList, error) {
	r := fmt.Sprintf("api/contacts/%v/events/sort", id)
	out := EventList{}
	_, err := c.get("GET", r, nil, nil, &out)
	return out, err
}

// CreateEvent ...
func (c *Client) CreateEvent(e EventUpsert) (*Event, error) {
	out := &Event{}
	_, err := c.send("POST", "api/events", nil, e, out)
	return out, err
}

// UpdateEvent ...
func (c *Client) UpdateEvent(id int64, e EventUpsert) (*Event, error) {
	out := &Event{}
	e.ID = &id
	_, err := c.send("PUT", "api/events", nil, e, out)
	return out, err
}

// DeleteEvent ...
func (c *Client) DeleteEvent(id int64) error {
	r := fmt.Sprintf("api/events/%v", id)
	return c.delete(r)
}
