// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provisioning

// UserHook receives lifecycle notifications from user management (async enqueue).
type UserHook interface {
	OnUserUpsert(userID uint)
	OnUserDelete(userID uint)
}

// OnUserUpsert implements UserHook.
func (c *Coordinator) OnUserUpsert(userID uint) {
	c.EnqueueUserSync(userID)
}

// OnUserDelete implements UserHook.
func (c *Coordinator) OnUserDelete(userID uint) {
	c.EnqueueUserDelete(userID)
}
