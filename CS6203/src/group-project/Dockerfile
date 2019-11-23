####### Begin BugsDB installation
FROM jkok005/bugsdbinstallation:1.0.0 AS BASELAYER
ARG dir=/app
WORKDIR ${dir}

ENV GOPATH=${dir}
ENV GOBIN=${dir}/bin
ENV PATH=$PATH:${GOBIN}:${GOPATH}

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

# Copy BugsDB & Dependencies into src for building
COPY . ${dir}/src/group-project

# Compile Protobuf
RUN /export/bin/protoc -I src/group-project/Protobuf/Template/ --go_out=plugins=grpc:src/group-project/Protobuf/Generate \
    src/group-project/Protobuf/Template/raft.proto

RUN go install src/group-project/main.go

####### Build new image containing only BugsDB binary
FROM jkok005/bugsdbinstallation:1.0.0
ARG source_dir=/app
ARG new_dir=/app
COPY --from=BASELAYER  ${source_dir}/bin/main ${new_dir}/main
# Set as default port for server
EXPOSE 8000
CMD ["/app/main"]