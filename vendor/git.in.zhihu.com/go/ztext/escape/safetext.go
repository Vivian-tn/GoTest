package escape

import "unicode"

var (
	//emptyRunes   = []rune{}
	spaceRunes = []rune{' '}
	//newlineRunes = []rune{'\n'}
)

// RuneMapper is a rune processing function, it takes a []rune as input,
// offset as starting pointer, limit as ending point, returns processed []rune
// and processed rune length
type RuneMapper func(text []rune, offset int, limit int) ([]rune, int)

func safeTextMapper(text []rune, offset int, limit int) (mapped []rune, step int) {
	codePoint := text[offset]
	step = 1
	switch codePoint {
	case 0x5c: /* \ */
		rest := limit - offset - 1
		if rest >= 3 {
			/* test if it's \x3c */
			if unicode.ToLower(text[offset+1]) == 'x' && text[offset+2] == '3' &&
				unicode.ToLower(text[offset+3]) == 'c' {
				step = 4
			} else if rest >= 5 {
				/* test if it's \u003c */
				if unicode.ToLower(text[offset+1]) == 'u' && text[offset+2] == '0' &&
					text[offset+3] == '0' && text[offset+4] == '3' &&
					unicode.ToLower(text[offset+5]) == 'c' {
					step = 6
				} else {
					return []rune{codePoint}, 1
				}
			} else {
				return []rune{codePoint}, 1
			}
		} else {
			return []rune{codePoint}, 1
		}
		fallthrough
	case 0xc:
		fallthrough
	case 0x2003:
		fallthrough
	case 0xa0:
		fallthrough
	case 0xb:
		return spaceRunes, step
	default:
		if codePoint >= 0x600 && codePoint <= 0x6ff {
			return spaceRunes, 1
		}
		// else branch
		return []rune{codePoint}, 1
	}
}

func safeTextBMPMapper(text []rune, offset int, limit int) (mapped []rune, step int) {
	codePoint := text[offset]
	if codePoint > 0xFFFF {
		return spaceRunes, 1
	} else if codePoint >= 0xd800 && codePoint <= 0xdfff {
		return spaceRunes, 1
	}
	return safeTextMapper(text, offset, limit)
}

// SafeText returns input replaced with safe character.
func SafeText(value string) string {
	return MapText(value, safeTextMapper)
}

// SafeTextBMP returns input replaced with safe character, it's a safer
// version than SafeText, for it only includes Basic Multilingual Plane(Plane
// 0) and excludes Plane 1-16. And in BMP itself, it excludes 0xd800-0xdfff,
// which is surrogate pairs in UTF-16.
// check https://en.wikipedia.org/wiki/Plane_(Unicode) for more information.
func SafeTextBMP(value string) string {
	return MapText(value, safeTextBMPMapper)
}

// MapText maps mapper function to value and return mapped version value
func MapText(value string, mapper RuneMapper) string {
	if len(value) == 0 {
		return value
	}
	input := []rune(value)
	ilen := len(input)
	output := make([]rune, ilen)
	olen := 0
	ocapacity := ilen
	for idx := 0; idx < ilen; {
		toCodePoints, step := mapper(input, idx, ilen)
		toCodePointsLen := len(toCodePoints)
		if toCodePointsLen > 0 {
			if ocapacity < olen+toCodePointsLen {
				ocapacity += (2 * toCodePointsLen)
				newOutput := make([]rune, ocapacity)
				copy(newOutput[:olen], output[:olen])
				output = newOutput
			}
			copy(output[olen:olen+toCodePointsLen], toCodePoints)
			olen += toCodePointsLen
		}
		idx += step
	}
	return string(output[:olen])
}

func doCollapseSpacesSafeTextMapper(text []rune, offset int, limit int, continuation RuneMapper) (mapped []rune, step int) {
	idx := offset
	for ; idx < limit && text[idx] == ' '; idx++ {
	}
	if idx != offset {
		return spaceRunes, idx - offset
	}
	return continuation(text, offset, limit)
}

func collapseSpacesSafeTextMapper(text []rune, offset int, limit int) (mapped []rune, step int) {
	return doCollapseSpacesSafeTextMapper(text, offset, limit, safeTextMapper)
}

func collapseSpacesSafeTextBMPMapper(text []rune, offset int, limit int) (mapped []rune, step int) {
	return doCollapseSpacesSafeTextMapper(text, offset, limit, safeTextBMPMapper)
}

// CollapseSpacesSafeText is SafeText with spaces collapsed
func CollapseSpacesSafeText(value string) string {
	return MapText(value, collapseSpacesSafeTextMapper)
}

// CollapseSpacesSafeTextBMP is SafeTextBMP with spaces collapsed
func CollapseSpacesSafeTextBMP(value string) string {
	return MapText(value, collapseSpacesSafeTextBMPMapper)
}
