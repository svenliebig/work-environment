package gitazuredevops

import (
	"fmt"
	"strings"
)

// parses the repository ID from the remote URL
// git@ssh.dev.azure.com:v3/{org}/{project}}/{repository}
func (c *client) getRepositoryId() (string, error) {
	remoteUrl := c.ctx.Project().Git.RemoteUrl

	if remoteUrl == "" {
		return "", fmt.Errorf("remote URL is empty")
	}

	if !strings.Contains(remoteUrl, "ssh.dev.azure.com") {
		return "", fmt.Errorf("remote URL is not an Azure DevOps repository")
	}

	splits := strings.Split(remoteUrl, "/")

	if len(splits) < 4 {
		return "", fmt.Errorf("remote URL is invalid")
	}

	return splits[3], nil
}

func (c *client) getProject() (string, error) {
	remoteUrl := c.ctx.Project().Git.RemoteUrl

	if remoteUrl == "" {
		return "", fmt.Errorf("remote URL is empty")
	}

	if !strings.Contains(remoteUrl, "ssh.dev.azure.com") {
		return "", fmt.Errorf("remote URL is not an Azure DevOps repository")
	}

	splits := strings.Split(remoteUrl, "/")

	if len(splits) < 4 {
		return "", fmt.Errorf("remote URL is invalid")
	}

	return splits[2], nil
}
