all:
	go install -v github.com/frankbraun/kitchensink/...

.PHONY: test
test:
	go test github.com/frankbraun/kitchensink/...
