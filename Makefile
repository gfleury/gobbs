PKGS = $$(go list ./... | grep -v /vendor/)

default:
	go build

test:
	go clean $(PKGS)
	TZ=UTC go test $(PKGS) -check.v -coverprofile=coverage.txt -covermode=atomic

race:
	go clean $(PKGS)
	go test -race $(PKGS) -check.v -coverprofile=coverage.txt -covermode=atomic

profile:
	go clean $(PKGS)
	make
	
clean:
	rm -rf *.prof
	go clean $(PKGS)

lint:
	docker run --rm -it -v ${PWD}:/go/src/github.com/gfleury/gobbs -w /go/src/github.com/gfleury/gobbs golangci/golangci-lint:latest golangci-lint run
