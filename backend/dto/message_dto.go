package dto

type Message struct {
	ID       int     `json:"id,omitempty"`
	MemberID int     `json:"member_id,omitempty"`
	Message  string  `json:"message,omitempty"`
	Member   *Member `json:"member,omitempty"`
}
