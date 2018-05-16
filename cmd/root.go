package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/47bytes/ver/model"
	"github.com/spf13/cobra"
	git "gopkg.in/libgit2/git2go.v25"
)

var RootCmd = &cobra.Command{
	Use:     "ver",
	Short:   "ver - a simple git tag semver version incrementer",
	Long:    "ver increments semver style git tags",
	Example: "$ ver -m -p\n Tag `v0.2.1` created successfully\n 27c1f1234188aa11585334726f8721d9a35038eb",
	RunE:    rootCmdFn,
}

func rootCmdFn(cmd *cobra.Command, args []string) error {

	model.PREFIX, _ = cmd.Flags().GetString("prefix")

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

	versions := model.Versions{}
	for _, tag := range tags {
		v, err := model.GetVersionFromTag(tag)
		if err != nil {
			return errors.New("Couldn't get version from tag. " + err.Error())
		}

		versions = append(versions, *v)
	}

	latestVer := versions.Latest()
	newVer := latestVer

	printLatest, _ := cmd.Flags().GetBool("latest")
	if printLatest || (latestVer == newVer) {
		fmt.Printf("%s\n", latestVer)
		return nil
	}

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
		v, err := model.GetVersionFromTag(setToVersion)
		if err != nil {
			return errors.New("Couldn't get version from tag. " + err.Error())
		}

		newVer = *v
	}

	// if latestVer == newVer {
	// 	return errors.New("No flags provided but at least (1) flag has to be set!")
	// }

	user, err := model.GetGitUser()
	if err != nil {
		return err
	}

	commit, err := model.GetHeadCommit(repo)
	if err != nil {
		return err
	}

	id, err := repo.Tags.Create(newVer.String(), commit, user, newVer.String())
	if err != nil {
		return errors.New("Unable to create tag. " + err.Error())
	}

	fmt.Printf("Tag `%s` created successfully\n%s\n", newVer, id)

	return nil

}

func init() {
	RootCmd.Flags().String("prefix", "v", "Prefix for git tag")
	RootCmd.Flags().BoolP("major", "M", false, "Increase major version number")
	RootCmd.Flags().BoolP("minor", "m", false, "Increase minor version number")
	RootCmd.Flags().BoolP("patch", "p", false, "Increase patch version number")
	RootCmd.Flags().StringP("set", "s", "", "Set version to this. e.g. ver -s \"v15.8.14\"")
	RootCmd.Flags().BoolP("latest", "l", false, "Print latest version")
}
