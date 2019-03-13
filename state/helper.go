package state

import "strings"

// StringsIdToStr helper for passing []string key
func StringsIdToStr(idSlice []string) string {
	return strings.Join(idSlice, "\000")
}

// StringsIdFromStr helper for restoring []string key
func StringsIdFromStr(idString string) []string {
	return strings.Split(idString, "\000")
}
