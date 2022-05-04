package message

const (
	OnlyJSON             = "Header 'Content-Type' must be 'application/json'."
	NotMeetAllRequest    = "This request does not meet all of the required elements."
	CannotUseAPI         = "You cannot use this API."
	CannotAddOutSideUser = "You cannot add outside user."

	TokenNotFound            = "The token is required."
	AuthorizationTypeInvalid = "This authorization type is not supported."
	WrongToken               = "This token is wrong."
	AlreadyExisted           = "This id's record is already existed."
	WrongIDorPassword        = "This login_id or this password is wrong."
	ValidJWT                 = "This JWT token is correct."
	InformedJWT              = "This JWT token is informed."
	ExpiredOrNotValidYetJWT  = "This JWT token is expired or not valid yet."
	InvalidJWT               = "This JWT token is invalid."

	DeviceNotFound = "This device is not found."
	UserNotFound   = "This user is not found."
	RaceNotFound   = "This race is not found."
	GroupNotFound  = "This group is not found."

	NotSupportWebSocket = "This client is not support WebSocket."
	NoUserIDContain     = "The user id must be contain."
	InvalidPointID      = "This point id is invalid."
)
