package message

const (
	SUCCESS_GET_ME             = "User profile retrieved successfully."
	SUCCESS_UPDATE_ME          = "User profile updated successfully."
	SUCCESS_REGISTER           = "Registration completed successfully."
	SUCCESS_LOGIN              = "Login successful."
	SUCCESS_LIST_GROUPS        = "Groups retrieved successfully."
	SUCCESS_SEND_FILES         = "Files sent successfully."
	SUCCESS_LOGOUT             = "Logout successful."
	SUCCESS_LIST_RECENT_GROUPS = "Recent groups retrieved successfully."
	SUCCESS_GET_MESSAGES       = "Messages retrieved successfully."

	ERR_BIND_REQ           = "Invalid request payload."
	ERR_GET_ME             = "Failed to retrieve user profile."
	ERR_UPDATE_ME          = "Failed to update user profile."
	ERR_REGISTER           = "Registration failed."
	ERR_LOGIN              = "Invalid email or password."
	ERR_LIST_GROUPS        = "Failed to retrieve groups."
	ERR_RETRIEVE_FILES     = "Failed to retrieve files."
	ERR_SEND_FILES         = "Failed to send files."
	ERR_LOGOUT             = "Failed to logout."
	ERR_LIST_RECENT_GROUPS = "Failed to retrieve recent groups."
	ERR_GET_MESSAGES       = "Failed to retrieve messages"

	WS_SUCCESS_SYNC_CHAT = "Chat synchronized successfully."
	WS_ERR_SYNC_CHAT     = "Failed to synchronize chat."
)
