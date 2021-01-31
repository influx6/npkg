ci: boot-redis
	go test -v -cover ./...
	go test -v -cover -tags=integration ./...

down:
	docker-compose down 

up:
	docker-compose up -d

build-docker-ci:
	docker build -t npkg-docker-tests -f ./dockerfile-test .

docker-ci: build-docker-ci
	docker run -it npkg-docker-tests make ci
