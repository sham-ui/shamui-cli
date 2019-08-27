install:
	wgo restore

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

test:
	wgo test -v ./src/{{ shortName }}/... && echo "Tests finished with success"