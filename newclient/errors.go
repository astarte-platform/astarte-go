package newclient

import "errors"

var (
	ErrConflictingUrls           error = errors.New("Conflicting URLs provided")
	ErrNoUrlsProvided            error = errors.New("No Astarte URL(s) provided")
	ErrNoPrivateKeyProvided      error = errors.New("No Astarte private key provided")
	ErrDifferentStatusCode       error = errors.New("Received unexpected status code")
	ErrRealmNameNotProvided      error = errors.New("Realm name was not provided")
	ErrRealmPublicKeyNotProvided error = errors.New("Realm public key was not provided")
	ErrTooManyReplicationFactors error = errors.New("Can't have both replication factor and datacenter replication factors")
	ErrNegativeReplicationFactor error = errors.New("Replication factor must be a strictly positive integer")
)
