package message

const (
	OnlyJSON             = "Header 'Content-Type' must be 'application/json'."
	NotMeetAllRequest    = "This request does not meet all of the required elements."
	CannotUseAPI         = "You cannot use this API."
	CannotAddOutSideUser = "You cannot add outside user."

	TokenNotFound            = "The token is required."
	AuthorizationTypeInvalid = "This authorization type is not supported."
	WrongToken               = "This token is wrong."
	DeviceNotFound           = "This device is not found."
	AlreadyExisted           = "This id's record is already existed."
	RaceNotFound             = "This race is not found."
)
