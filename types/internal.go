package types

import "github.com/blang/semver/v4"

// This is the database and app info
type DBInfo struct {
	// Version is limited to the Major version of the semver model
	Version semver.Version
}
