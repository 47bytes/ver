package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	git "gopkg.in/libgit2/git2go.v25"
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

func (versions Versions) latest() Version {
	var latest Version
	for i, v := range versions {
		if i == 0 {
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

func toVersion(s string) Version {
	tmp := strings.Split(s, ".")

	major, err := strconv.Atoi(tmp[0])
	CheckError(err)
	minor, err := strconv.Atoi(tmp[1])
	CheckError(err)

	tmp = strings.Split(tmp[2], "-")

	patch, err := strconv.Atoi(tmp[0])
	CheckError(err)

	var build string
	if len(tmp) == 2 {
		build = tmp[1]
	} else {
		build = ""
	}

	v := Version{
		Major: major,
		Minor: minor,
		Patch: patch,
		build: build,
	}

	return v
}

func cleanTag(t, prefix string) string {
	s := strings.Split(t, "/")
	tmp := s[len(s)-1]
	s = strings.Split(tmp, prefix)
	return s[len(s)-1]
}

func getVersionFromTag(s, prefix string) Version {
	tag := cleanTag(s, prefix)
	v := toVersion(tag)
	return v
}

var RootCmd = &cobra.Command{
	Use:   "ver",
	Short: "ver - simple git tag version incrementer",
	RunE:  rootCmdFn,
}

func rootCmdFn(cmd *cobra.Command, args []string) error {

	prefix, _ := cmd.Flags().GetString("prefix")
	PREFIX = prefix

	pwd, err := os.Getwd()
	if err != nil {
		return errors.New("Unable to get working directory. " + err.Error())
	}

	repo, err := git.OpenRepository(pwd)
	if err != nil {
		return errors.New("Directory doesn't appear to be a git repository. " + err.Error())
	}

	tags, err := repo.Tags.List()
	if err != nil {
		return errors.New("Tags could not be loaded. " + err.Error())
	}

	versions := Versions{}
	for _, tag := range tags {
		v := getVersionFromTag(tag, prefix)
		versions = append(versions, v)
	}

	latestVer := versions.latest()
	newVer := latestVer

	major, _ := cmd.Flags().GetBool("major")
	minor, _ := cmd.Flags().GetBool("minor")
	patch, _ := cmd.Flags().GetBool("patch")

	if major {
		newVer.Major += 1
	}
	if minor {
		newVer.Minor += 1
	}
	if patch {
		newVer.Patch += 1
	}

	fmt.Printf("old %#v\n", latestVer)
	fmt.Printf("new %#v\n", newVer)
	fmt.Println(newVer)

	confPath, _ := git.ConfigFindGlobal()
	if err != nil {
		return errors.New("Couldn't locate git config. " + err.Error())
	}

	simpleConf, err := git.NewConfig()
	if err != nil {
		return err
	}

	conf, err := git.OpenOndisk(simpleConf, confPath)
	if err != nil {
		return err
	}

	name, err := conf.LookupString("user.name")
	if err != nil {
		return errors.New("Couldn't find user.name git config key. " + err.Error())
	}
	email, err := conf.LookupString("user.email")
	if err != nil {
		return errors.New("Couldn't find user.email git config key. " + err.Error())
	}

	user := &git.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	oid := head.Target()

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		return err
	}

	blubb, _ := repo.Tags.Create(newVer.String(), commit, user, newVer.String())

	fmt.Printf("%#v", blubb)

	// rep.Tags.Create(newVer.String()	, commit *git.Commit, tagger *git.Signature, message string)

	// rep, err := git.PlainOpen(pwd)
	// if err != nil {
	// 	return errors.New("Directory doesn't appear to be a git repository. " + err.Error())
	// }

	// tagRefs, err := rep.Tags()
	// if err != nil {
	// 	return err
	// }

	// versions := Versions{}
	// err = tagRefs.ForEach(func(t *plumbing.Reference) error {
	// 	x := t.Strings()
	// 	v := getVersionFromTag(x[0], prefix)
	// 	v.tag = t
	// 	versions = append(versions, v)
	// 	return nil
	// })
	// CheckError(err)

	// latestVer := versions.latest()
	// newVer := latestVer

	// major, _ := cmd.Flags().GetBool("major")
	// minor, _ := cmd.Flags().GetBool("minor")
	// patch, _ := cmd.Flags().GetBool("patch")

	// if major {
	// 	newVer.Major += 1
	// }
	// if minor {
	// 	newVer.Minor += 1
	// }
	// if patch {
	// 	newVer.Patch += 1
	// }

	// fmt.Printf("old %#v\n", latestVer)
	// fmt.Printf("new %#v\n", newVer)

	// tag, err := rep.TagObject(newVer.tag.Hash())
	// rep.Worktree().Commit(msg string, opts *git.CommitOptions)

	// tags, err := rep.TagObjects()
	// err = tags.ForEach(func(t *object.Tag) error {
	// 	fmt.Println(t.Name)
	// 	return nil
	// })

	return nil

}

func init() {
	RootCmd.Flags().String("prefix", "v", "Prefix for git tag")
	RootCmd.Flags().BoolP("major", "M", false, "Increase major version number")
	RootCmd.Flags().BoolP("minor", "m", false, "Increase minor version number")
	RootCmd.Flags().BoolP("patch", "p", false, "Increase patch version number")
	// RootCmd.Flags().String("branch", "master", "Which branch should be tagged")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println("Unexpected error. " + err.Error())
		return
	}
}
