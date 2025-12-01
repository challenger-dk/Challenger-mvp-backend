package appError

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

// Generic Errors
var (
	ErrUnknownResource = errors.New("unknown resource")
	ErrServerError     = errors.New("internal server error")
)

// Authentication and Authorization Errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrAuthHeaderMissing  = errors.New("authorization header missing")
	ErrInvalidAuthHeader  = errors.New("invalid authorization header format")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidFriendship  = errors.New("invalid friendship")
	ErrSameUser           = errors.New("same user")
	ErrUserBlocked        = errors.New("you have been blocked by this user")
)

var (
	ErrUserExists    = errors.New("user with this email already exists")
	ErrInvalidSport  = errors.New("invalid sport name")
	ErrSportNotFound = errors.New("sport not found")
)

// Invitation Errors
var (
	ErrInviteSameUser            = errors.New("inviter and invitee cannot be the same user")
	ErrInvitationPending         = errors.New("invitation is already pending")
	ErrInvitationAccepted        = errors.New("user has already accepted this invitation")
	ErrInvitationDeclined        = errors.New("user has already declined this invitation")
	ErrInvitationProcessed       = errors.New("invitation already processed")
	ErrUnhandledInvitationStatus = errors.New("unhandled invitation status")
)

var (
	ErrMissingIdParam = errors.New("missing id parameter")
	ErrIdBelowOne     = errors.New("id parameter must be greater than 0")
)

// ErrorMap groups specific errors by the HTTP status code they should return.
// Validation errors are handled separately in the error handler.
var errorMap = map[int][]error{
	http.StatusNotFound: {
		gorm.ErrRecordNotFound,
		ErrSportNotFound,
	},
	http.StatusUnauthorized: {
		ErrInvalidCredentials,
		ErrInvalidToken,
		ErrUnauthorized,
		ErrUserNotFound,
	},
	http.StatusForbidden: {
		ErrUserBlocked,
	},
	http.StatusConflict: {
		ErrUserExists,
		ErrInvitationPending,
		ErrInvitationAccepted,
		ErrInvitationDeclined,
		ErrInvitationProcessed,
		ErrInviteSameUser,
	},
	http.StatusBadRequest: {
		ErrInvalidSport,
		ErrAuthHeaderMissing,
		ErrInvalidAuthHeader,
		ErrMissingIdParam,
		ErrIdBelowOne,
		ErrInvalidFriendship,
		ErrSameUser,
	},
	http.StatusInternalServerError: {
		ErrUnhandledInvitationStatus,
		ErrUnknownResource,
		ErrServerError,
	},
}
