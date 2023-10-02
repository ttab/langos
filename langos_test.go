package langos_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ttab/langos"
)

func ExampleGetLanguage() {
	info, err := langos.GetLanguage("en-gb")
	if err != nil {
		panic(err)
	}

	knownCombo := langos.IsKnownCombination(info.Code)

	fmt.Printf("code: %s\n", info.Code)
	fmt.Printf("language: %s - %s\n", info.Language, info.LanguageName)
	fmt.Printf("region: %s - %s\n", info.Region, info.RegionName)
	fmt.Printf("is known combination: %v", knownCombo)

	// Output:
	// code: en-GB
	// language: en - English
	// region: GB - United Kingdom
	// is known combination: true
}

type getCase struct {
	Input   string
	Invalid bool
	Unknown bool
	Info    langos.LanguageInformation
}

func TestGetLanguage(t *testing.T) {
	cases := map[string]getCase{
		"UppercaseLanguage": {
			Input: "SV",
			Info: langos.LanguageInformation{
				Code:         "sv",
				Language:     "sv",
				LanguageName: "Swedish",
			},
		},
		"LowercaseLanguage": {
			Input: "ta",
			Info: langos.LanguageInformation{
				Code:         "ta",
				Language:     "ta",
				LanguageName: "Tamil",
			},
		},
		"AllLower": {
			Input: "sv-fi",
			Info: langos.LanguageInformation{
				Code:         "sv-FI",
				Language:     "sv",
				LanguageName: "Swedish",
				HasRegion:    true,
				Region:       "FI",
				RegionName:   "Finland",
			},
		},
		"UnknownLanguage": {
			Input:   "xz",
			Invalid: true,
		},
		"UnknownRegion": {
			Input:   "sv-SW",
			Invalid: true,
		},
		"UnknownCombination": {
			Input:   "sv-RE",
			Unknown: true,
			Info: langos.LanguageInformation{
				Code:         "sv-RE",
				Language:     "sv",
				LanguageName: "Swedish",
				HasRegion:    true,
				Region:       "RE",
				RegionName:   "Réunion",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			info, err := langos.GetLanguage(tc.Input)

			switch {
			case !tc.Invalid && err != nil:
				t.Fatalf("invalid language code: %v", err)
			case tc.Invalid && err == nil:
				t.Fatalf("expected %q to be an invalid language code", tc.Input)
			}

			if err != nil {
				return
			}

			if diff := cmp.Diff(tc.Info, info); diff != "" {
				t.Fatalf("GetLanguage() info mismatch (-want +got):\n%s",
					diff)
			}

			known := langos.IsKnownCombination(tc.Input)

			switch {
			case tc.Unknown && known:
				t.Fatal("expected the language region combination to be unknown")
			case !tc.Unknown && !known:
				t.Fatal("expected the language region combination to be known")
			}
		})
	}
}
