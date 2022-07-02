package bilichat

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDao struct {
	db            *mongo.Database
	danMu         *mongo.Collection
	sc            *mongo.Collection
	gift          *mongo.Collection
	guard         *mongo.Collection
	entry         *mongo.Collection
	fans          *mongo.Collection
	rankCount     *mongo.Collection
	hotRank       *mongo.Collection
	roomChange    *mongo.Collection
	watchedChange *mongo.Collection
}

func newMongoDao(user, password, address string, port int, dbname string) (dao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var u string
	if user == "" && password == "" {
		u = fmt.Sprintf("mongodb://%s:%d", address, port)
	} else {
		u = fmt.Sprintf("mongodb://%s:%s@%s:%d", user, password, address, port)
	}
	opt := options.Client().ApplyURI(u)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, errors.Wrap(err, "mongo connect fail")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "mongo ping fail")
	}
	db := client.Database(dbname)
	return &mongoDao{
		db:            db,
		danMu:         db.Collection("danMu"),
		sc:            db.Collection("sc"),
		gift:          db.Collection("gift"),
		guard:         db.Collection("guard"),
		entry:         db.Collection("entry"),
		fans:          db.Collection("fans"),
		rankCount:     db.Collection("rankCount"),
		hotRank:       db.Collection("hotRank"),
		roomChange:    db.Collection("roomChange"),
		watchedChange: db.Collection("watchedChange"),
	}, nil
}

func (m *mongoDao) insertDanMuMsg(room Room, dms []*DanMuMessage) error {
	coll := m.danMu
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	docs := make([]interface{}, 0)
	for _, dm := range dms {
		doc := bson.D{
			{"cmd", dm.Cmd},
			{"timestamp", dm.Timestamp},
			{"room", bson.D{
				{"roomId", room.Id},
				{"liverUid", room.Liver.Uid},
				{"liverUname", room.Liver.Uname},
				{"liveStatus", room.IsLive},
			}},
			{"medal", bson.D{
				{"medalLevel", dm.MedalLevel},
				{"medalUid", dm.MedalUid},
				{"medalName", dm.MedalName},
			}},
			{"user", bson.D{
				{"userUid", dm.Uid},
				{"userName", dm.Uname},
				{"liveLevel", dm.LiveLevel},
			}},
			{"danMuText", dm.Text},
			{"types", dm.Types},
			{"fontsize", dm.FontSize},
			{"color", dm.Color},
		}
		docs = append(docs, doc)
	}
	_, err := coll.InsertMany(ctx, docs)
	return err
}

