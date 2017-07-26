package collection

import (
	"github.com/kuzzleio/sdk-go/types"
)

/*
  Instantiates a new Room object.
*/
func (dc Collection) Room(options types.RoomOptions) types.Room {
	room := types.Room{}

	scopes := map[string]bool{types.SCOPE_ALL: true, types.SCOPE_IN: true, types.SCOPE_OUT: true, types.SCOPE_NONE: true, types.SCOPE_UNKNOWN: true}
	states := map[string]bool{types.STATE_ALL: true, types.STATE_PENDING: true, types.STATE_DONE: true}
	user := map[string]bool{types.USER_ALL: true, types.USER_IN: true, types.USER_OUT: true, types.USER_NONE: true}

	if options != nil {
		scope := options.GetScope()
		if scopes[scope] {
			room.Scope = scope
		}
		state := options.GetState()
		if states[state] {
			room.State = state
		}
		u := options.GetUser()
		if user[u] {
			room.User = u
		}
		room.SubscribeToSelf = options.GetSubscribeToSelf()
	}

	return room
}
