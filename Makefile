
FILES = $(shell find . -type f -name '*.go' -not -path './vendor/*')

gofmt:
	@gofmt -s -w $(FILES)
	@gofmt -r '&α{} -> new(α)' -w $(FILES)
	@impsort . -p github.com/altipla-consulting/directus-go/v2

lint:
	go install ./...
	go vet ./...
	linter ./...

test:
	go test -v -race ./...
