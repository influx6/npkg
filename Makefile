ci: boot-redis
	go test -v -cover ./...
	go test -v -cover -tags=integration ./...

boot-redis:
	redis-server &
	sleep 5s

build-docker-test:
	docker build -t npkg-docker-tests -f ./dockerfile-test .

run-docker-test:
	docker run -it npkg-docker-tests make ci