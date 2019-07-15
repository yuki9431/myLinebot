package main

import "gopkg.in/mgo.v2"

type mongoDb struct {
	dial string
	name string
	db   *mgo.Database
}

type MongoDb interface {
	disconnectDb()
	insertDb(interface{}, string) error
	removeDb(interface{}, string) error
	searchDb(interface{}, string) error
}

func NewMongo(dial string, name string) (MongoDb, error) {
	session, err := mgo.Dial(dial)
	db := session.DB(name)

	return &mongoDb{
		dial: dial,
		name: name,
		db:   db,
	}, err
}

// mongoDB切断
func (m *mongoDb) disconnectDb() {
	m.db.Session.Close()
}

// mongoDB挿入
func (m *mongoDb) insertDb(obj interface{}, colectionName string) (err error) {
	defer m.disconnectDb()

	col := m.db.C(colectionName)
	return col.Insert(obj)
}

// mongoDB削除
func (m *mongoDb) removeDb(obj interface{}, colectionName string) (err error) {
	defer m.disconnectDb()

	col := m.db.C(colectionName)
	_, err = col.RemoveAll(obj)
	return
}

// mondoDB抽出
func (m *mongoDb) searchDb(obj interface{}, colectionName string) (err error) {
	defer m.disconnectDb()

	col := m.db.C(colectionName)
	return col.Find(nil).All(obj)
}
