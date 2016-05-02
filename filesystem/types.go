package filesystem

type Branches []string

func (branches *Branches) Contains(branch string) bool {
	for _, b := range *branches {
		if b == branch {
			return true
		}
	}
	return false
}
