package dto

type Group struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Bio  string `json:"bio,omitempty"`
	Type string `json:"type,omitempty"`
}

type SearchGroupRequest struct {
	Name string `query:"name"`
}

type SearchGroupResponse struct {
	Type      string  `json:"type"`
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Image     string  `json:"image"`
	Bio       *string `json:"bio,omitempty"`
	GroupType *string `json:"group_type,omitempty"`
	GroupID   *int    `json:"group_id"`
}

type LoadRecentGroup = SearchGroupResponse
