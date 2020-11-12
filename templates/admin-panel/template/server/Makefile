restore_deps:
	cd {{ shortName }}
	go mod download

get_go_bindata:
	cd {{ shortName }}
	go get -u github.com/go-bindata/go-bindata/...

install: restore_deps get_go_bindata

.ONESHELL:
build_js:
	cd ../client/
	yarn build

bind_dist: build_js
	go-bindata -o {{ shortName }}/assets/bindata.go -pkg assets -fs -prefix "../client" ../client/dist/ ../client/dist/font ../client/index.html ../client/favicon.ico

build_386: bind_dist
	cd {{ shortName }}
	env GOOS=linux GOARCH=386 go build -o ../bin/linux_386/{{ shortName }} {{ shortName }}

build_amd64: bind_dist
	cd {{ shortName }}
	env GOOS=linux GOARCH=amd64 go build -o ../bin/linux_amd64/{{ shortName }} {{ shortName }}

build: build_386 build_amd64

run_test:
	cd {{ shortName }}
	go test -v ./... -cover -coverprofile=coverage.out -coverpkg=./... && echo "Tests finished with success"

coverage:
	cd {{ shortName }}
	go tool cover -html=coverage.out

test: bind_dist run_test coverage
