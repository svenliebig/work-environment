package git

func RepositoryGetRemoteOriginUrl(p string) (string, error) {
	r, err := getRepository(p)

	if err != nil {
		return "", err
	}

	c, err := r.Config()

	if err != nil {
		return "", err
	}

	if origin, ok := c.Remotes["origin"]; ok {
		return origin.URLs[0], nil
	}

	return "", nil
}
