package def

import "errors"

var (
	KeyEmpty                error = errors.New("key is empty")
	EmptyKeyCannotBeDeleted error = errors.New("can't delete empty key")
	StartOrEndKeyEmpty      error = errors.New("start or end key is empty")
	StartKeyUnknown         error = errors.New("can't scan from last without knowing startKey")
	ResultsNotFound         error = errors.New("results not found")
)
