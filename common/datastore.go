package common

import (
	"gopkg.in/mgo.v2"
)
// CreateSession creates the main session to our mongodb instance
// 连接数据库
// default func
func CreateSession(host string) (*mgo.Session, error) {
	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}


