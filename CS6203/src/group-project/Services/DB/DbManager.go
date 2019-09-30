package DB

import (
	"github.com/golang/glog"
	"group-project/Utils"
	pb "group-project/Protobuf/Generate"
)

type DbManager struct {
	DbCli *Utils.RocksDbClient
}

func (d *DbManager) putKey(key string, val []byte) (bool, error) {
	/*
	Inserts key into DB
	*/
	if err := d.DbCli.Put(key, val); err != nil {
		glog.Fatal(err)
		return false, err
	}
	return true, nil
}

func (d *DbManager) getKey(key string) ([]byte, bool) {
	/*
	Attempts to get key. If key does not exist, returns false
	*/
	val, _ := d.DbCli.Get(key)
	if len(val) > 0 {return nil, false}
	return val, true
}

func (d *DbManager) putKeyRoutine() {
	for {
		select {
		case msg := <- Utils.PutKeyChannel.ReqCh:
			if resp, err := d.putKey(msg.Key, msg.Val); err != nil {
				Utils.PutKeyChannel.RespCh <- false
			} else {
				Utils.PutKeyChannel.RespCh <- resp
			}
		}
	}
}

func (d *DbManager) getKeyRoutine() {
	for {
		select {
		case msg := <- Utils.GetKeyChannel.ReqCh:
			val, isPresent := d.getKey(msg)
			Utils.GetKeyChannel.RespCh <- &pb.GetKeyResp{Ack:isPresent, Val:val}
		}
	}
}

func (d DbManager) Start() {
	glog.Info("Db Manager started")
	go d.getKeyRoutine()
	go d.putKeyRoutine()
}