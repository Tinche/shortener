IMG=tinche/shortener-frontend

image:
	docker build -t $(IMG) .

docker-run:
	docker run --rm -ti --net host $(img)

push:
	docker tag $(IMG) gcr.io/url-shortener-189609/$(IMG)
	gcloud docker -- push gcr.io/url-shortener-189609/$(IMG)
