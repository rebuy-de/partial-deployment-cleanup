package consul

import "time"

type Deployments []*Deployment

type Deployment struct {
	Project   string    `json:"-"`
	Branch    string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
