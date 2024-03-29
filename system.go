package directus

type Role struct {
	ID          string `json:"id,omitempty"`
	Icon        Icon   `json:"icon,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	AdminAccess bool `json:"admin_access"`
	AppAccess   bool `json:"app_access"`

	Users []string `json:"users,omitempty"`
}

type User struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`

	Provider           string `json:"provider,omitempty"`
	ExternalIdentifier string `json:"external_identifier,omitempty"`
}

type Icon string
