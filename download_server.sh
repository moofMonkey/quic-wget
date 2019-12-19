wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.13.5.linux-amd64.tar.gz
rm -rf go1.13.5.linux-amd64.tar.gz
wget https://github.com/MoofMonkey/quic-wget/raw/1.0/server.go
export GOROOT=/usr/local/go/bin
export PATH=$PATH:$GOROOT
go get github.com/lucas-clemente/quic-go
go build server.go
chmod +x server # in case it won't set execute flag
