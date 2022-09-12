package errors

const (
	userNotCreated    = "user not created"
	failJSONBind      = "failed to bind JSON"
	invalidID         = "invalid ID"
	userNotFound      = "user not found"
	userNotInContext  = "user not found in context. Has it been set in the middleware?"
	userNotAuthorized = "user not authorized"
	userNotModified   = "user not modified"
	userNotDeleted    = "user not deleted"
)

type UserNotCreated struct{}
type FailJSONBind struct{}
type InvalidID struct{}
type UserNotFound struct{}
type UserNotInContext struct{}
type UserNotAuthorized struct{}
type UserNotModified struct{}
type UserNotDeleted struct{}

func (e UserNotCreated) Error() string {
	return userNotCreated
}
func (e FailJSONBind) Error() string {
	return failJSONBind
}
func (e InvalidID) Error() string {
	return invalidID
}
func (e UserNotFound) Error() string {
	return userNotFound
}
func (e UserNotInContext) Error() string {
	return userNotInContext
}
func (e UserNotAuthorized) Error() string {
	return userNotAuthorized
}
func (e UserNotModified) Error() string {
	return userNotModified
}
func (e UserNotDeleted) Error() string {
	return userNotDeleted
}
