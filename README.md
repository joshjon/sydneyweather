# ☁️ Sydney Weather

An HTTP Service that reports on the weather in Sydney using [weatherstack](https://weatherstack.com)
and [OpenWeather](https://openweathermap.org).

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

## 🧰 Tools Used

- Go 1.18.1
- Docker 20.10.10 CE
- Make
