package gateway

import "errors"

var (
	ErrEventChannelClosed = errors.New(`event channel is closed`)
)