func (m *mongoDao) insertScMsg(room Room, sc *SuperChatMessage) error {
	coll := m.sc
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	doc := bson.D{
		{"cmd", sc.Cmd},
		{"timestamp", sc.Timestamp},
		{"room", bson.D{
			{"roomId", room.Id},
			{"liverUid", room.Liver.Uid},
			{"liverUname", room.Liver.Uname},
			{"liveStatus", room.IsLive},
		}},
		{"medal", bson.D{
			{"medalLevel", sc.MedalLevel},
			{"medalUid", sc.MedalUid},
			{"medalName", sc.MedalName},
		}},
		{"user", bson.D{
			{"userUid", sc.Uid},
			{"userName", sc.Uname},
			{"liveLevel", sc.LiveLevel},
		}},
		{"scText", sc.Text},
		{"price", sc.Price},
	}
	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (m *mongoDao) insertGiftMsg(room Room, gm *GiftMessage) error {
	coll := m.gift
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	doc := bson.D{
		{"cmd", gm.Cmd},
		{"timestamp", gm.Timestamp},
		{"room", bson.D{
			{"roomId", room.Id},
			{"liverUid", room.Liver.Uid},
			{"liverUname", room.Liver.Uname},
			{"liveStatus", room.IsLive},
		}},
		{"medal", bson.D{
			{"medalLevel", gm.MedalLevel},
			{"medalUid", gm.MedalUid},
			{"medalName", gm.MedalName},
		}},
		{"user", bson.D{
			{"userUid", gm.Uid},
			{"userName", gm.Uname},
		}},
		{"giftId", gm.GiftId},
		{"giftName", gm.GiftName},
		{"price", gm.Price},
		{"num", gm.Num},
	}
	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (m *mongoDao) insertGuardMsg(room Room, gm *GuardMessage) error {
	coll := m.guard
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	doc := bson.D{
		{"cmd", gm.Cmd},
		{"timestamp", gm.Timestamp},
		{"room", bson.D{
			{"roomId", room.Id},
			{"liverUid", room.Liver.Uid},
			{"liverUname", room.Liver.Uname},
			{"liveStatus", room.IsLive},
		}},
		{"user", bson.D{
			{"userUid", gm.Uid},
			{"userName", gm.Uname},
		}},
		{"roleName", gm.Name},
		{"price", gm.Price},
	}
	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (m *mongoDao) insertEntryMsg(room Room, em *EntryMessage) error {
	coll := m.entry
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	doc := bson.D{
		{"cmd", em.Cmd},
		{"timestamp", em.Timestamp},
		{"room", bson.D{
			{"roomId", room.Id},
			{"liverUid", room.Liver.Uid},
			{"liverUname", room.Liver.Uname},
			{"liveStatus", room.IsLive},
		}},
		{"user", bson.D{
			{"userUid", em.Uid},
			{"userName", em.Uname},
		}},
		{"medal", bson.D{
			{"medalLevel", em.MedalLevel},
			{"medalUid", em.MedalUid},
			{"medalName", em.MedalName},
		}},
	}
	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (m *mongoDao) insertFansMsg(room Room, rfm *RoomFansMessage) error {
	coll := m.fans
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	doc := bson.D{
		{"cmd", rfm.Cmd},
		{"timestamp", rfm.Timestamp},
		{"room", bson.D{
			{"roomId", room.Id},
			{"liverUid", room.Liver.Uid},
			{"liverUname", room.Liver.Uname},
			{"liveStatus", room.IsLive},
		}},
		{"fans", rfm.Fans},
		{"fansClub", rfm.FansClub},
	}
	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (m *mongoDao) insertRankCountMsg(room Room, rcm *RankCountMessage) error {
	coll := m.rankCount
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	doc := bson.D{
		{"cmd", rcm.Cmd},
		{"timestamp", rcm.Timestamp},
		{"room", bson.D{
			{"roomId", room.Id},
			{"liverUid", room.Liver.Uid},
			{"liverUname", room.Liver.Uname},
			{"liveStatus", room.IsLive},
		}},
		{"countNum", rcm.Count},
	}
	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (m *mongoDao) insertHotRankMsg(room Room, hrm *HotRankMessage) error {
	coll := m.hotRank
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	doc := bson.D{
		{"cmd", hrm.Cmd},
		{"timestamp", hrm.Timestamp},
		{"room", bson.D{
			{"roomId", room.Id},
			{"liverUid", room.Liver.Uid},
			{"liverUname", room.Liver.Uname},
			{"liveStatus", room.IsLive},
		}},
		{"rankNum", hrm.Rank},
		{"areaNum", hrm.Area},
	}
	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (m *mongoDao) insertRoomChangeMsg(room Room, rcm *RoomChangeMessage) error {
	coll := m.roomChange
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	doc := bson.D{
		{"cmd", rcm.Cmd},
		{"timestamp", rcm.Timestamp},
		{"room", bson.D{
			{"roomId", room.Id},
			{"liverUid", room.Liver.Uid},
			{"liverUname", room.Liver.Uname},
			{"liveStatus", room.IsLive},
		}},
		{"title", rcm.Title},
		{"areaName", rcm.AreaName},
		{"parentAreaName", rcm.ParentAreaName},
	}
	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (m *mongoDao) insertWatchedChangeMsg(room Room, wcm *WatchedChangeMessage) error {
	coll := m.watchedChange
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	doc := bson.D{
		{"cmd", wcm.Cmd},
		{"timestamp", wcm.Timestamp},
		{"room", bson.D{
			{"roomId", room.Id},
			{"liverUid", room.Liver.Uid},
			{"liverUname", room.Liver.Uname},
			{"liveStatus", room.IsLive},
		}},
		{"watchedNum", wcm.Num},
	}
	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (m *mongoDao) Close() error {
	client := m.db.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}
