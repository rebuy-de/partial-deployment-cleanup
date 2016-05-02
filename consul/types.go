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

func (branches *Branches) Slice() []string {
	slice := make([]string, 0)

	for _, b := range *branches {
		slice = append(slice, b.Name)
	}

	return slice
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

func (d *Distribution) BranchSlice() []string {
	pseudoMap := make(map[string]bool)
	for _, v := range *d {
		pseudoMap[v] = true
	}

	slice := make([]string, 0)
	for b, _ := range pseudoMap {
		slice = append(slice, b)
	}

	return slice
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
