package DB

import (
	"github.com/golang/glog"
	"group-project/Utils"
	pb "group-project/Protobuf/Generate"
)

type DbManager struct {
	DbCli *Utils.RocksDbClient
}

func (d *DbManager) PutKeyRoutine(msg *pb.PutKeyMsg) (bool, error) {
	/*
	Inserts key into DB
	*/
	if err := d.DbCli.Put(msg.Key, msg.Val); err != nil {
		glog.Fatal(err)
		return false, err
	}
	return true, nil
}

func (d *DbManager) GetKeyRoutine(key *pb.GetKeyMsg) error {
	/*
	Attempts to get key
	*/
}