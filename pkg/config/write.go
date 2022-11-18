package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/svenliebig/work-environment/pkg/core"
)

func Write(path string, projects []*core.Project) error {
	result, err := json.MarshalIndent(&workEnvironmentConfig{Projects: projects}, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	d, _ := filepath.Split(path)

	if err := os.Mkdir(d, os.ModePerm); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	err = os.WriteFile(path, result, 0644)

	if err != nil {
		return err
	}

	return nil
}
