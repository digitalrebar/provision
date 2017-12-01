package models

import (
	"sort"
	"strings"
)

// Meta holds information about arbritary things.
// Initial usage will be for UX elements.
//
// swagger: model
type MetaData struct {
	Meta map[string]string
}

func (m *MetaData) fill() {
	if m.Meta == nil {
		m.Meta = map[string]string{}
	}
}

func (m *MetaData) ClearFeatures() {
	m.Meta["feature-flags"] = ""
}

func (m *MetaData) Features() []string {
	m.fill()
	if flags, ok := m.Meta["feature-flags"]; ok && len(flags) > 0 {
		return strings.Split(flags, ",")
	}
	return []string{}
}

func (m *MetaData) HasFeature(flag string) bool {
	flag = strings.TrimSpace(flag)
	for _, testFlag := range m.Features() {
		if flag == strings.TrimSpace(testFlag) {
			return true
		}
	}
	return false
}

func (m *MetaData) AddFeature(flag string) {
	flag = strings.TrimSpace(flag)
	if m.HasFeature(flag) || flag == "" {
		return
	}
	flags := m.Features()
	flags = append(flags, flag)
	sort.Strings(flags)
	m.Meta["feature-flags"] = strings.Join(flags, ",")
}

func (m *MetaData) RemoveFeature(flag string) {
	flag = strings.TrimSpace(flag)
	newFlags := []string{}
	for _, testFlag := range m.Features() {
		if flag == testFlag {
			continue
		}
		newFlags = append(newFlags, testFlag)
	}
	m.Meta["feature-flags"] = strings.Join(newFlags, ",")
}

func (m *MetaData) MergeFeatures(other []string) {
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
	m.Meta["feature-flags"] = strings.Join(newFlags, ",")
}
