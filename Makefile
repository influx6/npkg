ci:
	go test -v -cover ./...
	go test -v -cover -tags=integration ./...

build-docker-test:
	docker build -t npkg-docker-tests -f ./dockerfile-test .

run-docker-test:
	docker run npkg-docker-tests