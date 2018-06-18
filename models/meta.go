package models

import (
	"sort"
	"strings"
)

// Meta holds information about arbitrary things.
// Initial usage will be for UX elements.
//
// swagger: model
type Meta map[string]string

type MetaHaver interface {
	Model
	GetMeta() Meta
	SetMeta(Meta)
}

func (m Meta) ClearFeatures() {
	m["feature-flags"] = ""
}

func (m Meta) Features() []string {
	if flags, ok := m["feature-flags"]; ok && len(flags) > 0 {
		return strings.Split(flags, ",")
	}
	return []string{}
}

func (m Meta) HasFeature(flag string) bool {
	flag = strings.TrimSpace(flag)
	for _, testFlag := range m.Features() {
		if flag == strings.TrimSpace(testFlag) {
			return true
		}
	}
	return false
}

func (m Meta) AddFeature(flag string) {
	flag = strings.TrimSpace(flag)
	if m.HasFeature(flag) || flag == "" {
		return
	}
	flags := m.Features()
	flags = append(flags, flag)
	sort.Strings(flags)
	m["feature-flags"] = strings.Join(flags, ",")
}

func (m Meta) RemoveFeature(flag string) {
	flag = strings.TrimSpace(flag)
	newFlags := []string{}
	for _, testFlag := range m.Features() {
		if flag == testFlag {
			continue
		}
		newFlags = append(newFlags, testFlag)
	}
	m["feature-flags"] = strings.Join(newFlags, ",")
}

func (m Meta) MergeFeatures(other []string) {
	flags := map[string]struct{}{}
	for _, flag := range m.Features() {
		flags[flag] = struct{}{}
	}
	for _, flag := range other {
		flag = strings.TrimSpace(flag)
		if flag != "" {
			flags[flag] = struct{}{}
		}
	}
	newFlags := make([]string, 0, len(flags))
	for k := range flags {
		newFlags = append(newFlags, k)
	}
	sort.Strings(newFlags)
	m["feature-flags"] = strings.Join(newFlags, ",")
}
