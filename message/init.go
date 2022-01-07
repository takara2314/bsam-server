package message

const (
	OnlyJSON             = "Header 'Content-Type' must be 'application/json'."
	NotMeetAllRequest    = "This request does not meet all of the required elements."
	CannotUseAPI         = "You cannot use this API."
	CannotAddOutSideUser = "You cannot add outside user."

	TokenNotFound            = "The token is required."
	AuthorizationTypeInvalid = "The authorization type is not supported."
	WrongToken               = "This token is wrong."
)
