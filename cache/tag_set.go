package cache

import (
	"strings"
)

type TagSet struct {
	names []string
	store Store
}

func NewTagSet(store Store, names []string) *TagSet {
	return &TagSet{names: names, store: store}
}

func (t *TagSet) Reset() error {
	for _, name := range t.names {
		_, err := t.ResetTag(name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TagSet) ResetTag(name string) (string, error) {
	id := Sha1(name)

	err := t.store.Forever(t.TagKey(name), id)

	return id, err
}

func (t *TagSet) GetNamespace() (string, error) {
	var err error
	var names = make([]string, len(t.names))
	for i, name := range t.names {
		name, err = t.TagId(name)
		if err != nil {
			return "", err
		}
		names[i] = name
	}
	return strings.Join(names, "|"), nil
}

func (t *TagSet) TagId(name string) (string, error) {
	var id string
	tagKey := t.TagKey(name)
	err := t.store.Get(tagKey, &id)
	if err != nil {
		return t.ResetTag(name)
	}

	return id, nil
}

func (t *TagSet) TagKey(name string) string {
	return "tag:" + name + ":key"
}
