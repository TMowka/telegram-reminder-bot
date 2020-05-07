ifeq "$(strip $(shell go env GOARCH))" "amd64"
RACE_FLAG := -race
endif

lint:
	cd /tmp && GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

pretest: lint

gotest:
	go test $(RACE_FLAG) -vet all ./...

test: pretest gotest
