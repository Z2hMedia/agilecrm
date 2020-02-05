package agilecrm

import "fmt"

type Document struct {
	ID           int64  `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	DummyName    string `json:"dummy_name,omitempty"`
	UploadedTime int    `json:"uploaded_time,omitempty"`
	Extension    string `json:"extension,omitempty"`
	DocType      string `json:"doc_type,omitempty"`
	Text         string `json:"text,omitempty"`
	TemplateType string `json:"template_type,omitempty"`
	Size         int    `json:"size,omitempty"`
	NetworkType  string `json:"network_type,omitempty"`
	URL          string `json:"url,omitempty"`
	EntityType   string `json:"entity_type,omitempty"`

	ContactIds []string `json:"contact_ids,omitempty"`
	CaseIds    []string `json:"case_ids,omitempty"`
	DealIds    []string `json:"deal_ids,omitempty"`

	Update int `json:"update,omitempty"`

	Owner *ContactUser `json:"owner,omitempty"`

	Contacts        ContactList `json:"contacts,omitempty"`
	Deals           DealList    `json:"deals,omitempty"`
	RelatedContacts ContactList `json:"related_contacts,omitempty"`
}

type DocumentList []Document

type UpsertDoc struct {
	ID          *int64   `json:"id,omitempty"`
	Extension   string   `json:"extension,omitempty"`
	DocType     string   `json:"doc_type,omitempty"`
	Name        string   `json:"name,omitempty"`
	URL         string   `json:"url,omitempty"`
	Size        int      `json:"size,omitempty"`
	NetworkType string   `json:"network_type,omitempty"`
	ContactIds  []string `json:"contact_ids,omitempty"`
	DealIds     []string `json:"deal_ids,omitempty"`
}

// GetContactDocuments ...
func (c *Client) GetContactDocuments(id int64) (DocumentList, error) {
	r := fmt.Sprintf("api/documents/contact/%v/docs", id)

	out := DocumentList{}
	_, err := c.get("GET", r, nil, nil, &out)

	return out, err
}

// CreateContactDocument ...
func (c *Client) CreateDocument(doc UpsertDoc) (*Document, error) {
	doc.ID = nil
	out := &Document{}
	_, err := c.send("POST", "api/documents", nil, doc, out)
	return out, err
}

// UpdateDocument ...
func (c *Client) UpdateDocument(id int64, doc UpsertDoc) (*Document, error) {
	doc.ID = &id
	out := &Document{}
	_, err := c.send("PUT", "api/documents", nil, doc, out)
	return out, err
}
