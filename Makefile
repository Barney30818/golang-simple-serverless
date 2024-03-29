.PHONY: build clean deploy gomodgen

build: gomodgen
	go mod tidy
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -mod=readonly -ldflags="-s -w" -o bin/sendEmail sendEmail/main.go
	env GOARCH=amd64 GOOS=linux go build -mod=readonly -ldflags="-s -w" -o bin/createMember createMember/main.go
	env GOARCH=amd64 GOOS=linux go build -mod=readonly -ldflags="-s -w" -o bin/getMember getMember/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
