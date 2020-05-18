restore_wgo:
	wgo restore

get_go_bindata:
	wgo get -u github.com/go-bindata/go-bindata/...

install: restore_wgo get_go_bindata

.ONESHELL:
build_js:
	cd ../client/
	yarn build

bind_dist: build_js
	./vendor/bin/go-bindata -o src/assets/bindata.go -pkg assets -fs -prefix "../client" ../client/dist/ ../client/dist/font ../client/index.html ../client/favicon.ico

build_386: bind_dist
	env GOOS=linux GOARCH=386 wgo install {{ shortName }}

build_amd64: bind_dist
	env GOOS=linux GOARCH=amd64 wgo install {{ shortName }}

build: build_386 build_amd64

run_test:
	wgo test -v ./src/{{ shortName }}/... -cover -coverprofile=coverage.out -coverpkg ./src/{{ shortName }}/... && echo "Tests finished with success"

coverage:
	wgo tool cover -html=coverage.out

test: bind_dist run_test coverage
