package agilecrm

const (
	TypeSystem = "SYSTEM"

	SubtypeWork = "work"
)

// Property ...
type Property struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Value   string `json:"value,omitempty"`
	Subtype string `json:"subtype,omitempty"`
}

type PropertyList []Property

// Find ...
func (pl PropertyList) Find(i string) (string, string) {
	for _, v := range pl {
		if v.Name == i {
			return v.Type, v.Value
		}
	}

	return "", ""
}
