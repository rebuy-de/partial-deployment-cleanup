package consul

import "strings"

type Key string

func (k Key) Add(sub string) Key {
	return Key(k.Clean() + "/" + Key(sub).Clean())
}

func (k Key) Base(root Key) Key {
	key := strings.TrimPrefix(k.Clean(), root.Clean())
	key = strings.Trim(key, "/")
	return Key(key)
}

func (k Key) Get(i int) string {
	parts := strings.Split(k.Clean(), "/")
	if len(parts) <= i {
		return ""
	}
	return parts[i]
}

func (k Key) Clean() string {
	return strings.Trim(string(k), "/")
}
