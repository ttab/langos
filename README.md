# Langos

[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)][godev]
[![Build Status](https://github.com/ttab/langos/actions/workflows/test.yaml/badge.svg?branch=main)][actions]

Langos validates a language (ISO 639-1) region (ISO 3166-1 alpha-2) combinations in the format "language[-region]" and returns basic information about the language and region.

Example:

``` go
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
```

## Data sources

Langos uses the languages and regions defined in the [IANA language subtag registry](https://www.iana.org/assignments/language-subtag-registry/language-subtag-registry) together with the the [list of known language and region codes](http://www.sean-crist.com/professional/pages/language_region_codes/index.html) compiled by Sean Christ. 

[godev]: https://pkg.go.dev/github.com/ttab/langos
[actions]: https://github.com/ttab/langos/actions
