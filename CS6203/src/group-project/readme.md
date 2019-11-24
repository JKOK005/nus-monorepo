## BugsDB: Building the ultra generic and scalable database
We are pretty sure theres still a long way to go in terms of development effort. If you do find any bugs in the program, please raise it in the issue tracking section rather than deducting marks :0

A detailed description of how we built BugsDB can be found [here](...)

### Description
BugsDB is a simple key value data store. We use RocksDB as our storage layer and consistent hashing for load balancing requests. 
Clients communicate to the database via a series of gRPC calls. 

A request handling example is illustrated below
``` go
if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url, bootstrap_port), grpc.WithInsecure()); err != nil {
	glog.Error(err)
} else {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollTimeOutMs) * time.Millisecond)
	client := pb.NewPutKeyServiceClient(conn)
	defer conn.Close(); defer cancel()
	resp, err := client.PutKey(ctx, &pb.PutKeyMsg{Key: "keyA", Val: []byte("Hello world!"}; err != nil {
		panic(err)
	} else if resp.Ack != true {
		glog.Error("Failed to insert key: ", key)
	}
}
```

We begin by dailing up a BugsDB deployment over at `bootstrap_url:bootstrap_port`. We then instruct the system to insert a key-value (string, btye) datapoint of ("keyA", "Hello world!"). 

To read the message from the database, we simply issue a GetKey request to any deployment. 

``` go
if conn, err := grpc.Dial(fmt.Sprintf("%s:%d", bootstrap_url, bootstrap_port), grpc.WithInsecure()); err != nil {
	glog.Error(err)
} else {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pollTimeOutMs) * time.Millisecond)
	client := pb.NewGetKeyServiceClient(conn)
	defer conn.Close(); defer cancel()
	client.GetKey(ctx, &pb.GetKeyMsg{Key: "keyA"}; err != nil {
		panic(err)
	} else if resp.Ack != true {
		glog.Infof(fmt.Sprintf("Retrieved key: %d, val: %s", key, string(resp.Val)))
	}
}
```

If all goes well, we should see the print out: 
```
Retrieved key: keyA, val: Hello world!
```

