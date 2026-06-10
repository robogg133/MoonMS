CURRENT_COMMIT = $(shell git rev-parse HEAD)

build:
	go build .

release:
	go build -trimpath -buildvcs=true -ldflags="-s -w -X 'main.AdittionalVersionData=(CommitID: $(CURRENT_COMMIT))'" .
