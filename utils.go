package chromedriver2har

func safeStringDereference(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}
