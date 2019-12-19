# quic-wget

## Server installation & start (prebuilt)

### Linux
```sh
wget https://github.com/MoofMonkey/quic-wget/releases/download/1.0/server
chmod +x ./server
./server --password="your_super_secret_password" --target="0.0.0.0:12345"
```
You must change 12345 (port) and password example to some random generated ones.

Also you can run server in screen by adding "screen " in front of "./server".
It highly recommended to do this on unstable connections.

### Windows

1. Download [server executable](https://github.com/MoofMonkey/quic-wget/releases/download/1.0/server.exe)
2. Open cmd/PowerShell in the same folder (or open it by Win+X and cd to folder with executable)
3.
```sh
.\server --password="your_super_secret_password" --target="0.0.0.0:12345"
```
You must change 12345 (port) and password example to some random generated ones.

## Client installation & start (prebuilt)

### Linux
```sh
wget https://github.com/MoofMonkey/quic-wget/releases/download/1.0/client
chmod +x ./client
./client --password="your_super_secret_password" --target="8.8.8.8:12345" --downloadPath="/backup.tar.bz2" --localPath="backup.tar.bz2"
```
You must change 12345 (port) and password example to ones used on server.

### Windows

1. Download [client executable](https://github.com/MoofMonkey/quic-wget/releases/download/1.0/client.exe)
2. Open cmd/PowerShell in the same folder (or open it by Win+X and cd to folder with executable)
3.
```sh
.\client --password="your_super_secret_password" --target="8.8.8.8:12345"  --downloadPath="/backup.tar.bz2" --localPath="backup.tar.bz2"
```
You must change 12345 (port) and password example to some random generated ones.


## Server installation & start (***NOT*** prebuilt)

### Linux
```sh
wget https://github.com/MoofMonkey/quic-wget/raw/master/download_server.sh
chmod +x download_server.sh
./download_server.sh
./server --password="your_super_secret_password" --target="0.0.0.0:12345"
```
You must change 12345 (port) and password example to some random generated ones.

Also you can run server in screen by adding "screen " in front of "./server".
It highly recommended to do this on unstable connections.

### Windows

1. Download and install Go from https://golang.org/dl/
1. Download [server source code](https://github.com/MoofMonkey/quic-wget/raw/master/server.go)
2. Open cmd/PowerShell in the same folder (or open it by Win+X and cd to folder with executable)
3.
```sh
go get github.com/lucas-clemente/quic-go
go build server.go
.\server --password="your_super_secret_password" --target="0.0.0.0:12345"
```
You must change 12345 (port) and password example to some random generated ones.

## Client installation & start (***NOT*** prebuilt)

### Linux
```sh
wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.13.5.linux-amd64.tar.gz
rm -rf go1.13.5.linux-amd64.tar.gz
wget https://github.com/MoofMonkey/quic-wget/raw/master/client.go
export GOROOT=/usr/local/go/bin
export PATH=$PATH:$GOROOT
go get github.com/lucas-clemente/quic-go
go build client.go
chmod +x client # in case it won't set execute flag
./client --password="your_super_secret_password" --target="8.8.8.8:12345" --downloadPath="/backup.tar.bz2" --localPath="backup.tar.bz2"
```
You must change 12345 (port) and password example to some random generated ones.

Also you can run server in screen by adding "screen " in front of "./server".
It highly recommended to do this on unstable connections.

### Windows

1. Download and install Go from https://golang.org/dl/
1. Download [server source code](https://github.com/MoofMonkey/quic-wget/raw/master/client.go)
2. Open cmd/PowerShell in the same folder (or open it by Win+X and cd to folder with executable)
3.
```sh
go get github.com/lucas-clemente/quic-go
go build server.go
.\client --password="your_super_secret_password" --target="8.8.8.8:12345" --downloadPath="/backup.tar.bz2" --localPath="backup.tar.bz2"
```
You must change 12345 (port) and password example to some random generated ones.
