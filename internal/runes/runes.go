package runes

type TestFunc func(r rune) bool

// HasPrefix tests whether 's' begins with 'prefix'.
func HasPrefix(s, prefix []rune) bool {
	return len(s) >= len(prefix) && Compare(s[:len(prefix)], prefix) == 0
}

// HasSuffix tests whether 's' ends with 'suffix'.
func HasSuffix(s, suffix []rune) bool {
	return len(s) >= len(suffix) && Compare(s[len(s)-len(suffix):], suffix) == 0
}

func SplitRune(s []rune, r rune) [][]rune {
	l := make([][]rune, 0, 5)

	prev_idx := 0
	idx := 0
	for ; idx < len(s); idx++ {
		sr := s[idx]
		if sr == r {
			l = append(l, s[prev_idx:idx])
			prev_idx = idx + 1
		}
	}
	// append remainder
	return append(l, s[prev_idx:])
}

// Index returns the index of the first occurrence of 'needle' in 's'.
func Index(s, needle []rune) int {
	switch {
	case len(needle) == 0:
		return 0
	case len(needle) == 1: // just a single rune, same as IndexRune()
		return IndexRune(s, needle[0])
	case len(needle) > len(s): // this can never succeed
		return -1
	case len(needle) == len(s): // same length, operation becomes a simple comparison
		if Compare(s, needle) == 0 {
			return 0
		}
		return -1
	}

head:
	for idx := 0; idx <= len(s)-len(needle); idx++ {
		// attempt to match 'needle' from here
		for nidx, nr := range needle {
			if s[idx+nidx] != nr {
				continue head
			}
		}
		return idx
	}
	return -1
}

// IndexAny returns the index of the first occurrence in 's' that is in 'any'.
func IndexAny(s, any []rune) int {
	if len(any) == 0 {
		return -1
	}
	for idx, c := range s {
		for _, m := range any {
			if c == m {
				return idx
			}
		}
	}
	return -1
}

// IndexNotAny returns the index of the first occurrence in 's' that is not in 'any'.
func IndexNotAny(s, any []rune) int {
	if len(any) == 0 {
		return -1
	}

next_rune:
	for idx, r := range s {
		// does 'any' contain 'r' ?
		for _, ar := range any {
			if ar == r {
				continue next_rune
			}
		}
		// no: 'any' does not contain 'r'
		return idx
	}
	return -1
}

// IndexRune returns the index to the first occurrence in 's' that equals 'r'.
func IndexRune(s []rune, r rune) int {
	for idx, c := range s {
		if c == r {
			return idx
		}
	}
	return -1
}

// IndexNotRune returns the index to the first occurrance in 's' that does not equal 'r'.
func IndexNotRune(s []rune, r rune) int {
	for idx, c := range s {
		if c != r {
			return idx
		}
	}
	return -1
}

// IndexFunc returns the index to the first occurrence in 's' where 'f' returns true.
func IndexFunc(s []rune, f TestFunc) int {
	for idx, c := range s {
		if f(c) {
			return idx
		}
	}
	return -1
}

// IndexNotFunc returns the index to the first occurrence in 's' where 'f' returns false.
func IndexNotFunc(s []rune, f TestFunc) int {
	for idx, c := range s {
		if !f(c) {
			return idx
		}
	}
	return -1
}

// LastIndex returns the index to the last occurrence of 'needle' in 's'.
func LastIndex(s, needle []rune) int {
	switch {
	case len(needle) == 0:
		return len(s)
	case len(needle) > len(s): // this will never succeed
		return -1
	case len(needle) == 1: // just a single rune, same as LastIndexRune()
		return LastIndexRune(s, needle[0])
	case len(needle) == len(s): // same length, operation becomes a simple comparison
		if Compare(s, needle) == 0 {
			return 0
		}
		return -1
	}

head:
	for idx := len(s) - len(needle); idx >= 0; idx-- {
		// attempt to match 'needle' from here
		for nidx := 0; nidx < len(needle); nidx++ {
			sr := s[idx+len(needle)-nidx-1]
			nr := needle[len(needle)-nidx-1]
			if sr != nr {
				continue head
			}
		}

		return idx // index to start of 'needle'
	}

	return -1
}

// Compare returns:
//   a == b -> 0
//    a > b -> 1
//    a < b -> -1
func Compare(a, b []rune) int {
	if len(a) != len(b) {
		if len(a) > len(b) {
			return 1
		}
		return -1
	}
	for idx, ar := range a {
		diff := ar - b[idx]
		if diff != 0 {
			if diff > 0 {
				return 1
			} else if diff < 0 {
				return -1
			}
		}
	}

	return 0
}

// LastIndexAny returns the index of the last occurrence in 's' that is in 'any'.
func LastIndexAny(s []rune, any []rune) int {
	if len(any) == 0 {
		return -1
	}
	for idx := len(s) - 1; idx >= 0; idx-- {
		c := s[idx]
		for _, m := range any {
			if c == m {
				return idx
			}
		}
	}
	return -1
}

// LastIndexNotAny returns the index of the last occurrence in 's' that is not in 'any'.
func LastIndexNotAny(s []rune, any []rune) int {
	if len(any) == 0 {
		return -1
	}
	//fmt.Fprintf(os.Stderr, "last !any: '%s' in '%s'...\n", string(any), string(s))
next_rune:
	for idx := len(s) - 1; idx >= 0; idx-- {
		c := s[idx]
		for _, ac := range any {
			//fmt.Fprintf(os.Stderr, "  '%s' != '%s'\n", string(c), string(ac))
			if c == ac {
				continue next_rune
			}
		}
		//fmt.Fprintf(os.Stderr, "  ==> %d\n", idx)
		return idx
	}
	return -1
}

// LastIndexRune returns the index to the last occurrence in 's' that equals 'r'.
func LastIndexRune(s []rune, r rune) int {
	for idx := len(s) - 1; idx >= 0; idx-- {
		if s[idx] == r {
			return idx
		}
	}
	return -1
}

// LastIndexNotRune returns the index to the last occurrence in 's' that does not equal 'r'.
func LastIndexNotRune(s []rune, r rune) int {
	for idx := len(s) - 1; idx >= 0; idx-- {
		if s[idx] != r {
			return idx
		}
	}
	return -1
}

// LastIndexFunc returns the index to the last occurrence in 's' where 'f' returns true.
func LastIndexFunc(s []rune, f TestFunc) int {
	for idx := len(s) - 1; idx >= 0; idx-- {
		if f(s[idx]) {
			return idx
		}
	}
	return -1
}

// LastIndexNotFunc returns the index to the last occurrence in 's' where 'f' returns false.
func LastIndexNotFunc(s []rune, f TestFunc) int {
	for idx := len(s) - 1; idx >= 0; idx-- {
		if !f(s[idx]) {
			return idx
		}
	}
	return -1
}
