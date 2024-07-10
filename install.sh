
echo start build and install
go install -ldflags="-s -w" -trimpath -v .
