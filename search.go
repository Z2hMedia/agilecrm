package agilecrm

import (
	"encoding/json"
	"net/url"
	"strings"
)

type filterCondition string

const (
	FilterEqual filterCondition = "EQUALS"
	FilterOn    filterCondition = "ON"
)

type filterRule struct {
	Left      string          `json:"LHS"`
	Condition filterCondition `json:"CONDITION"`
	Right     string          `json:"RHS"`
}

type filterJson struct {
	Rules       []filterRule `json:"rules"`
	OrRules     []filterRule `json:"or_rules"`
	ContactType ContactType  `json:"contact_type"`
}

// addRule ...
func (fj *filterJson) addRule(l, r string, c filterCondition) {
	fj.Rules = append(fj.Rules, filterRule{l, c, r})
}

// addOrRule ...
func (fj *filterJson) addOrRule(l, r string, c filterCondition) {
	fj.OrRules = append(fj.OrRules, filterRule{l, c, r})
}

// setType ...
func (fj *filterJson) setType(ct ContactType) {
	fj.ContactType = ct
}

// dynamicFilter ...
func (c *Client) dynamicFilter(f filterJson, out interface{}) error {
	bits, err := json.Marshal(f)
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Add("page_size", "10")
	v.Add("global_sort_key", "-created_time")
	v.Add("filterJson", string(bits))
	q := v.Encode()

	req, err := c.postForm("POST", "api/filters/filter/dynamic-filter", strings.NewReader(q), nil)
	if err != nil {
		return err
	}

	_, err = c.processResults(req, out)
	return err
}

// FindContactsByTag ...
func (c *Client) FindContactsByTag(tag string) (ContactList, error) {
	fj := filterJson{}
	fj.addRule("tags", tag, FilterEqual)
	fj.setType(TypeContact)

	out := ContactList{}
	err := c.dynamicFilter(fj, &out)
	return out, err
}

// FindCompaniesByTag ...
func (c *Client) FindCompaniesByTag(tag string) (ContactList, error) {
	fj := filterJson{}
	fj.addRule("tags", tag, FilterEqual)
	fj.setType(TypeCompany)

	out := ContactList{}
	err := c.dynamicFilter(fj, &out)
	return out, err
}

// FindDealsByTag ...
func (c *Client) FindDealsByTag(tag string) (DealList, error) {
	fj := filterJson{}
	fj.addRule("tags", tag, FilterEqual)
	fj.setType(TypeDeal)

	out := DealList{}
	err := c.dynamicFilter(fj, &out)
	return out, err
}
