package config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	componentNamePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	classNamePattern     = regexp.MustCompile(`^([a-z][a-z0-9-]*|[0-9]+xl)$`)
)

// validateComponents rejects component definitions that could not produce
// valid, predictable class names. Catching this at config load time means a
// typo is a build error with a file name, not a silently missing class.
func validateComponents(path string, components map[string]ComponentConfig) error {
	names := make([]string, 0, len(components))
	for name := range components {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		component := components[name]
		if !componentNamePattern.MatchString(name) {
			return fmt.Errorf("parse config %s: component name %q must be lowercase letters, digits and single dashes", path, name)
		}
		if len(component.Base) == 0 && len(component.Variants) == 0 &&
			len(component.Parts) == 0 && len(component.States) == 0 &&
			component.Extends == "" {
			return fmt.Errorf("parse config %s: component %q defines no base, variants, parts, states or extends", path, name)
		}
		if component.Extends != "" && !componentNamePattern.MatchString(component.Extends) &&
			!strings.HasPrefix(component.Extends, "fui-") {
			return fmt.Errorf("parse config %s: component %q extends %q, which is not a component name", path, name, component.Extends)
		}
		if err := validateDecls(path, "components."+name+".base", component.Base); err != nil {
			return err
		}
		for _, key := range sortedMapKeys(component.Variants) {
			if !classNamePattern.MatchString(key) {
				return fmt.Errorf("parse config %s: component %q variant %q must be lowercase letters, digits and dashes", path, name, key)
			}
			if err := validateDecls(path, "components."+name+".variants."+key, component.Variants[key]); err != nil {
				return err
			}
		}
		for _, key := range sortedMapKeys(component.Parts) {
			if !classNamePattern.MatchString(key) {
				return fmt.Errorf("parse config %s: component %q part %q must be lowercase letters, digits and dashes", path, name, key)
			}
			if err := validateDecls(path, "components."+name+".parts."+key, component.Parts[key]); err != nil {
				return err
			}
		}
		for _, key := range sortedMapKeys(component.States) {
			if !validStateSelector(key) {
				return fmt.Errorf("parse config %s: component %q state %q must start with one of : [ . space > + ~", path, name, key)
			}
			if err := validateDecls(path, "components."+name+".states."+key, component.States[key]); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateDecls(path, where string, decls map[string]string) error {
	props := make([]string, 0, len(decls))
	for prop := range decls {
		props = append(props, prop)
	}
	sort.Strings(props)
	for _, prop := range props {
		value := decls[prop]
		if strings.TrimSpace(prop) == "" {
			return fmt.Errorf("parse config %s: %s has an empty property name", path, where)
		}
		if strings.ContainsAny(prop, "{};:") {
			return fmt.Errorf("parse config %s: %s property %q contains an invalid character", path, where, prop)
		}
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("parse config %s: %s property %q has an empty value", path, where, prop)
		}
		if strings.ContainsAny(value, "{};") {
			return fmt.Errorf("parse config %s: %s property %q value %q must not contain { } or ;", path, where, prop, value)
		}
	}
	return nil
}

func validStateSelector(selector string) bool {
	trimmed := strings.TrimSpace(selector)
	if trimmed == "" {
		return false
	}
	switch trimmed[0] {
	case ':', '[', '.', '>', '+', '~':
		return true
	}
	return false
}

func sortedMapKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
