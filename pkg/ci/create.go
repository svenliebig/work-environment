package ci

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/svenliebig/work-environment/pkg/utils/bamboo"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/wepath"
)

func Create(p string, url string, ciType string, name string, auth string) error {
	wer, err := wepath.GetWorkEnvironmentRoot(p)

	if err != nil {
		return err
	}

	// TODO validate parameters

	if ciType != "bamboo" {
		fmt.Printf("the type %q is not a valid ci type\n", ciType)
		return nil
	}

	client := &bamboo.Client{
		BaseUrl:   url,
		AuthToken: auth,
	}

	version, err := client.GetInfo(context.Background())

	if err != nil {
		return err
	}

	cp := filepath.Join(wer, ConfigFilename)

	if wepath.Exists(cp) {
		c, err := ReadConfig(cp)

		if err != nil {
			return err
		}

		if c.Contains(name) {
			answer := cli.Question(fmt.Sprintf("\nthe identifier %q is already declared in the config of your work environment, do you want to override it? (y/n) ", name), []string{"y", "n"})
			if answer == "n" {
				return nil
			}
		}
	}

	ci := &CI{
		CiType:     ciType,
		Identifier: name,
		AuthToken:  auth,
		Url:        url,
	}

	result, err := json.MarshalIndent(&CIConfig{Environments: []*CI{ci}}, "", "  ")

	if err != nil {
		return err
	}

	if err := os.WriteFile(cp, result, 0644); err != nil {
		return err
	}

	// TODO make this pretty
	fmt.Printf("\nsuccessfully added a new CI to the work environment:\n\tidentifier: %s\n\ttype: %s\n\tversion: %s\n\turl: %s\n\n", name, ciType, version.Version, url)

	return nil
}
