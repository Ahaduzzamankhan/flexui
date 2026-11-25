package extractor

// scanScript extracts candidates from string and template literals in a script
// source. Comments are skipped and template literal substitutions are scanned
// recursively.
func scanScript(src string, add func(string)) {
	i := 0
	n := len(src)
	for i < n {
		switch c := src[i]; {
		case c == '/' && i+1 < n && src[i+1] == '/':
			for i < n && src[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < n && src[i+1] == '*':
			i += 2
			for i+1 < n && !(src[i] == '*' && src[i+1] == '/') {
				i++
			}
			i += 2
		case c == '"' || c == '\'':
			i++
			start := i
			for i < n && src[i] != c {
				if src[i] == '\\' && i+1 < n {
					i++
				}
				i++
			}
			Tokens(src[start:minInt(i, n)], add)
			if i < n {
				i++
			}
		case c == '`':
			i++
			start := i
			for i < n && src[i] != '`' {
				if src[i] == '\\' && i+1 < n {
					i += 2
					continue
				}
				if src[i] == '$' && i+1 < n && src[i+1] == '{' {
					Tokens(src[start:i], add)
					end := balanced(src, i+1)
					if end < 0 {
						return
					}
					scanScript(src[i+2:end], add)
					i = end + 1
					start = i
					continue
				}
				i++
			}
			Tokens(src[start:minInt(i, n)], add)
			if i < n {
				i++
			}
		default:
			i++
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
