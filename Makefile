.PHONY: gen-pdf

gen-pdf:
	npx --yes md-to-pdf $(FILE)
