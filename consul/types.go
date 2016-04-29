package consul

import "time"

type Branch struct {
	Project   string    `json:"-"`
	Name      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Branches []*Branch

func (branches *Branches) Contains(branch string) bool {
	for _, b := range *branches {
		if b.Name == branch {
			return true
		}
	}
	return false
}

type Distribution map[string]string

func (d *Distribution) Contains(branch string) bool {
	for _, v := range *d {
		if v == branch {
			return true
		}
	}

	return false
}

type Projects []string

func (p *Projects) Contains(project string) bool {
	for _, v := range *p {
		if v == project {
			return true
		}
	}

	return false
}
