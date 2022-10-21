
test:
	go test -v -race ./...

test-short:
	go test -short -race ./...
# todo: COVERPROFILE
