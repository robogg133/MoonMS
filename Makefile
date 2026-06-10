GOOS ?=""
GOARCH ?=""

CURRENT_COMMIT = $(shell git rev-parse HEAD)



build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build .

release:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -trimpath -buildvcs=true -ldflags="-s -w -X 'main.AdittionalVersionData=(CommitID: $(CURRENT_COMMIT))'" .
