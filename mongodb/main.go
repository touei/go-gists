package main

import (
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	session *mgo.Session

	address        string = "127.0.0.1:27017"
	dbName         string = "demodb"
	collectionName string = "democollection"
)

func main() {
	var err error
	session, err = mgo.Dial(address)
	if err != nil {
		beego.Error("链接MongoDB服务失败:", err.Error())
		return
	}

	c := session.DB(dbName).C(collectionName)

	//插入
	err = c.Insert(struct {
		Name           string `bson:"name"`
		LastUpdateTime string `bson:"last_update_time"`
	}{
		Name:           "xdt",
		LastUpdateTime: time.Now().Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		beego.Error("插入失败:", err.Error())
	}

	//查询
	total, _ := c.Find(nil).Count()
	beego.Info("总数:", total)

	var d = make([]struct {
		ID             bson.ObjectId `bson:"_id"`
		Name           string        `bson:"name"`
		LastUpdateTime string        `bson:"last_update_time"`
	}, 0)

	//匹配name=xdt的记录，若传nil查询所有记录
	err = c.Find(bson.M{"name": "xdt"}).All(&d)
	if err != nil {
		beego.Error("查询失败:", err.Error())
		return
	}

	for _, v := range d {
		beego.Info("ID:", v.ID.Hex(), "Name:", v.Name,
			"LastUpdateTime:", v.LastUpdateTime)
	}

	//删除
	if len(d) > 0 {
		err = c.Remove(bson.M{"_id": d[0].ID})
		if err != nil {
			beego.Error("删除失败:", err.Error())
		}
	}
}
