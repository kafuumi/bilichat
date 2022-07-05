package main

import (
	"context"
	"database/sql"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v3"
)

type source struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Addr   string `yaml:"addr"`
	Port   int    `yaml:"port"`
	Dbname string `yaml:"dbname"`
}
type config struct {
	Mysql           source `yaml:"mysql"`
	Mongo           source `yaml:"mongo"`
	DanMuId         string `yaml:"danMuId"`
	ScId            string `yaml:"scId"`
	GiftID          string `yaml:"giftID"`
	GuardId         string `yaml:"guardId"`
	EntryId         string `yaml:"entryId"`
	FansId          string `yaml:"fansId"`
	RankCountId     string `yaml:"rankCountId"`
	HotRankId       string `yaml:"hotRankId"`
	RoomChangedId   string `yaml:"roomChangedId"`
	WatchedChangeId string `yaml:"WatchedChangeId"`
}

func main() {
	configFile, err := os.OpenFile("./config.yaml", os.O_RDWR, 0777)
	if err != nil {
		log.Fatalln(err)
	}
	con := readConfig(configFile)
	defer func() {
		_ = configFile.Close()
	}()

	mongoDB, err := connectMongoDB(con.Mongo.User, con.Mongo.Pass,
		con.Mongo.Addr, con.Mongo.Port, con.Mongo.Dbname)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = mongoDB.Client().Disconnect(context.TODO())
	}()

	mysqlDB, err := connectMysqlDao(con.Mysql.User, con.Mysql.Pass,
		con.Mysql.Addr, con.Mysql.Port, con.Mysql.Dbname)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = mysqlDB.Close()
	}()

	sync(&con, mysqlDB, mongoDB)
	_, _ = configFile.Seek(0, 0)
	writeConfig(con, configFile)
}
func sync(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	syncDanMu(con, mysqlDB, mongoDB)
	syncSC(con, mysqlDB, mongoDB)
	syncGift(con, mysqlDB, mongoDB)
	syncGuard(con, mysqlDB, mongoDB)
	syncEntry(con, mysqlDB, mongoDB)
	syncFans(con, mysqlDB, mongoDB)
	syncRankCount(con, mysqlDB, mongoDB)
	syncHotRank(con, mysqlDB, mongoDB)
	syncRoomChange(con, mysqlDB, mongoDB)
	syncWatchedChange(con, mysqlDB, mongoDB)
}

