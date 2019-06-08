default: install

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
	GOOS=linux GOARCH=amd64 go build
	docker build -t mjpitz/cluster-preset:latest .

dockerx:
	docker buildx rm mjpitz--cluster-preset || echo "mjpitz--cluster-preset does not exist"
	docker buildx create --name mjpitz--cluster-preset --use
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -f Dockerfile.pub -t mjpitz/cluster-preset:latest .
