package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vvvvv/ver/pkg/ver"
	git "gopkg.in/libgit2/git2go.v25"
)

var (
	BUILD_VERSION = ""
	BUILD_HASH    = ""
	BUILD_DATE    = ""
)

var RootCmd = &cobra.Command{
	Use:     "ver",
	Short:   "ver - a simple git tag semver version incrementer",
	Long:    "ver increments semver style git tags",
	Example: "$ ver -m -p\n Tag `v0.2.1` created successfully\n 27c1f1234188aa11585334726f8721d9a35038eb",
	RunE:    rootCmdFn,
}

func rootCmdFn(cmd *cobra.Command, args []string) error {
	ver.Prefix, _ = cmd.Flags().GetString("prefix")

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

	setToVersion, _ := cmd.Flags().GetString("set")
	if setToVersion != "" {
		// has no version prefix
		// set it explicitly
		if !strings.HasPrefix(setToVersion, ver.Prefix) {
			setToVersion = ver.Prefix + setToVersion
		}
		v, err := ver.GetVersionFromTag(setToVersion)
		if err != nil {
			return errors.New("Couldn't get version from tag. " + err.Error())
		}

		user, err := ver.GetGitUser()
		if err != nil {
			return err
		}

		commit, err := ver.GetHeadCommit(repo)
		if err != nil {
			return err
		}

		id, err := repo.Tags.Create(v.String(), commit, user, v.String())
		if err != nil {
			return errors.New("Unable to create tag. " + err.Error())
		}

		fmt.Printf("Tag `%s` created successfully\n%s\n", v, id)

		if pushTags, _ := cmd.Flags().GetBool("push"); pushTags {
			// this is fucking embarrassing
			// please look away
			push := exec.Command("git", "push", "--tags")
			err = push.Run()
			if err != nil {
				return err
			}
		}

		return nil
	}

	versions := ver.Versions{}
	for _, tag := range tags {
		v, err := ver.GetVersionFromTag(tag)
		if err != nil {
			continue
		}

		versions = append(versions, *v)
	}

	fmt.Printf("%s\n", versions.Latest())

	return nil

}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of ver",
	Run:   versionCmdFn,
}

func versionCmdFn(cmd *cobra.Command, args []string) {
	fmt.Printf("ver version %s %s %s\n", BUILD_VERSION, BUILD_HASH, BUILD_DATE)
}

var incrementCmd = &cobra.Command{
	Use:   "i",
	Short: "Used to increment version",
	RunE:  incrementCmdFn,
}

func incrementCmdFn(cmd *cobra.Command, args []string) error {
	ver.Prefix, _ = cmd.Flags().GetString("prefix")

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

	versions := ver.Versions{}
	for _, tag := range tags {
		v, err := ver.GetVersionFromTag(tag)
		if err != nil {
			continue
			//return errors.New("Couldn't get version from tag. " + err.Error())
		}

		versions = append(versions, *v)
	}

	latestVer := versions.Latest()
	newVer := latestVer

	major, _ := cmd.Flags().GetBool("major")
	minor, _ := cmd.Flags().GetBool("minor")
	patch, _ := cmd.Flags().GetBool("patch")

	if major {
		newVer.Major += 1
		newVer.Minor = 0
		newVer.Patch = 0
	}
	if minor {
		newVer.Minor += 1
		newVer.Patch = 0
	}
	if patch {
		newVer.Patch += 1
	}

	setToVersion, _ := cmd.Flags().GetString("set")
	if setToVersion != "" {
		// has no version prefix
		// set it explicitly
		if !strings.HasPrefix(setToVersion, ver.Prefix) {
			setToVersion = ver.Prefix + setToVersion
		}
		v, err := ver.GetVersionFromTag(setToVersion)
		if err != nil {
			return errors.New("Couldn't get version from tag. " + err.Error())
		}

		newVer = *v
	}

	user, err := ver.GetGitUser()
	if err != nil {
		return err
	}

	commit, err := ver.GetHeadCommit(repo)
	if err != nil {
		return err
	}

	id, err := repo.Tags.Create(newVer.String(), commit, user, newVer.String())
	if err != nil {
		return errors.New("Unable to create tag. " + err.Error())
	}

	fmt.Printf("Tag `%s` created successfully\n%s\n", newVer, id)

	if pushTags, _ := cmd.Flags().GetBool("push"); pushTags {
		// this is fucking embarrassing
		// please look away
		push := exec.Command("git", "push", "--tags")
		err = push.Run()
		if err != nil {
			return err
		}
	}

	return nil

}

func init() {
	RootCmd.PersistentFlags().String("prefix", "v", "Prefix for git tag")
	RootCmd.PersistentFlags().StringP("set", "s", "", "Set version to this. e.g. ver -s \"v15.8.14\"")
	RootCmd.PersistentFlags().Bool("push", true, "Set to disable pushing tag to origin")

	incrementCmd.Flags().BoolP("major", "M", false, "Increase major version number")
	incrementCmd.Flags().BoolP("minor", "m", false, "Increase minor version number")
	incrementCmd.Flags().BoolP("patch", "p", false, "Increase patch version number")

	RootCmd.AddCommand(
		versionCmd,
		incrementCmd,
	)

}
func main() {
	if err := RootCmd.Execute(); err != nil {
		// dont print error
		// it's printed anyways
		return
	}
}
