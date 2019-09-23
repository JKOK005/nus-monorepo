### Installation instructions
Install rocksdb library on machine
```shell script
# Instructions at: https://brewinstall.org/Install-rocksdb-on-Mac-with-Brew/
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)" < /dev/null 2> /dev/null
brew install rocksdb
```

Install rocksdb Go library. Documentation can be found [here](https://godoc.org/github.com/tecbot/gorocksdb)
```go
CGO_CFLAGS="-I/path/to/rocksdb/include" \
CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" \
go get github.com/tecbot/gorocksdb
```

### Generate protobuf file via:
```proto
protoc -I Protobuf/Template/ --go_out=plugins=grpc:Protobuf/Generate Protobuf/Template/raft.proto 
```