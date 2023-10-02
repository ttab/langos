package langos

import (
	"fmt"
	"strings"
)

// LanguageInformation contains information about a language code.
type LanguageInformation struct {
	// Code in the canonical xx[-YY] format.
	Code string
	// Language (ISO 639-1) code.
	Language string
	// LanguageName as defined by IANA.
	LanguageName string
	// HasRegion is true if a region was specified for the language.
	HasRegion bool
	// Region is the uppercase ISO 3166-1 alpha-2 code.
	Region string
	// RegionName as defined by IANA.
	RegionName string
}

// GetLanguage validates a language (ISO 639-1) region (ISO 3166-1 alpha-2)
// combination in the format "language[-region]" and returns basic information
// about the language.
func GetLanguage(code string) (LanguageInformation, error) {
	lc, rc, hasRegion := strings.Cut(code, "-")

	lc = strings.ToLower(lc)

	lName, ok := languageMap[lc]
	if !ok {
		return LanguageInformation{}, fmt.Errorf("unknown language code %q", lc)
	}

	if !hasRegion {
		return LanguageInformation{
			Code:         lc,
			Language:     lc,
			LanguageName: lName,
		}, nil
	}

	rc = strings.ToUpper(rc)

	rName, ok := regionMap[rc]
	if !ok {
		return LanguageInformation{}, fmt.Errorf("unknown region code %q", rc)
	}

	return LanguageInformation{
		Code:         lc + "-" + rc,
		Language:     lc,
		LanguageName: lName,
		HasRegion:    true,
		Region:       rc,
		RegionName:   rName,
	}, nil
}

// IsKnownCombination checks a language code against known language region
// pairs. Returns true if the language is valid and no region is specified or
// the language-region combination is known.
func IsKnownCombination(code string) bool {
	lc, rc, hasRegion := strings.Cut(code, "-")

	lc = strings.ToLower(lc)

	_, ok := languageMap[lc]
	if !ok {
		return false
	}

	if !hasRegion {
		return true
	}

	rc = strings.ToUpper(rc)

	return pairMap[lc+"-"+rc]
}
