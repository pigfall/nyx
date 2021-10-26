# You have more choices
# Table
* [Download](#Download)
* [Usage](#Usage)
* [Build](#Build)



# [Download](https://github.com/pigfall/yingying/releases)
# Usage
Server

modify config.json, change the port to your wanted listen port, default is 10101

Run it in your server
```
./server.exe --conf config.json

```

Client 

modify config.json , change the server_addr to your server address
```
client.exe -conf config.json
```

# Build
## Prerequires
* go > 1.17, set env GOPATH, and set $GOPATH/bin to your env PATH
```
git clone https://github.com/pigfall/yingying.git
cd yingying
```
Build server
```
cd cmd/server
go build .
```

Build Client
```
cd cmd/client/build
go build .
./build.exe
```
