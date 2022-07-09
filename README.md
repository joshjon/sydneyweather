# ☁️ Sydney Weather

An HTTP Service that reports on the weather in Sydney using [weatherstack](https://weatherstack.com)
and [OpenWeather](https://openweathermap.org).

## 🚀 Running

1. Start the server 

    ```shell
    go run main.go
    ``` 

2. Make a request

    ```shell
    curl http://localhost:8080/v1/weather?city=sydney
    ```

## 🔬 Testing

```shell
go test -count=1 ./...
```

## 🧰 Tools Used

- Go 1.18.1
