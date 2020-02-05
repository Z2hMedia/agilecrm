package agilecrm

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	TypeContact = "PERSON"
	TypeCompany = "COMPANY"
)

// ContactUser ...
type ContactUser struct {
	ID          int64  `json:"id"`
	Domain      string `json:"domain"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	PictureURL  string `json:"pic"`
	ScheduleID  string `json:"schedule_id"`
	CalendarURL string `json:"calendar_url"`
	CalURL      string `json:"calendarURL"`
}

type Viewed struct {
	ViewedTime int   `json:"viewed_time,omitempty"`
	ViewerId   int64 `json:"viewer_id,omitempty"`
}

// Contact ...
type Contact struct {
	ID                  int64  `json:"id,omitempty"`
	Type                string `json:"type,omitempty"`
	StarValue           int    `json:"star_value,omitempty"`
	LeadScore           int    `json:"lead_score,omitempty"`
	EntityType          string `json:"entity_type,omitempty"`
	ContactCompanyId    string `json:"contact_company_id,omitempty"`
	FormId              int64  `json:"formId,omitempty"`
	LastContacted       int    `json:"last_contacted,omitempty"`
	LastEmailed         int    `json:"last_emailed,omitempty"`
	LastCampaignEmailed int    `json:"last_campaign_emailed,omitempty"`
	LastCalled          int    `json:"last_called,omitempty"`
	CreatedAt           int    `json:"created_time,omitempty"`
	UpdatedAt           int    `json:"updated_time,omitempty"`

	Viewed *Viewed `json:"viewed,omitempty"`

	Tags []string `json:"tags,omitempty"`

	TagsWithTime []struct {
		Tag            string    `json:"tag,omitempty"`
		CreatedAt      time.Time `json:"created_at,omitempty"`
		AvailableCount int       `json:"available_count,omitempty"`
		EntityType     string    `json:"entity_type,omitempty"`
	} `json:"tags_with_time,omitempty"`

	Properties PropertyList `json:"properties,omitempty"`

	CampaignStatus    []interface{} `json:"campaignStatus,omitempty"`
	UnsubscribeStatus []interface{} `json:"unsubscribeStatus,omitempty"`
	EmailBounceStatus []interface{} `json:"emailBounceStatus,omitempty"`

	Owner *ContactUser `json:"owner,omitempty"`

	Cursor string `json:"cursor,omitempty"`
}

type ContactList []*Contact

// Email ...
func (c Contact) Email() string {
	for _, p := range c.Properties {
		if strings.ToLower(p.Name) == "email" {
			return p.Value
		}
	}
	return ""
}

// Cursor ...
func (cl ContactList) Cursor() string {
	if len(cl) <= 0 {
		return ""
	}
	c := cl[len(cl)-1]
	return c.Cursor
}

// ListContacts ...
func (c *Client) ListContacts(perPage int, cursor string) (ContactList, error) {
	out := ContactList{}

	params := map[string]string{}
	if perPage > 0 {
		params["page_size"] = fmt.Sprintf("%v", perPage)
	}
	if cursor != "" {
		params["cursor"] = cursor
	}

	st, err := c.get("GET", "api/contacts", nil, params, &out)

	if st == http.StatusOK && err == nil {
		// everything's fine, just return
		return out, nil
	}

	// if the status code is one of the defined status from the API,
	// return the proper error message
	switch st {
	case http.StatusNoContent:
		return ContactList{}, ErrNoContacts
	case http.StatusUnauthorized:
		return ContactList{}, ErrUnauthorized
	}

	return ContactList{}, err
}

// FindContactById ...
func (c *Client) FindContactById(id int) (*Contact, error) {
	r := fmt.Sprintf("api/contacts/%v", id)
	out := Contact{}
	err := c.findByID(r, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// FindContactByEmail ...
func (c *Client) FindContactByEmail(email string) (*Contact, error) {
	r := fmt.Sprintf("api/contacts/search/email/%v", email)

	out := &Contact{}
	st, err := c.get("GET", r, nil, nil, out)

	if st == http.StatusOK && err == nil {
		return out, nil
	}

	switch st {
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusBadRequest:
		return nil, ErrWrongFormat
	}

	return nil, err
}

// FindContactsByEmail ...
func (c *Client) FindContactsByEmail(emails []string) (map[string]*Contact, error) {
	tmp := make([]string, len(emails))
	for i := range emails {
		tmp[i] = fmt.Sprintf("\"%v\"", emails[i])
	}
	q := strings.Join(tmp, ",")

	vals := url.Values{}
	vals.Add("email_ids", fmt.Sprintf("[%v]", q))
	q = vals.Encode()

	// fmt.Printf("query: \n%v\n", out)

	out := map[string]*Contact{}

	req, err := c.postForm("POST", "api/contacts/search/email", strings.NewReader(q), nil)
	if err != nil {
		return out, err
	}

	cl := ContactList{}
	st, err := c.processResults(req, &cl)
	if st == http.StatusOK && err == nil {
		for i, v := range emails {
			out[v] = cl[i]
		}

		return out, nil
	}

	switch st {
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusBadRequest:
		return nil, ErrWrongFormat
	}
	return nil, err
}

// createContact ...
func (c *Client) createContact(in Contact) (*Contact, error) {
	st, err := c.send("POST", "api/contacts", nil, &in, &in)
	if st == http.StatusOK && err == nil {
		return &in, nil
	}

	switch st {
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusBadRequest:
		return nil, ErrWrongFormat
	case http.StatusNotAcceptable:
		return nil, ErrContactLimit
	}

	return nil, err
}

// CreateContact ...
func (c *Client) CreateContact(in Contact) (*Contact, error) {
	in.Type = TypeContact
	return c.createContact(in)
}

// _updateContact ...
func (c *Client) _updateContact(in Contact, route string) (*Contact, error) {
	st, err := c.send("PUT", route, nil, &in, &in)
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

// UpdateContactProperties can update the properties of a contact using this method. To update
// lead score, star value, or tags, use the appropriate method
func (c *Client) UpdateContactProperties(id int64, in Contact) (*Contact, error) {
	in.ID = id
	return c._updateContact(in, "api/contacts/edit-properties")
}

// UpdateLeadScore ...
func (c *Client) UpdateContactLeadScore(id, score int64) (*Contact, error) {
	ctc := Contact{ID: id, LeadScore: int(score)}
	return c._updateContact(ctc, "api/contacts/edit/lead-score")
}

// UpdateContactStarValue ...
func (c *Client) UpdateContactStarValue(id, star int64) (*Contact, error) {
	ctc := Contact{ID: id, StarValue: int(star)}
	return c._updateContact(ctc, "api/contacts/edit/add-star")
}

// UpdateContactTags ...
func (c *Client) UpdateContactTags(id int64, tags []string) (*Contact, error) {
	ctc := Contact{ID: id, Tags: tags}
	return c._updateContact(ctc, "api/contacts/edit/tags")
}

// DeleteContactTags ...
func (c *Client) DeleteContactTags(id int64, tags []string) (*Contact, error) {
	ctc := Contact{ID: id, Tags: tags}
	return c._updateContact(ctc, "api/contacts/delete/tags")
}

// DeleteContact ...
func (c *Client) DeleteContact(id int) error {
	r := fmt.Sprintf("api/contacts/%v", id)
	return c.delete(r)
}

// SearchContacts ...
func (c *Client) SearchContacts(query string) (ContactList, error) {
	return ContactList{}, nil
}
