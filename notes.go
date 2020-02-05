package agilecrm

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Note struct {
	ID          int64    `json:"id,omitempty"`
	Subject     string   `json:"subject"`
	Description string   `json:"description"`
	ContactIDs  []string `json:"contact_ids,omitempty"`
	DealIDs     []string `json:"deal_ids,omitempty"`
	CreatedTime int      `json:"created_time,omitempty"`
	EntityType  string   `json:"entity_type,omitempty"`

	// Count is only used by the API when getting the notes for a contact
	Count    *int        `json:"count,omitempty"`
	Contacts ContactList `json:"contacts,omitempty"`
}

type NoteList []*Note

// CreateNote ...
func (c *Client) CreateNote(in Note) (*Note, error) {
	in.Count = nil
	_, err := c.send("POST", "api/notes", nil, in, &in)
	if err != nil {
		return nil, err
	}
	return &in, nil
}

// AddNoteToContact ...
func (c *Client) AddNoteToContact(email string, in Note) (*Note, error) {
	in.Count = nil

	bits, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	vals := url.Values{}
	vals.Add("email", email)
	vals.Add("note", string(bits))
	q := vals.Encode()

	req, err := c.postForm("POST", "api/contacts/email/note/add", strings.NewReader(q), nil)
	if err != nil {
		return nil, err
	}

	_, err = c.processResults(req, &in)
	if err != nil {
		return nil, err
	}

	return &in, nil
}

// GetContactNotes ...
func (c *Client) GetContactNotes(id int) (NoteList, error) {
	r := fmt.Sprintf("api/contacts/%v/notes", id)
	out := NoteList{}

	_, err := c.get("GET", r, nil, nil, &out)
	return out, err
}

// DeleteContactNote ...
func (c *Client) DeleteContactNote(contact_id, note_id int) error {
	r := fmt.Sprintf("api/contacts/%v/notes/%v", contact_id, note_id)
	return c.delete(r)
}

// CreateDealNote ...
func (c *Client) CreateDealNote(id int, in Note) (*Note, error) {
	in.Count = nil
	did := fmt.Sprintf("%v", id)
	in.DealIDs = append(in.DealIDs, did)
	in.ContactIDs = []string{}

	_, err := c.send("PUT", "api/opportunity/deals/notes", nil, in, &in)
	if err != nil {
		return nil, err
	}
	return &in, nil
}

// GetDealNotes ...
func (c *Client) GetDealNotes(id int64) (NoteList, error) {
	r := fmt.Sprintf("api/opportunity/%v/notes", id)
	out := NoteList{}
	_, err := c.get("GET", r, nil, nil, &out)
	return out, err
}

// TODO: delete notes from specific deal
