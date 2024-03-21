package core

import (
	"errors"
)

var (
	ErrCIAlreadyExists  = errors.New("a ci with that identifier does already exists")
	ErrVCSAlreadyExists = errors.New("a vcs with that identifier does already exists")
	ErrNoSuchCI         = errors.New("there is no ci with that identifier available")
)

const ConfigurationFileName = ".work-environment.json"

type Configuration struct {
	Projects        []*Project
	CIEnvironments  []*CI
	VCSEnvironments []*VCS

	// tells if the configuration has been modified since
	// the last reading or writing from the configuration file
	dirty bool
}

// returns the ci environment with the given identifier, if it exists
func (c *Configuration) GetCIEnvironmentById(identifier string) (*CI, error) {
	if len(c.CIEnvironments) == 0 {
		return nil, errors.New("no ci environments defined")
	} else {
		if identifier == "" {
			return nil, errors.New("parameter identifier is empty, but required when there are more then one ci environment")
		}

		for _, e := range c.CIEnvironments {
			if e.Identifier == identifier {
				return e, nil
			}
		}
		return nil, errors.New("no such ci environment with the given identifier")
	}
}

func (c *Configuration) GetVCSEnvironmentById(identifier string) (*VCS, error) {
	if len(c.VCSEnvironments) == 0 {
		return nil, errors.New("no vcs environments defined")
	} else {
		if identifier == "" {
			return nil, errors.New("parameter identifier is empty, but required when there are more then one vcs environment")
		}

		for _, e := range c.VCSEnvironments {
			if e.Identifier == identifier {
				return e, nil
			}
		}

		return nil, errors.New("no such vcs environment with the given identifier")
	}
}

func (c *Configuration) HasCI(id string) bool {
	for _, ci := range c.CIEnvironments {
		if ci.Identifier == id {
			return true
		}
	}
	return false
}

// Put this into context?
func (c *Configuration) UpdateProjectCI(identifier string, pci *ProjectCI) error {
	c.dirty = true

	for _, p := range c.Projects {
		if p.Identifier == identifier {
			p.CI = pci
			return nil
		}
	}

	return errors.New("project not found")
}

func (c *Configuration) UpdateProjectCD(identifier string, pcd *ProjectCD) error {
	c.dirty = true

	for _, p := range c.Projects {
		if p.Identifier == identifier {
			p.CD = pcd
			return nil
		}
	}

	return errors.New("project not found")
}

func (c *Configuration) UpdateProjectVCS(identifier string, pvcs *ProjectVCS) error {
	c.dirty = true

	for _, p := range c.Projects {
		if p.Identifier == identifier {
			p.VCS = pvcs
			return nil
		}
	}

	return errors.New("project not found")
}

// checks the configuration for a project identifier
// returns true and a pointer to the project, in case
// a project got found
func (c *Configuration) ContainsProject(identifier string) (bool, *Project) {
	for _, project := range c.Projects {
		if project.Identifier == identifier {
			return true, project
		}
	}
	return false, nil
}

// adds a project to the configuration and sets it's dirty flag
func (c *Configuration) AddProject(p *Project) {
	c.dirty = true
	c.Projects = append(c.Projects, p)
}

func (c *Configuration) RemoveProject(p *Project) error {
	c.dirty = true
	for i, pr := range c.Projects {
		if pr.Identifier == p.Identifier {
			t := make([]*Project, 0, len(c.Projects)-1)
			p := c.Projects[:i]
			t = append(t, p...)
			e := c.Projects[i+1:]
			t = append(t, e...)
			c.Projects = t
			return nil
		}
	}
	return errors.New("project not found")
}

func (c *Configuration) AddVCS(vcs *VCS) error {
	c.dirty = true

	if c.VCSEnvironments == nil {
		c.VCSEnvironments = []*VCS{vcs}
	} else {
		for _, e := range c.VCSEnvironments {
			if e.Identifier == vcs.Identifier {
				return ErrVCSAlreadyExists
			}
		}

		c.VCSEnvironments = append(c.VCSEnvironments, vcs)
	}

	return nil
}

func (c *Configuration) AddCI(ci *CI) error {
	c.dirty = true

	if c.CIEnvironments == nil {
		c.CIEnvironments = []*CI{ci}
	} else {
		for _, cie := range c.CIEnvironments {
			if cie.Identifier == ci.Identifier {
				return ErrCIAlreadyExists
			}
		}

		c.CIEnvironments = append(c.CIEnvironments, ci)
	}

	return nil
}

func (c *Configuration) UpdateCI(ci *CI) error {
	c.dirty = true

	for i, cie := range c.CIEnvironments {
		if cie.Identifier == ci.Identifier {
			c.CIEnvironments[i] = ci
			return nil
		}
	}

	return ErrNoSuchCI
}

func (c *Configuration) IsDirty() bool {
	return c.dirty
}
