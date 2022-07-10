.PHONY: start stop unit

start:
	docker build -t local/sydneyweather .
	docker run --rm --name sydneyweather -p 8080:8080 -d \
		-e WEATHER_STACK_KEY=${WEATHER_STACK_KEY} \
		-e OPEN_WEATHER_KEY=${OPEN_WEATHER_KEY} \
		local/sydneyweather -config config.yaml

stop:
	docker stop sydneyweather

unit:
	go test -count=1 ./...
