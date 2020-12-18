package common

// StringValueOrDefault returns original value if not empty. Default otherwise.
func StringValueOrDefault(str string, def string) string {
	if str != "" {
		return str
	}
	return def
}

// StringFromRef returns original value if not empty. Default otherwise.
func StringFromRef(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
