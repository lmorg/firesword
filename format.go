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

/*
    Lower(int, int)int{} is currently unused.
    I was debating having the following feature:
    ` --lower str    Round field down to int. eg "size:1024;" for size in KB.
                  (semicolon delimited, -hf for field names)`
    However it's a surprisingly large amount of work for a feature that would
    only really be useful for piping, eg | uniq -c
*/
//func Lower(val, rounding int) int {
//	return int(val/rounding) * rounding
//}
