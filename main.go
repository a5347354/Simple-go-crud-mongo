package main

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// 資料庫描述 Descrition of Database
const (
	host   = "localhost:27017"
	source = "<SOURCE>"
	user   = "<USERS>"
	pass   = "<PASSWORD>"
)

// 資料格式 Data Model
type Movies struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
}

func main() {
	dbName := "DBNAME"
	collection := "COLLECTION"

	// 新增一筆資料 Insert Data
	var movie Movies
	movie.ID = bson.NewObjectId()
	movie.Name = "玩命關頭8"
	movie.Description = "以街頭賽車和家人朋友之間的羈絆為主題"
	err1 := Insert(dbName, collection, movie)
	if err1 != nil {
		log.Fatal(err1)
	}

	//尋找資料 Search Data
	var result Movies
	err2 := FindMovieByName(dbName, collection, bson.M{"name": "玩命關頭9"}, nil, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("單筆資料")
	fmt.Println(result)

	//尋找全部電影資料 Search All Data
	var resultAllMovie []Movies
	err3 := FindAllMovies(dbName, collection, nil, nil, &resultAllMovie)
	if err3 != nil {
		log.Fatal(err3)
	}
	fmt.Println("全部資料")
	fmt.Println(resultAllMovie)

	//修改電影資料 Update Data
	//Be careful
	err4 := Update(dbName, collection, bson.M{"name": "玩命關頭9"}, bson.M{"$set": bson.M{"name": "9頭關命玩"}})
	if err4 != nil {
		log.Fatal(err4)
	}

	//刪除電影資料 Remove data
	err5 := Remove(dbName, collection, bson.M{"name": "玩命關頭9"})
	if err5 != nil {
		log.Fatal(err5)
	}
}

var globalS *mgo.Session

// 設定資料庫帳號密碼 Setting of Account and Password
func init() {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{host},
		Source:   source,
		Username: user,
		Password: pass,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalln("create session error ", err)
	}
	globalS = session
}

func connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	session := globalS.Copy()
	c := session.DB(db).C(collection)
	return session, c
}

func Insert(db, collection string, movie Movies) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Insert(movie)
}

func Update(db, collection string, query, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Update(query, update)
}

func Remove(db, collection string, query interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Remove(query)
}

func FindMovieByName(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}

func FindAllMovies(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).All(result)
}
