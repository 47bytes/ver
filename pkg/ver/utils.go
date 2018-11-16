package ver

import (
	"errors"
	"time"

	git "gopkg.in/libgit2/git2go.v25"
)

func GetGitUser() (*git.Signature, error) {
	confPath, err := git.ConfigFindGlobal()
	if err != nil {
		return nil, errors.New("Couldn't locate git config. " + err.Error())
	}

	simpleConf, err := git.NewConfig()
	if err != nil {
		return nil, err
	}

	conf, err := git.OpenOndisk(simpleConf, confPath)
	if err != nil {
		return nil, err
	}

	name, err := conf.LookupString("user.name")
	if err != nil {
		return nil, errors.New("Couldn't find user.name git config key. " + err.Error())
	}
	email, err := conf.LookupString("user.email")
	if err != nil {
		return nil, errors.New("Couldn't find user.email git config key. " + err.Error())
	}

	user := &git.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}

	return user, nil
}

func GetHeadCommit(repo *git.Repository) (*git.Commit, error) {
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	oid := head.Target()

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
