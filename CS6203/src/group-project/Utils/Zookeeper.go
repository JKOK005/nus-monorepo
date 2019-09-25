package Utils

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/golang/glog"
	"strings"
	"time"
)

var (
	node_path 			string 		= GetEnvStr("NODE_PATH", "nodes")
	follower_path 		string 		= GetEnvStr("FOLLOWER_PATH", "followers")
	servers_zk 			[]string 	= GetEnvStrSlice("SERVERS_ZK", []string{"localhost:2181"})
	conn_timeout 		int 		= GetEnvInt("CONN_TIMEOUT", 10)
)

type SdClient struct {
	zk_servers 	[]string
	zk_root    	string
	conn      	*zk.Conn
}

func NewZkClient() (*SdClient, error) {
	/* Registers node with ZK cluster */
	glog.Info("Creating client to ZK server", servers_zk)
	client := new(SdClient)
	conn, _, err := zk.Connect(servers_zk, time.Duration(conn_timeout) * time.Second)
	if err != nil {return nil, err}

	glog.Info("Successfully connected to ZK at", servers_zk)
	client.conn = conn
	return client, nil
}

func (s SdClient) checkPathExists(path string) (bool, error) {
	exists, _, err := s.conn.Exists(path)
	if err != nil {return false, err}
	return exists, nil
}

func (s SdClient) constructNode(path string, data []byte) error {
	/*
		Checks if node exists at path
		Else attempts to create a permanent node with data
	*/
	if exists, err := s.checkPathExists(path); err != nil {
		return err
	} else if exists == false {
		glog.Info("Attempting to create node at path", path)
		_, err := s.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s SdClient) constructEphemeralNode(path string, data []byte) error {
	/*
	Constructs an ephemeral node at path with data
	*/
	glog.Info("Creating Ephemeral node at path ", path)
	_, err := s.conn.CreateProtectedEphemeralSequential(path, data, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {return err}
	return nil
}

func (s SdClient) ConstructNodesInPath(path string, delimiter string, data []byte) error {
	/*
		Creates a ZK path of nested nodes from path and delimiter
		If the node created is not the end node, we will populate its data with an empty []byte
		If the node created is the end node, we will populate its data with given data
		Path 		- /distributed_trace/nodes
		Demiliter 	- '/'
	*/
	var err error
	pathSlice := strings.Split(path, delimiter)
	pathTrace := "/"
	for _, eachPath := range pathSlice[ : len(pathSlice)-1] {
		pathTrace = pathTrace + eachPath
		if err = s.constructNode(pathTrace, nil); err != nil {return err}
	}
	if err = s.constructNode(path, data); err != nil {return err}
	return nil
}

func (s SdClient) PrependNodePath(fromPath string) string {
	if fromPath == "" {return fmt.Sprint("/%s", node_path)}
	return fmt.Sprintf("/%s/%s", node_path, fromPath)
}

func (s SdClient) PrependFollowerPath(fromPath string) string {
	if fromPath == "" {return fmt.Sprint("/%s", follower_path)}
	return fmt.Sprintf("/%s/%s", follower_path, fromPath)
}

func (s SdClient) GetNodePaths(from_path string) ([]string, error) {
	glog.Info("GetNodePaths called at", from_path)
	childs, _, err := s.conn.Children(from_path)
	if err != nil {return nil, err}
	return childs, nil
}

func (s SdClient) GetNodeValue (from_path string) ([]byte, error) {
	/* Passes in a node paths and returns the value of the node */
	data, _, err := s.conn.Get(from_path)
	if err != nil {return nil, err}
	return data, nil
}

func (s SdClient) SetNodeValue (from_path string, data []byte) error {
	/* Sets the node value from from_path */
	_, err := s.conn.Set(from_path, data, 0)
	return err
}

func (s SdClient) RegisterEphemeralNode(client_path string, data []byte) error {
	/*
		Registers ephemeral node at client_path with data
		If intermediate paths do not exist, we simply create them as a permanent node with empty data
	*/
	glog.Info("Registering worker ephemeral node address at ", client_path)
	full_path_without_last_slice := strings.Split(client_path, "/")
	full_path_without_last := strings.Join(full_path_without_last_slice[ : len(full_path_without_last_slice)-1], "/")
	if err := s.ConstructNodesInPath(full_path_without_last, "/", nil); err != nil {return err}
	if err := s.constructEphemeralNode(client_path, data); err != nil {return err}
	glog.Info("Created")
	return nil
}