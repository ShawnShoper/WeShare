package mongo

import (
	mgo "gopkg.in/mgo.v2"
)

type MgoOpera struct {
	session *mgo.Session
	mgodb   *mgo.Database
	mgocol  *mgo.Collection
}

func (mgoo *MgoOpera) checkSession() {
	mgoo.session.Refresh()
}
func (m *MgoOpera) UseDB(db string) mgo.Database {
	m.mgodb = m.session.DB(db)
}
