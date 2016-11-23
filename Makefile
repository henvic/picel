.SILENT: get-dependencies list-packages build build-with-pprof test check-go
.PHONY: get-dependencies list-packages build build-with-pprof test check-go
build:
	go build
build-with-pprof:
	go build -tags=pprof
get-dependencies: check-go
	if ! which glide &> /dev/null; \
	then >&2 echo "Missing dependency: Glide is required https://glide.sh/"; \
	fi;

	glide install
list-packages:
	go list ./... | grep -v /vendor/
test:
	echo "Checking for unchecked errors."
	errcheck $(go list ./... | grep -v /vendor/)
	echo "Linting code."
	test -z "$(golint ./... | grep -v "^vendor" | tee /dev/stderr)"
	echo "Examining source code against code defect."
	go vet $(go list ./... | grep -v /vendor/)
	echo "Running tests (may take a while)."
	go test `go list ./... | grep -v /vendor/`
check-go:
	if ! which go &> /dev/null; \
	then >&2 echo "Missing dependency: Go is required https://golang.org/"; \
	fi;
