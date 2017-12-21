package model

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var PREFIX string

type Version struct {
	Major int
	Minor int
	Patch int
	build string
}

func (v Version) String() string {
	if v.build != "" {
		return fmt.Sprintf("%s%d.%d.%d-%s", PREFIX, v.Major, v.Minor, v.Patch, v.build)
	}
	return fmt.Sprintf("%s%d.%d.%d", PREFIX, v.Major, v.Minor, v.Patch)
}

type Versions []Version

func (versions Versions) Latest() Version {
	var latest Version
	for i, v := range versions {
		if i == 0 {
			latest = v
		}

		if v.Major > latest.Major {
			latest = v
		}
		if (v.Major == latest.Major) && (v.Minor > latest.Minor) {
			latest = v
		}
		if (v.Major == latest.Major) && (v.Minor == latest.Minor) && (v.Patch > latest.Patch) {
			latest = v
		}

		if (v.Major >= latest.Major) && (v.Minor >= latest.Minor) && (v.Patch >= latest.Patch) {
			latest = v
		}
	}

	return latest
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func toVersion(s string) (*Version, error) {
	tmp := strings.Split(s, ".")

	major, err := strconv.Atoi(tmp[0])
	if err != nil {
		return nil, errors.New("Major has to be an int. " + err.Error())
	}

	minor, err := strconv.Atoi(tmp[1])
	if err != nil {
		return nil, errors.New("Minor has to be an int. " + err.Error())
	}

	tmp = strings.Split(tmp[2], "-")

	patch, err := strconv.Atoi(tmp[0])
	if err != nil {
		return nil, errors.New("Patch has to be an int. " + err.Error())
	}

	var build string
	if len(tmp) == 2 {
		build = tmp[1]
	} else {
		build = ""
	}

	v := &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
		build: build,
	}

	return v, nil
}

func cleanTag(t string) string {
	s := strings.Split(t, "/")
	tmp := s[len(s)-1]
	s = strings.Split(tmp, PREFIX)
	return s[len(s)-1]
}

func GetVersionFromTag(s string) (*Version, error) {
	tag := cleanTag(s)
	v, err := toVersion(tag)
	if err != nil {
		return nil, err
	}
	return v, nil
}
