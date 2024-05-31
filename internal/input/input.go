package input

import "maps"

type StateInput map[string]any

func (i StateInput) GetStringSlice(key string) []string {
	v, ok := i[key]
	if !ok {
		return []string{}
	}
	s, ok := v.([]string)
	if !ok {
		return []string{}
	}
	return s
}

func (i StateInput) GetString(key string) string {
	v, ok := i[key]
	if !ok {
		return ""
	}

	s, ok := v.(string)
	if !ok {
		return ""
	}

	return s
}

func (i StateInput) GetInt(key string) int {
	v, ok := i[key]
	if !ok {
		return 0
	}

	n, ok := v.(int)
	if !ok {
		return 0
	}

	return n
}

func (i StateInput) GetFloat64(key string) float64 {
	v, ok := i[key]
	if !ok {
		return 0
	}

	n, ok := v.(float64)
	if !ok {
		return 0
	}

	return n
}

func (i StateInput) GetBool(key string) bool {
	v, ok := i[key]
	if !ok {
		return false
	}

	n, ok := v.(bool)
	if !ok {
		return false
	}

	return n
}

func (i StateInput) Get(key string) any {
	return i[key]
}

func (i StateInput) Merged(other StateInput) StateInput {
	clone := maps.Clone(i)
	for key, val := range other {
		if _, ok := clone[key]; !ok {
			clone[key] = val
		}
	}
	return clone
}
