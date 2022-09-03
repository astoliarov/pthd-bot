package services

import "strings"

func normalizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, "ё", "e", -1)
	return name
}
