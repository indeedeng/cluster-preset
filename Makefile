default: install

build-deps:
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups

fmt:
	go-groups -w .
	gofmt -s -w .

deps:
	go get -v ./...

test:
	go vet ./...
	golint -set_exit_status ./...
	go test -v ./...

install:
	go install

deploy:
	mkdir -p bin
	gox -os="windows linux darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	gox -os="linux" -arch="arm" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	GOOS=linux GOARCH=arm64 go build -o bin/cluster-preset_linux_arm64

docker:
	docker build -t indeedoss/cluster-preset:latest .

dockerx:
	docker buildx rm indeedoss--cluster-preset || echo "indeedoss--cluster-preset does not exist"
	docker buildx create --name indeedoss--cluster-preset --use
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -f Dockerfile.pub -t indeedoss/cluster-preset:latest .
