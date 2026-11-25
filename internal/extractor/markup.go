package extractor

import "strings"

// classAttributes are attribute names whose values are class lists. Names are
// compared lower-cased.
var classAttributes = map[string]bool{
	"class":        true,
	"classname":    true,
	"class:list":   true,
	":class":       true,
	"v-bind:class": true,
	"ng-class":     true,
	"ngclass":      true,
	"[ngclass]":    true,
	"[class]":      true,
	"tw":           true,
}

// expressionAttributes are class-carrying attributes whose value is a script
// expression rather than a plain class list.
var expressionAttributes = map[string]bool{
	":class":       true,
	"v-bind:class": true,
	"class:list":   true,
	"ng-class":     true,
	"ngclass":      true,
	"[ngclass]":    true,
	"[class]":      true,
}

// scanMarkup walks an HTML/JSX/Vue/Svelte/Astro file.
func scanMarkup(src string, add func(string)) {
	i := 0
	n := len(src)
	for i < n {
		if src[i] != '<' {
			next := strings.IndexByte(src[i:], '<')
			if next < 0 {
				return
			}
			i += next
			continue
		}
		switch {
		case strings.HasPrefix(src[i:], "<!--"):
			end := strings.Index(src[i:], "-->")
			if end < 0 {
				return
			}
			i += end + 3
			continue
		case strings.HasPrefix(src[i:], "<!"):
			end := strings.IndexByte(src[i:], '>')
			if end < 0 {
				return
			}
			i += end + 1
			continue
		case strings.HasPrefix(src[i:], "</"):
			end := strings.IndexByte(src[i:], '>')
			if end < 0 {
				return
			}
			i += end + 1
			continue
		}
		i++
		nameStart := i
		for i < n && isTagNameByte(src[i]) {
			i++
		}
		name := strings.ToLower(src[nameStart:i])
		i, selfClosing := scanAttributes(src, i, add)
		if i >= n {
			return
		}
		if selfClosing {
			continue
		}
		switch name {
		case "script":
			end := indexFold(src[i:], "</script")
			if end < 0 {
				return
			}
			scanScript(src[i:i+end], add)
			i = skipPast(src, i+end, ">")
		case "style":
			end := indexFold(src[i:], "</style")
			if end < 0 {
				return
			}
			i = skipPast(src, i+end, ">")
		}
	}
}

// scanAttributes reads attributes starting at i (just after the tag name) and
// returns the index after the closing '>' plus whether the tag self-closed.
func scanAttributes(src string, i int, add func(string)) (int, bool) {
	n := len(src)
	for i < n {
		for i < n && isSpace(src[i]) {
			i++
		}
		if i >= n {
			return n, false
		}
		if src[i] == '>' {
			return i + 1, false
		}
		if src[i] == '/' && i+1 < n && src[i+1] == '>' {
			return i + 2, true
		}
		if src[i] == '{' {
			end := balanced(src, i)
			if end < 0 {
				return n, false
			}
			i = end + 1
			continue
		}
		start := i
		for i < n && !isSpace(src[i]) && src[i] != '=' && src[i] != '>' && src[i] != '/' && src[i] != '{' {
			i++
		}
		if i == start {
			i++
			continue
		}
		attr := strings.ToLower(src[start:i])
		for i < n && isSpace(src[i]) {
			i++
		}
		if i >= n {
			return n, false
		}
		if src[i] != '=' {
			if strings.HasPrefix(attr, "class:") {
				if token := attr[len("class:"):]; token != "" && ShapeOK(token) {
					add(token)
				}
			}
			continue
		}
		i++
		for i < n && isSpace(src[i]) {
			i++
		}
		if i >= n {
			return n, false
		}
		switch src[i] {
		case '"', '\'':
			quote := src[i]
			i++
			valueStart := i
			for i < n && src[i] != quote {
				if src[i] == '\\' && i+1 < n {
					i++
				}
				i++
			}
			value := src[valueStart:i]
			if i < n {
				i++
			}
			if classAttributes[attr] {
				if expressionAttributes[attr] {
					scanScript(value, add)
				} else {
					Tokens(value, add)
				}
			}
		case '{':
			end := balanced(src, i)
			if end < 0 {
				return n, false
			}
			if classAttributes[attr] {
				scanScript(src[i+1:end], add)
			}
			i = end + 1
		default:
			valueStart := i
			for i < n && !isSpace(src[i]) && src[i] != '>' {
				i++
			}
			if classAttributes[attr] {
				Tokens(src[valueStart:i], add)
			}
		}
	}
	return n, false
}

// balanced returns the index of the bracket closing the one at open, honouring
// strings, template literals and comments.
func balanced(src string, open int) int {
	if open >= len(src) {
		return -1
	}
	var closer byte
	switch src[open] {
	case '{':
		closer = '}'
	case '[':
		closer = ']'
	case '(':
		closer = ')'
	default:
		return -1
	}
	depth := 0
	for i := open; i < len(src); i++ {
		switch c := src[i]; c {
		case '\'', '"', '`':
			i = skipString(src, i)
		case '/':
			if i+1 < len(src) && src[i+1] == '/' {
				for i < len(src) && src[i] != '\n' {
					i++
				}
			} else if i+1 < len(src) && src[i+1] == '*' {
				i += 2
				for i+1 < len(src) && !(src[i] == '*' && src[i+1] == '/') {
					i++
				}
				i++
			}
		case '{':
			depth++
		case '}':
			depth--
			if depth < 0 {
				return -1
			}
			if depth == 0 && closer == '}' {
				return i
			}
		case '[':
			if closer == ']' {
				depth++
			}
		case ']':
			if closer == ']' {
				depth--
				if depth == 0 {
					return i
				}
			}
		case '(':
			if closer == ')' {
				depth++
			}
		case ')':
			if closer == ')' {
				depth--
				if depth == 0 {
					return i
				}
			}
		}
	}
	return -1
}

// skipString returns the index of the closing quote of the string starting at i.
func skipString(src string, i int) int {
	quote := src[i]
	if quote == '`' {
		for i++; i < len(src); i++ {
			if src[i] == '\\' {
				i++
				continue
			}
			if src[i] == '`' {
				return i
			}
			if src[i] == '$' && i+1 < len(src) && src[i+1] == '{' {
				end := balanced(src, i+1)
				if end < 0 {
					return len(src) - 1
				}
				i = end
			}
		}
		return len(src) - 1
	}
	for i++; i < len(src); i++ {
		if src[i] == '\\' {
			i++
			continue
		}
		if src[i] == quote {
			return i
		}
	}
	return len(src) - 1
}

func indexFold(src, lower string) int {
	return strings.Index(strings.ToLower(src), lower)
}

func skipPast(src string, from int, want string) int {
	if from >= len(src) {
		return len(src)
	}
	idx := strings.Index(src[from:], want)
	if idx < 0 {
		return len(src)
	}
	return from + idx + len(want)
}

func isTagNameByte(c byte) bool {
	switch {
	case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		return true
	case c == '-' || c == ':' || c == '_' || c == '.':
		return true
	}
	return false
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}
