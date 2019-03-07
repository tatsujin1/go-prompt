package prompt

import "strings"

// Filter is the type to filter the prompt.Choiceion array.
type Filter func([]Choice, string, bool) []Choice

// FilterHasPrefix checks whether the string choices.Text begins with sub.
func FilterHasPrefix(choices []Choice, sub string, ignoreCase bool) []Choice {
	return filterChoiceions(choices, sub, ignoreCase, strings.HasPrefix)
}

// FilterHasSuffix checks whether the completion.Text ends with sub.
func FilterHasSuffix(choices []Choice, sub string, ignoreCase bool) []Choice {
	return filterChoiceions(choices, sub, ignoreCase, strings.HasSuffix)
}

// FilterContains checks whether the completion.Text contains sub.
func FilterContains(choices []Choice, sub string, ignoreCase bool) []Choice {
	return filterChoiceions(choices, sub, ignoreCase, strings.Contains)
}

// FilterFuzzy checks whether the completion.Text fuzzy matches sub.
// Fuzzy searching for "dog" is equivalent to "*d*o*g*". This search term
// would match, for example, "Good food is gone"
//                               ^  ^      ^
func FilterFuzzy(choices []Choice, sub string, ignoreCase bool) []Choice {
	return filterChoiceions(choices, sub, ignoreCase, fuzzyMatch)
}

func fuzzyMatch(s, sub string) bool {
	sChars := []rune(s)
	subChars := []rune(sub)
	sIdx := 0

	for _, c := range subChars {
		found := false
		for ; sIdx < len(sChars); sIdx++ {
			if sChars[sIdx] == c {
				found = true
				sIdx++
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func filterChoiceions(suggestions []Choice, sub string, ignoreCase bool, function func(string, string) bool) []Choice {
	if sub == "" {
		return suggestions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]Choice, 0, len(suggestions))
	for i := range suggestions {
		c := suggestions[i].Text
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if function(c, sub) {
			ret = append(ret, suggestions[i])
		}
	}
	return ret
}
