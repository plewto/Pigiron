package main

import "fmt"

type version struct {
	major int
	minor int
	revision int
	level string
}

var VERSION = version{0, 2, 1, "Beta"}

func (v *version) String() string {
	mj, mn, rv, lev  :=  v.major, v.minor, v.revision, v.level
	return fmt.Sprintf("Version %d.%d.%d %s", mj, mn, rv, lev)
}

