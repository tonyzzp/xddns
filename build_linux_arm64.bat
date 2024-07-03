set goos=linux
set goarch=arm64
go build -ldflags="-s -w" -trimpath -v .