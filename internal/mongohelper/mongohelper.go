package mongohelper

import (
	"math/rand"
	"time"

	"gopkg.in/mgo.v2"
)

type mongoDb struct {
	dial string
	name string
	db   *mgo.Database
}

// MongoDb インターフェース定義
type MongoDb interface {
	DisconnectDb()
	InsertDb(interface{}, string) error
	RemoveDb(interface{}, string) error
	SearchDb(interface{}, interface{}, string) error
	RandomSearchDb(interface{}, string) error
	UpdateDb(interface{}, interface{}, string) error
	Count(colectionName string) (n int, err error)
}

// NewMongo mongoDbを返す
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
func (m *mongoDb) DisconnectDb() {
	m.db.Session.Close()
}

// mongoDB挿入
func (m *mongoDb) InsertDb(obj interface{}, colectionName string) (err error) {
	col := m.db.C(colectionName)
	return col.Insert(obj)
}

// mongoDB削除
func (m *mongoDb) RemoveDb(selector interface{}, colectionName string) (err error) {
	col := m.db.C(colectionName)
	_, err = col.RemoveAll(selector)
	return
}

// mondoDB抽出
func (m *mongoDb) SearchDb(obj, selector interface{}, colectionName string) (err error) {
	col := m.db.C(colectionName)
	return col.Find(selector).All(obj)
}

// mondoDB 1件ランダム抽出
func (m *mongoDb) RandomSearchDb(obj interface{}, colectionName string) (err error) {
	rand.Seed(time.Now().UnixNano())

	col := m.db.C(colectionName)
	colectionCount, err := m.Count(colectionName)
	return col.Find(nil).Skip(rand.Intn(colectionCount)).Limit(1).All(obj)
}

// mongoDB更新
func (m *mongoDb) UpdateDb(selector, update interface{}, colectionName string) (err error) {
	col := m.db.C(colectionName)
	return col.Update(selector, update)
}

// ドキュメント数を返す
func (m *mongoDb) Count(colectionName string) (n int, err error) {
	col := m.db.C(colectionName)
	return col.Count()
}
