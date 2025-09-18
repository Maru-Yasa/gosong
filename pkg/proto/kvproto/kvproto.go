package kvproto

import (
	"fmt"
	"strings"
)

// encode map[string]string -> string key=value key2=value2
func Encode(data map[string]string) string {
	parts := make([]string, 0, len(data))
	for k, v := range data {
		if strings.Contains(v, " ") {
			v = fmt.Sprintf("\"%s\"", v)
		}
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, " ")
}

// decode string "key=value key2=value2" -> map[string]string
func Decode(s string) map[string]string {
	result := make(map[string]string)
	parts := strings.Fields(s)

	for _, part := range parts {
		if !strings.Contains(part, "=") {
			continue
		}
		pair := strings.SplitN(part, "=", 2)
		key := pair[0]
		val := strings.Trim(pair[1], `"`)
		result[key] = val
	}
	return result
}