| APIs over gRPC | Example |
| -------------  | ------------- |
| PutKey(context, protobuf.PutKeyMsg(key str, val []byte)) | PutKeyMsg{Key: "keyA", Val: []byte("Hello world!"} | 
| GetKey(context, protobuf.GetKeyMsg(key str, val []byte)) | GetKeyMsg{Key: "keyA"}| 


## Installation instructions
BugsDB has the following dependencies
*Go packages*
- samuel/go-zookeeper 	[ref](https://github.com/samuel/go-zookeeper)
- techbot/rocksdb 		[ref](https://github.com/tecbot/gorocksdb)
- grpc 					[ref](https://github.com/grpc/grpc-go)
- protoc-gen-go 		[ref](https://github.com/golang/protobuf)
- proto 				[ref](https://github.com/golang/protobuf)
- glog 					[ref](https://github.com/golang/glog)

*System packages*
- Protocol buffers
- RocksDB 

### Local
To install BugsDB locally, first clone this repo

```
git clone https://github.com/JKOK005/nus-monorepo.git

cd CS6203/src/group-project
```

Install rocksdb library on machine
```shell script
# For Mac users 
# Instructions at: https://brewinstall.org/Install-rocksdb-on-Mac-with-Brew/
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)" < /dev/null 2> /dev/null
brew install rocksdb

# OR 

# For linux users
sudo apt-get -y install librocksdb-dev
```

Install the Protoc compiler (this process will take a while to complete)
``` shell script
curl -L -o /tmp/protobuf.tar.gz https://github.com/google/protobuf/releases/download/v3.7.1/protobuf-cpp-3.7.1.tar.gz
tar xvzf protobuf.tar.gz /tmp
cd /tmp
mkdir export
./autogen.sh && ./configure --prefix=/export && make -j 3 && make check && make install
```

Install Go packages
``` go
# Zookeeper
RUN go get -u github.com/samuel/go-zookeeper/zk

# gRPC
RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/golang/protobuf/proto

# glog
RUN go get -u github.com/golang/glog

# go RocksDB
RUN go get -u github.com/tecbot/gorocksdb
```

Additional linkages might need to be done for Mac OS:
```
CGO_CFLAGS="-I/path/to/rocksdb/include" \
CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" \
```

### Generate protobuf file via:
The next step would be to compile the protocol buffers to generate the binaries
```proto
protoc -I Protobuf/Template/ --go_out=plugins=grpc:Protobuf/Generate Protobuf/Template/raft.proto 
```

### Executing the progamme
Several run time environment variables need to be defined when creating a BugsDB deployment within a hash group. 

| Env. var | Definition | Default |
| -------------  | ------------- | ------------- |
| REGISTER_LISTENER_DNS | The URL which all nodes in the hash group registers themselves as in Zookeeper | localhost |
| REGISTER_LISTENER_PORT | The PORT which all nodes in the hash group registers themselves as in Zookeeper | 8000 |
| LISTENER_ADR | The URL which is used for API calls to the server / client interface. We need to make the distinction between *LISTENER_ADR* and *REGISTER_LISTENER_DNS* owing to the fact that deployments over Kubernetes requre the node to listen to its localhost address, yet has to register its service (in Zookeeper) as its service address | "0.0.0.0" |
| HASH_GROUP | Hash group count of each node. Must be <= CHORD_HASH_SIZE | 1 |
| START_CYCLE_NO | Starting cycle number for RAFT protocol | 0 |
| CYCLES_TO_TIMEOUT | Number of cycles before node times out and requests for election | 10 |
| CHORD_HASH_SIZE | The size of the CHORD ring | 4 |
| STORAGE_LOC | Location for RocksDB storage | ./storage/leader |
| NODE_PATH | Path in Zookeeper which nodes with similar hash groups registers themselves | /nodes |
| FOLLOWER_PATH | Path in Zookeeper which nodes following another in their finger table registers themselves | /followers |
| SERVERS_ZK | List of Zookeeper bootstrap urls | []string{"localhost:2181"} |
| CONN_TIMEOUT | Zookeeper connection timeout | 10 |


#### Deploying a single instance
The simplest deployment is a single instance deployment. This consists of a single BugsDB node with no replication to slaves. 

```go 
# Standard deployment
go run main.go

# Or in debug mode
go run main.go --stderrthreshold=INFO
```

#### Scaling hash groups 
To scale with workload, we can deploy multiple nodes within the consistent hash ring. For a 3 server local deployment

```go
export CHORD_HASH_SIZE=3

# Node 1 - In separate shell session
export REGISTER_LISTENER_PORT=8000
export HASH_GROUP=1
export STORAGE_LOC="./storage/node_1"
go run main.go

# Node 2 - In separate shell session
export REGISTER_LISTENER_PORT=8100
export HASH_GROUP=2
export STORAGE_LOC="./storage/node_2"
go run main.go

# Node 3 - In separate shell session
export REGISTER_LISTENER_PORT=8200
export HASH_GROUP=3
export STORAGE_LOC="./storage/node_3"
go run main.go
```

#### Replication of data 
To attach slaves to a leader in the hash group for replication, we simply create a new session with the _HASH_GROUP_ pointing to the leader's _HASH_GROUP_. 

Assuming that a leader is currently servicing _HASH_GROUP_=1 at PORT=8000

```go
export CHORD_HASH_SIZE=3

# Slave 1 - In separate shell session
export REGISTER_LISTENER_PORT=8100
export HASH_GROUP=1
export STORAGE_LOC="./storage/slave_1"
go run main.go
```

### Deployments on Docker & Kubernetes
BugsDB is on Kubernetes !!!!

#### Docker setup
The first step is to compile both docker images for BugsDB. The following docker images are available:

| Name | Description |
| -------------  | ------------- |
| preSetupDockerFile | Base image for building the main BugsDB image. Sets up the environment for RocksDB and gRPC. *Note* Compiling this image will take a long time to try not to modify anything in the image. |
| Dockerfile | BugsDB docker image. Used to generate the binaries for the system. | 

Compiling the docker files is as simple as:
``` shell script
# Compile preSetupDockerFile
docker build -t <user>/bugsdbinstallation:<version> -f preSetupDockerFile .

# Compile main DockerFile
docker build -t <user>/bugsdb:<version> -f Dockerfile .

# Optional push
docker push <user>/bugsdbinstallation:<version>
docker push <user>/bugsdb:<version> -f Dockerfile .
```
*Note* Edit the first line of `Dockerfile` accordingly to point to the pre installation image that was defined, as it currently points to my docker hub id.

In any case, pre compiled images are present in docker hub
- jkok005/bugsdbinstallation:1.0.0
- jkok005/bugsdb:1.0.0

#### Helm charts
Our helm chart released has been tested on `minikube v0.24.1` and `helm v2.16.0`

To launch the chart, ensure that your Kubernetes cluster is ready and Tiller is active (for versions < helm 3.0.0)

Spin up the helm chart via:
```shell script
helm install --name deploy-bugsdb helm-charts/
```

Tear down the chart via:
```shell script
helm delete deploy-bugsdb --purge
```