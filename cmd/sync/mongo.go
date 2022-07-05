package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type msg interface {
	danMu | sc | gift | guard | entry | fans | rankCount | hotRank | roomChanged | watchedChange
}

type room struct {
	RoomId     int    `bson:"roomId"`
	LiverUid   int64  `bson:"liverUid"`
	LiverUname string `bson:"liverUname"`
	LiveStatus bool   `bson:"liveStatus"`
}
type baseMsg struct {
	ID        primitive.ObjectID `bson:"_id"`
	Cmd       string             `bson:"cmd"`
	Timestamp int64              `bson:"timestamp"`
	Room      room               `bson:"room"`
}
type medal struct {
	MedalLevel int    `bson:"medalLevel"`
	MedalUid   int64  `bson:"medalUid"`
	MedalName  string `bson:"medalName"`
}

type user struct {
	UserUid   int64  `bson:"userUid"`
	UserName  string `bson:"userName"`
	LiveLevel int    `bson:"liveLevel,omitempty"`
}
type danMu struct {
	BaseMsg   baseMsg `bson:"inline"`
	Medal     medal   `bson:"medal"`
	User      user    `bson:"user"`
	DanMuText string  `bson:"danMuText"`
	Types     int     `bson:"types"`
	Fontsize  int     `bson:"fontsize"`
	Color     int     `bson:"color"`
}

type sc struct {
	BaseMsg baseMsg `bson:"inline"`
	Medal   medal   `bson:"medal"`
	User    user    `bson:"user"`
	ScText  string  `bson:"scText"`
	Price   float32 `bson:"price"`
}

type gift struct {
	BaseMsg  baseMsg `bson:"inline"`
	Medal    medal   `bson:"medal"`
	User     user    `bson:"user"`
	GiftId   int     `bson:"giftId"`
	GiftName string  `bson:"giftName"`
	Price    float32 `bson:"price"`
	Num      int     `bson:"num"`
}

type guard struct {
	BaseMsg  baseMsg `bson:"inline"`
	User     user    `bson:"user"`
	RoleName string  `bson:"roleName"`
	Price    float32 `bson:"price"`
}

type entry struct {
	BaseMsg baseMsg `bson:"inline"`
	User    user    `bson:"user"`
	Medal   medal   `bson:"medal"`
}

type fans struct {
	BaseMsg  baseMsg `bson:"inline"`
	Fans     int     `bson:"fans"`
	FansClub int     `bson:"fansClub"`
}

type rankCount struct {
	BaseMsg  baseMsg `bson:"inline"`
	CountNum int     `bson:"countNum"`
}

type hotRank struct {
	BaseMsg  baseMsg `bson:"inline"`
	RankNum  int     `bson:"rankNum"`
	AreaName string  `bson:"areaNum"` //铸币了，数据库中这个字段的名称整错了，懒得改了，就这样
}

type roomChanged struct {
	BaseMsg        baseMsg `bson:"inline"`
	Title          string  `bson:"title"`
	AreaName       string  `bson:"areaName"`
	ParentAreaName string  `bson:"parentAreaName"`
}

type watchedChange struct {
	BaseMsg    baseMsg `bson:"inline"`
	WatchedNum int     `bson:"watchedNum"`
}

func connectMongoDB(user, pass, addr string, port int, dbname string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s:%d", user, pass, addr, port)))
	if err != nil {
		return nil, errors.Wrap(err, "创建mongodb客户端失败")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "连接mongodb失败")
	}
	return client.Database(dbname), nil
}

func find[T msg](db *mongo.Database, colName string, idStart primitive.ObjectID, limit int) ([]T, error) {
	col := db.Collection(colName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var filter = bson.D{}
	if !idStart.IsZero() {
		filter = bson.D{
			{"_id", bson.D{
				{"$gt", idStart},
			}},
		}
	}
	cursor, err := col.Find(ctx, filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, errors.Wrap(err, "获取数据失败")
	}
	var result []T
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, errors.Wrap(err, "解析数据失败")
	}
	return result, nil
}
