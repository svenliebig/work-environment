package utils

import (
	"log"
	"os"
	"path"
)

func GetPath(s []string) (string, error) {
	if len(s) == 1 {
		p := s[0]
		if path.IsAbs(p) {
			return p, nil
		} else {
			wd, err := os.Getwd()

			if err != nil {
				return "", err
			}

			return path.Join(wd, p), nil
		}
	}

	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	return cwd, nil
}
