.SILENT: main get-dependencies list-packages build test check-go
.PHONY: get-dependencies list-packages build test
main:
	echo "picel CLI build tool commands:"
	echo "get-dependencies, list-packages, build, test"
get-dependencies: check-go
	if ! which glide &> /dev/null; \
	then >&2 echo "Missing dependency: Glide is required https://glide.sh/"; \
	fi;

	glide install
list-packages:
	go list ./... | grep -v /vendor/
build:
	go build
test:
	go test `go list ./... | grep -v /vendor/`
check-go:
	if ! which go &> /dev/null; \
	then >&2 echo "Missing dependency: Go is required https://golang.org/"; \
	fi;
