package util

import (
	"os"
	"os/user"
	path2 "path"
)

func GetHomePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	path := path2.Join(usr.HomeDir, ".fb_crm")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return "", nil
		}
	}

	return path, nil
}
