package git

func BranchName(path string) (string, error) {
	r, err := getRepository(path)

	if err != nil {
		return "", err
	}

	h, err := r.Head()

	if err != nil {
		return "", err
	}

	if h.Name().IsBranch() {
		return h.Name().Short(), nil
	}

	return "", nil
}
