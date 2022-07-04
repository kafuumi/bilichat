package main

import (
	"bytes"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type info struct {
	Datetime   string
	LiverUname string
	Uid        int64
	Uname      string
	Cmd        string
	Msg        string
	Price      float32
	MedalLevel int
	MedalName  string
}

type Results struct {
	DanMu []*info
	SC    []*info
	Gift  []*info
	Entry []*info
}

func main() {
	tempFile, err := os.Open("./template.html")
	if err != nil {
		log.Fatalln(err)
	}
	sb := &strings.Builder{}
	_, err = io.Copy(sb, tempFile)
	if err != nil {
		log.Fatalln(err)
	}
	_ = tempFile.Close()

	indexFile, err := os.Open("./index.html")
	if err != nil {
		log.Fatalln(err)
	}
	indexBuf := &bytes.Buffer{}
	_, err = io.Copy(indexBuf, indexFile)
	if err != nil {
		log.Fatalln(err)
	}
	_ = indexFile.Close()
	temp := template.Must(template.New("result").Parse(sb.String()))
	sb.Reset()
	db := connectDB()
	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		writer.Header().Set("Content-Language", "zh-CN")
		v := request.URL.Query()
		uid, err := strconv.ParseInt(v.Get("uid"), 10, 64)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
		}
		danMuInfo, err := query(db, uid, "danMu", "danMuText")
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
		}
		giftInfo, err := query(db, uid, "gift", "giftName")
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))

		}
		scInfo, err := query(db, uid, "sc", "scText")
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))

		}
		entryInfo, err := query(db, uid, "entry", "")
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))

		}
		r := Results{
			DanMu: danMuInfo,
			SC:    scInfo,
			Gift:  giftInfo,
			Entry: entryInfo,
		}
		if err = temp.Execute(writer, r); err != nil {
			_, _ = writer.Write([]byte(err.Error()))
		}
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		_, err = io.Copy(writer, indexBuf)
		if err != nil {
			log.Printf("%v\n", err)
		}
	})
	log.Println("开启服务")
	log.Fatalln(http.ListenAndServe("localhost:80", nil))
}

func query(db *mongo.Database, uid int64, col string, msgKey string) ([]*info, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	collection := db.Collection(col)
	cursor, err := collection.Find(ctx, bson.D{
		{"user.userUid", uid},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make([]*info, 0)
	for cursor.Next(ctx) {
		var r bson.D
		err = cursor.Decode(&r)
		if err != nil {
			return nil, err
		}
		rMap := r.Map()
		room := rMap["room"].(bson.D).Map()
		user := rMap["user"].(bson.D).Map()
		medal := rMap["medal"].(bson.D).Map()
		rr := info{
			Datetime:   time.Unix(rMap["timestamp"].(int64), 0).Format("2006-01-02 15:04:05"),
			LiverUname: room["liverUname"].(string),
			Uid:        uid,
			Uname:      user["userName"].(string),
			Cmd:        rMap["cmd"].(string),
			MedalLevel: int(medal["medalLevel"].(int32)),
			MedalName:  medal["medalName"].(string),
		}
		if msgKey == "" {
			rr.Msg = "进场"
		} else {
			rr.Msg = rMap[msgKey].(string)
			if rMap["price"] != nil {
				rr.Price = float32(rMap["price"].(float64))
			}
		}
		result = append(result, &rr)
	}
	return result, nil
}

func connectDB() *mongo.Database {
	ctx := context.Background()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI("mongodb://carol:mongodbcarol@localhost:27017"))
	if err != nil {
		log.Fatalln(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return client.Database("liveInfo")
}
