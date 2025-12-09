package posting_queue

import "errors"

var (
	ErrPostNotFound = errors.New("post not found in queue")
)

