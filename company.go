package agilecrm

// ListCompanies ...
func (c *Client) ListCompanies() (ContactList, error) {
	out := ContactList{}
	_, err := c.get("POST", "api/contacts/companies/list", nil, nil, &out)

	return out, err
}

// CreateCompany ...
func (c *Client) CreateCompany(in Contact) (*Contact, error) {
	in.Type = TypeCompany
	return c.createContact(in)
}

// UpdateCompany ...
func (c *Client) UpdateCompanyProperties(id int64, in Contact) (*Contact, error) {
	in.ID = id
	return c._updateContact(in, "api/contacts/edit-properties")
}

// FindCompanyById ...
func (c *Client) FindCompanyById(id int) (*Contact, error) {
	return c.FindContactById(id)
}

// DeleteCompany ...
func (c *Client) DeleteCompany(id int) error {
	return c.DeleteContact(id)
}
