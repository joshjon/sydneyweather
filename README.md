# ☁️ Sydney Weather

An HTTP Service that reports on the weather in Sydney using [weatherstack](https://weatherstack.com)
and [OpenWeather](https://openweathermap.org).

## 🧰 Tools Used

- Go 1.18.1
- Docker 20.10.10 CE
- Make

## 🚀 Running

Before proceeding, please ensure you have Docker running and that you have set `WEATHER_STACK_KEY`
and `OPEN_WEATHER_KEY` environment variables e.g. `export OPEN_WEATHER_KEY=some-key`.

1. Start the server

    ```shell
    make start
    ``` 

2. Make a request

    ```shell
    curl http://localhost:8080/v1/weather?city=sydney
    ```

3. Stop the server

    ```shell
    make stop

## 🔬 Testing

- Unit tests

   ```shell
   make unit
   ```

- Integration test

   ```shell
   make integration
   ```

## 🌩 Tradeoffs

The project was time boxed to the recommended duration. Below are some trade-offs and improvements that could have been
made if more time was permitted.

- Only Sydney weather is supported as per the requirements, but with slight modification any city can easily be
  accepted.
- The service was not deployed anywhere. The next step would have been creating a new service deployment using
  Kubernetes with multiple replicas for high availability.
- API Key secrets are read from environment variables. If the service was deployed, it would be ideal to use a secret
  manager to store and retrieve the keys e.g. Google Secret Manager.
- If more weather source clients were required, they could have been implemented in a slightly more generic manner.
  There is currently a primary and fail over client which are accessed directly to get weather data. This could be
  improved by accepting an ordered list of weather clients that implement the same interface, which are then looped and
  until a successful weather data response is returned. This would make it extremely easy to add/remove more weather
  sources to further mitigate failure or delivery of stale results.
