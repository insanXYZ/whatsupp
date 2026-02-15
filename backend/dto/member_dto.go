package dto

type Member struct {
	ID    int    `json:"id,omitempty"`
	Role  string `json:"role,omitempty"`
	Group *Group `json:"group,omitempty"`
	User  *User  `json:"user,omitempty"`
}
