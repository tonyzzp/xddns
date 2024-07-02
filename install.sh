echo start build...
go build -ldflags "-s -w" -trimpath .
echo build complete
echo copy to /usr/bin
sudo cp xddns /usr/bin/xddns
sudo chmod +x /usr/bin/xddns

