package common

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrInvalidVersionFormat = errors.New("invalid version format")

type Version struct {
	Major string
	Minor string
	Patch string
}

func NewVersion(args ...string) Version {
	expanded := append(args, []string{"*", "*", "*"}...)
	return Version{
		Major: expanded[0],
		Minor: expanded[1],
		Patch: expanded[2],
	}
}

func (v Version) String() string {
	return fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, v.Patch)
}

func (v Version) ge(a, b string) bool {
	return (len(a) > len(b)) || (len(a) == len(b) && a >= b)
}

func (v Version) equal(a, b string) bool {
	return a == b
}

func (v Version) GreaterEqual(ov Version) bool {
	return v.ge(v.Major, ov.Major) ||
		(v.equal(v.Major, ov.Major) && v.ge(v.Minor, ov.Minor)) ||
		(v.equal(v.Major, ov.Major) && v.ge(v.Minor, ov.Minor) && v.ge(v.Patch, ov.Patch))
}

func (v *Version) Equal(ov Version) bool {
	return v.equal(v.Major, ov.Major) && v.equal(v.Minor, ov.Minor) && v.equal(v.Patch, ov.Patch)
}

func (v *Version) Compatiable(ov Version) bool {
	return v.equal(v.Major, ov.Major) && (v.ge(v.Minor, ov.Minor) || v.Minor == "*")
}

var (
	versionPattern string = `^(\d+\.)?(\d+\.)?(\*|\d+)$`
	versionRegex   *regexp.Regexp
)

func init() {
	versionRegex = regexp.MustCompile("^" + versionPattern + "$")
}

func ParseVersion(s string) (Version, error) {
	matches := versionRegex.FindStringSubmatch(s)
	if matches == nil {
		return Version{}, ErrInvalidVersionFormat
	}
	splitted := strings.Split(matches[0], ".")
	return NewVersion(splitted...), nil
}
