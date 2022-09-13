package errors

const (
	failJSONBind         = "failed to bind JSON"
	invalidID            = "invalid ID"
	userNotCreated       = "user not created"
	userNotFound         = "user not found"
	userNotInContext     = "user not found in context. Has it been set in the middleware?"
	userNotAuthorized    = "user not authorized"
	userNotModified      = "user not modified"
	userNotDeleted       = "user not deleted"
	userAlreadyInRoom    = "user already in room"
	roomNotCreated       = "room not created"
	roomNotFound         = "room not found"
	roomNotInContext     = "room not found in context. Has it been set in the middleware?"
	roomNotModified      = "room not modified"
	roomNotDeleted       = "room not deleted"
	streamNotCreated     = "stream not created"
	ownerCantKickHimself = "cannot kick owner from room"
)

type FailJSONBind struct{}
type InvalidID struct{}
type UserNotCreated struct{}
type UserNotFound struct{}
type UserNotInContext struct{}
type UserNotAuthorized struct{}
type UserNotModified struct{}
type UserNotDeleted struct{}
type RoomNotCreated struct{}
type RoomNotFound struct{}
type RoomNotInContext struct{}
type RoomNotModified struct{}
type RoomNotDeleted struct{}
type UserAlreadyInRoom struct{}
type StreamNotCreated struct{}
type OwnerCantKickHimself struct{}

func (e FailJSONBind) Error() string {
	return failJSONBind
}
func (e InvalidID) Error() string {
	return invalidID
}
func (e UserNotCreated) Error() string {
	return userNotCreated
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
func (e UserAlreadyInRoom) Error() string {
	return userAlreadyInRoom
}
func (e RoomNotCreated) Error() string {
	return roomNotCreated
}
func (e RoomNotFound) Error() string {
	return roomNotFound
}
func (e RoomNotInContext) Error() string {
	return roomNotInContext
}
func (e RoomNotModified) Error() string {
	return roomNotModified
}
func (e RoomNotDeleted) Error() string {
	return roomNotDeleted
}
func (e StreamNotCreated) Error() string {
	return streamNotCreated
}
func (e OwnerCantKickHimself) Error() string {
	return ownerCantKickHimself
}
