.PHONY: stress run build clean logs stop

stress:
	go test -v ./test/... -timeout=10m

run:
	docker-compose up -d --build

build:
	docker-compose build

clean:
	docker-compose down

logs:
	docker-compose logs -f

stop:
	docker-compose down

restart: stop run

status:
	docker-compose ps