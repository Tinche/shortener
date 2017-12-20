IMG=tinche/shortener-go

build:
	go build -o shortener src/shortener/*.go

run:
	./shortener

image:
	docker build -t $(IMG) .

docker-run:
	docker run --rm -ti -e REDIS_URL=127.0.0.1:6379 --net host $(IMG)

push:
	docker tag $(IMG) gcr.io/url-shortener-189609/$(IMG)
	gcloud docker -- push gcr.io/url-shortener-189609/$(IMG)