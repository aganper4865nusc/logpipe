package aggregate

import "errors"

var (
	errNilFunc          = errors.New("aggregate: flush function must not be nil")
	errNoFlushCondition = errors.New("aggregate: at least one of maxSize or interval must be positive")
)
