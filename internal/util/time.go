package util

import (
    "time"
)

func MustLocation(name string) *time.Location {
    loc, err := time.LoadLocation(name)
    if err != nil { return time.FixedZone(name, -8*3600) }
    return loc
}
