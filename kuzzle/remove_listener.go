package kuzzle

// RemoveListener removes a listener
func (k Kuzzle) RemoveListener(event int) {
	k.socket.RemoveListener(event)
}