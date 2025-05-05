# memo

``` windows
go build -o ../demo/custom_lib.dll -buildmode=c-shared .
```

``` linux
docker run -ti -v 專案位置:/app -w /app golang:1.24-bullseye bash
go build -o ../demo/custom_lib.so -buildmode=c-shared .
```
