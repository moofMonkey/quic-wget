# QUIC-wget

## Description
QUIC-wget is utility to transfer files over QUIC (which uses UDP, and uses checksums).

QUIC-wget was built to download files (especially backups) over unstable network with frequent
timeouts/low bandwidth. Usually it's 2-3x faster than TCP/SFTP connection for us.

## Server installation & start (prebuilt)

### Linux
```sh
wget https://github.com/MoofMonkey/quic-wget/releases/download/2.2/quic-wget
chmod +x ./quic-wget
./quic-wget --password="your_super_secret_password" --target="0.0.0.0:12345"
```
You must change 12345 (port) and password example to some random generated ones.

You can also run server in screen by adding "screen " in front of "./quic-wget".
It is highly recommended doing this on unstable connections.

### Windows

1. Download [quic-wget executable](https://github.com/MoofMonkey/quic-wget/releases/download/2.2/quic-wget.exe)
2. Open cmd/PowerShell in the same folder (or open it by Win+X and cd to folder with executable)
3.
```sh
.\quic-wget --password="your_super_secret_password" --target="0.0.0.0:12345"
```
You must change 12345 (port) and password example to some random generated ones.

## Client installation & start (prebuilt)

### Linux
```sh
wget https://github.com/MoofMonkey/quic-wget/releases/download/2.2/quic-wget
chmod +x ./quic-wget
./quic-wget --password="your_super_secret_password" --target="8.8.8.8:12345" --downloadPath="/backup.tar.bz2" --localPath="backup.tar.bz2"
```
You must change 12345 (port) and password example to ones used on server.

### Windows

1. Download [quic-wget executable](https://github.com/MoofMonkey/quic-wget/releases/download/2.2/quic-wget.exe)
2. Open cmd/PowerShell in the same folder (or open it by Win+X and cd to folder with executable)
3.
```sh
.\quic-wget --password="your_super_secret_password" --target="8.8.8.8:12345"  --downloadPath="/backup.tar.bz2" --localPath="backup.tar.bz2"
```
You must change 12345 (port) and password example to some random generated ones.