func syncDanMu(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		danMuMsg, err := find[danMu](mongoDB, "danMu", getObjectID(con.DanMuId), 256)
		if err != nil {
			log.Printf("danmu %v\n", err)
			break
		}
		lens := len(danMuMsg)
		total += lens
		if lens == 0 {
			log.Printf("danMu ok, %d\n", total)
			break
		}
		if err = insertDanMuMsg(mysqlDB, danMuMsg); err != nil {
			log.Printf("danmu %v\n", err)
			break
		}
		con.DanMuId = danMuMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func syncSC(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		scMsg, err := find[sc](mongoDB, "sc", getObjectID(con.ScId), 256)
		if err != nil {
			log.Printf("sc %v\n", err)
			break
		}
		lens := len(scMsg)
		total += lens
		if lens == 0 {
			log.Printf("sc ok, %d\n", total)
			break
		}
		if err = insertScMsg(mysqlDB, scMsg); err != nil {
			log.Printf("sc %v\n", err)
			break
		}
		con.ScId = scMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func syncGift(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		giftMsg, err := find[gift](mongoDB, "gift", getObjectID(con.GiftID), 256)
		if err != nil {
			log.Printf("gift %v\n", err)
			break
		}
		lens := len(giftMsg)
		total += lens
		if lens == 0 {
			log.Printf("gift ok, %d\n", total)
			break
		}
		if err = insertGiftMsg(mysqlDB, giftMsg); err != nil {
			log.Printf("gift %v\n", err)
			break
		}
		con.GiftID = giftMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func syncGuard(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		guardMsg, err := find[guard](mongoDB, "guard", getObjectID(con.GuardId), 256)
		if err != nil {
			log.Printf("guard %v\n", err)
			break
		}
		lens := len(guardMsg)
		total += lens
		if lens == 0 {
			log.Printf("guard ok, %d\n", total)
			break
		}
		if err = insertGuardMsg(mysqlDB, guardMsg); err != nil {
			log.Printf("guard %v\n", err)
			break
		}
		con.GuardId = guardMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func syncEntry(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		entryMsg, err := find[entry](mongoDB, "entry", getObjectID(con.EntryId), 256)
		if err != nil {
			log.Printf("entry %v\n", err)
			break
		}
		lens := len(entryMsg)
		total += lens
		if lens == 0 {
			log.Printf("entry ok, %d\n", total)
			break
		}
		if err = insertEntryMsg(mysqlDB, entryMsg); err != nil {
			log.Printf("entry %v\n", err)
			break
		}
		con.EntryId = entryMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func syncFans(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		fansMsg, err := find[fans](mongoDB, "fans", getObjectID(con.FansId), 256)
		if err != nil {
			log.Printf("fans %v\n", err)
			break
		}
		lens := len(fansMsg)
		total += lens
		if lens == 0 {
			log.Printf("fans ok, %d\n", total)
			break
		}
		if err = insertFansMsg(mysqlDB, fansMsg); err != nil {
			log.Printf("fans %v\n", err)
			break
		}
		con.FansId = fansMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func syncRankCount(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		rankCountMsg, err := find[rankCount](mongoDB, "rankCount", getObjectID(con.RankCountId), 256)
		if err != nil {
			log.Printf("rankCount %v\n", err)
			break
		}
		lens := len(rankCountMsg)
		total += lens
		if lens == 0 {
			log.Printf("rankCount ok, %d\n", total)
			break
		}
		if err = insertRankCountMsg(mysqlDB, rankCountMsg); err != nil {
			log.Printf("rankCount %v\n", err)
			break
		}
		con.RankCountId = rankCountMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func syncHotRank(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		hotRankMsg, err := find[hotRank](mongoDB, "hotRank", getObjectID(con.HotRankId), 256)
		if err != nil {
			log.Printf("hotRank %v\n", err)
			break
		}
		lens := len(hotRankMsg)
		total += lens
		if lens == 0 {
			log.Printf("hotRank ok, %d\n", total)
			break
		}
		if err = insertHotRankMsg(mysqlDB, hotRankMsg); err != nil {
			log.Printf("hotRank %v\n", err)
			break
		}
		con.HotRankId = hotRankMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func syncRoomChange(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		roomMsg, err := find[roomChanged](mongoDB, "roomChange", getObjectID(con.RoomChangedId), 256)
		if err != nil {
			log.Printf("roomChanged %v\n", err)
			break
		}
		lens := len(roomMsg)
		total += lens
		if lens == 0 {
			log.Printf("roomChange ok, %d\n", total)
			break
		}
		if err = insertRoomChangeMsg(mysqlDB, roomMsg); err != nil {
			log.Printf("roomChange %v\n", err)
			break
		}
		con.RoomChangedId = roomMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func syncWatchedChange(con *config, mysqlDB *sql.DB, mongoDB *mongo.Database) {
	total := 0
	for {
		watchMsg, err := find[watchedChange](mongoDB, "watchedChange", getObjectID(con.WatchedChangeId), 256)
		if err != nil {
			log.Printf("watchedChange %v\n", err)
			break
		}
		lens := len(watchMsg)
		total += lens
		if lens == 0 {
			log.Printf("watchedChange ok, %d\n", total)
			break
		}
		if err = insertWatchedChangeMsg(mysqlDB, watchMsg); err != nil {
			log.Printf("watchedChange %v\n", err)
			break
		}
		con.WatchedChangeId = watchMsg[lens-1].BaseMsg.ID.Hex()
	}
}

func getObjectID(src string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(src)
	if err != nil {
		return primitive.NilObjectID
	}
	return id
}

func readConfig(r io.Reader) config {
	var c config
	err := yaml.NewDecoder(r).Decode(&c)
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

func writeConfig(con config, w io.Writer) {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	err := encoder.Encode(con)
	if err != nil {
		log.Fatalln(err)
	}
}
