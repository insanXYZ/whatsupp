package dto

type Member struct {
	ID      int    `json:"id,omitempty"`
	GroupID string `json:"group_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	Role    string `json:"role,omitempty"`
	Group   *Group `json:"group,omitempty"`
	User    *User  `json:"user,omitempty"`
}
