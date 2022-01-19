package code

type Code string

const (
	OK Code = "ok"

	// BadRequest      Code = "bad_request"
	// InvalidArgument Code = "invalid_argument"
	// Unauthorized    Code = "unauthorized"
	// NotFound        Code = "not_found"
	// Conflict        Code = "conflict"

	InvalidQuery    Code = "invalid_query"
	InvalidArgument Code = "invalid_argument"
	NotFound        Code = "not_fount"

	// server error
	UUID     Code = "uuid"
	JSON     Code = "json"
	Slack    Code = "slack"
	Internal Code = "internal"

	Unknown Code = "unknown"
)
