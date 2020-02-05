package agilecrm

import (
	"fmt"
	"net/http"
)

type Deal struct {
	ID            int64       `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	Description   string      `json:"description,omitempty"`
	ExpectedValue float64     `json:"expected_value,omitempty"`
	PipelineID    int64       `json:"pipeline_id,omitempty"`
	Milestone     string      `json:"milestone,omitempty"`
	Probabilty    int         `json:"probabilty,omitempty"`
	CloseDate     int         `json:"close_date,omitempty"`
	CreatedTime   int         `json:"created_time,omitempty"`
	OwnerID       string      `json:"owner_id,omitempty"`
	Prefs         string      `json:"prefs,omitempty"`
	Contacts      ContactList `json:"contacts,omitempty"`
	ContactIds    []string    `json:"contact_ids,omitempty"`
	Cursor        string      `json:"cursor,omitempty"`
}

type DealList []Deal

// Cursor ...
func (dl DealList) Cursor() string {
	if len(dl) <= 0 {
		return ""
	}
	d := dl[len(dl)-1]
	return d.Cursor
}

// ListDeals ...
func (c *Client) ListDeals(perPage int, cursor string) (DealList, error) {
	out := DealList{}

	params := map[string]string{}
	if perPage > 0 {
		params["page_size"] = fmt.Sprintf("%v", perPage)
	}
	if cursor != "" {
		params["cursor"] = cursor
	}

	st, err := c.get("GET", "api/opportunity", nil, params, &out)

	if st == http.StatusOK && err == nil {
		// everything's fine, just return
		return out, nil
	}

	// if the status code is one of the defined status from the API,
	// return the proper error message
	switch st {
	case http.StatusNoContent:
		return out, ErrNoContacts
	case http.StatusUnauthorized:
		return out, ErrUnauthorized
	}

	return out, nil
}

// FindDealByID ...
func (c *Client) FindDealByID(id int) (*Deal, error) {
	r := fmt.Sprintf("api/opportunity/%v", id)
	out := Deal{}
	err := c.findByID(r, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDeal ...
func (c *Client) CreateDeal(in Deal) (*Deal, error) {
	st, err := c.send("POST", "api/opportunity", nil, in, &in)

	if st == http.StatusOK && err == nil {
		return &in, nil
	}

	switch st {
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusBadRequest:
		return nil, ErrWrongFormat
	}

	return nil, err
}

// UpdateDeal ...
func (c *Client) UpdateDeal(id int64, in Deal) (*Deal, error) {
	in.ID = id
	st, err := c.send("PUT", "api/opportunity/partial-update", nil, in, &in)
	if st == http.StatusOK && err == nil {
		return &in, nil
	}

	switch st {
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusBadRequest:
		return nil, ErrWrongFormat
	}
	return nil, err
}

// DeleteDeal ...
func (c *Client) DeleteDeal(id int) error {
	r := fmt.Sprintf("api/opportunity/%v", id)
	return c.delete(r)
}

// TODO: create deal for a contact using email

// TODO: bulk delete

// TODO: get deals from default track grouped by milestone

// TODO: get deals for specific track grouped by milestone

// TODO: get deals for specific track

// TODO: get deals for specific contact

// TODO: get deals for current user

// TODO: remove deal contacts
