.PHONY: generate
generate: codes.go

.PHONY: clean
clean:
	rm -f data/*

codes.go: generate/main.go data/language-subtag-registry.txt data/language_region_codes_crist.tsv
	go generate

data/language-subtag-registry.txt:
	curl -o data/language-subtag-registry.txt \
		https://www.iana.org/assignments/language-subtag-registry/language-subtag-registry

data/language_region_codes_crist.tsv:
	curl -o data/language_region_codes_crist.tsv \
		http://www.sean-crist.com/professional/pages/language_region_codes/language_region_codes_crist_201805.txt
