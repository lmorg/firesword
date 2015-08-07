package main

func Trim(s string, length int) string {
	switch {
	case length > 0:
		return lTrim(s, length)

	case length < 0:
		return rTrim(s, -length)
	}

	return s
}

func lTrim(s string, length int) string {
	if len(s) <= length {
		return s
	}

	return "…" + s[len(s)-length+1:]
}

func rTrim(s string, length int) string {
	if len(s) <= length {
		return s
	}

	return s[:length-1] + "…"
}
