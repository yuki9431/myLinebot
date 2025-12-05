mongohelper
====
## Overview
- Go言語用のmongoDBドライバ [mgo](https://github.com/go-mgo/mgo)を自分が使いやすいようにカスタマイズ

## Descriptionz
挿入、抽出、削除を関数一つで実行できる  

## Requirement
- Go 1.10 or later

## Install
```bash:#
go get github.com/yuki9431/mongohelper
```

## Configuration
```go:main.go
import (
	"fmt"

	"github.com/yuki9431/mongohelper"
)

func main() {
	mongoDial         = "mongodb://localhost/mongodb"
	mongoName         = "databaseName"
	
	// DB設定
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	if err != nil {
		fmt.Println(err)
	}
}
```

### NewMongo(dial string, name string) 
DBの接続を開始し、MongoDb型を返します。  

### DisconnectDb()
DBを切断します。  

### InsertDb(obj interface{}, colectionName string)
objをcolectionNameに格納します。  

### RemoveDb(selector interface{}, colectionName string)
selectorに合致するドキュメントを削除する。  
一意の値を渡すと、1つのドキュメントを削除するが、複数に一致する場合は、一致するドキュメントを全て削除する。

```go:main.go
// example, remove document
userId := Hatsune39
selector := bson.M{"userid": userId}

if err := mongo.RemoveDb(selector, "userInfos"); err != nil {
	fmt.Println(err)
}
```

### SearchDb(obj, selector interface{}, colectionName string)
selectorに合致する全てのドキュメントをobjに格納する。

### UpdateDb(selector, update interface{}, colectionName string)
selectorに合致するドキュメントを全て更新する。  
updateに更新したいデータを渡す。  

```
// example, update name Hatsune39
userId := Hatsune39
newName := "Megurine Luka"
selector := bson.M{"userid": userId}
update := bson.M{"$set": bson.M{"name": newName}}

if err := mongo.UpdateDb(selector, update, "userInfos"); err != nil {
	fmt.Println(err)
}
```


## Contribution
1. Fork ([https://github.com/yuki9431/mongohelper](https://github.com/yuki9431/mongohelper))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create new Pull Request


## Author
[Dillen H. Tomida](https://twitter.com/t0mihir0)
