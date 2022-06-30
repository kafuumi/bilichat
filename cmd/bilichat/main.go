package main

import (
	"bufio"
	"fmt"
	"github.com/Hami-Lemon/bilichat"
	"os"
)

func main() {
	roomId := 22637261
	c, err := bilichat.GetChatServer(roomId)
	if err != nil {
		panic(err)
	}
	room := c.Room()
	fmt.Printf("title:%s, roomId:%d, uname=%s, uid=%d, isLive=%t\n",
		room.Title, room.Rid, room.Liver.Uname, room.Liver.Uid, room.IsLive)
	err = c.Connect()
	if err != nil {
		panic(err)
	}
	defer c.Disconnect()
	out := make(chan bilichat.Message, 64)
	go c.ReceiveMsg(out)
	f, _ := os.Create("./data_dina.json")
	writer := bufio.NewWriter(f)
	filter := make(map[string]struct{})
	filter["STOP_LIVE_ROOM_LIST"] = struct{}{}
	filter["LIVE_INTERACTIVE_GAME"] = struct{}{}
	filter["DANMU_MSG"] = struct{}{}
	filter["ONLINE_RANK_V2"] = struct{}{}
	filter[bilichat.CmdEntryEffect] = struct{}{}
	filter[bilichat.CmdInteractWord] = struct{}{}
	filter["HOT_ROOM_NOTIFY"] = struct{}{}
	filter["NOTICE_MSG"] = struct{}{}
	filter[bilichat.CmdSendGift] = struct{}{}
	filter["ONLINE_RANK_COUNT"] = struct{}{}
	filter["WATCHED_CHANGE"] = struct{}{}
	filter["HOT_RANK_CHANGED_V2"] = struct{}{}
	filter["WIDGET_BANNER"] = struct{}{}
	filter[bilichat.CmdComboSend] = struct{}{}
	filter["GUARD_BUY"] = struct{}{}
	for {
		select {
		case message := <-out:
			if _, ok := filter[message.MsgType()]; ok {
				break
			}
			_, _ = writer.Write(message.Raw())
			_ = writer.WriteByte('\n')
			_ = writer.Flush()
			switch msg := message.(type) {
			case *bilichat.DanMuMessage:
				//fmt.Printf("[danmu]uname=%s, text=%s, uid=%d, medalUid=%d, medalUname=%s, medalName=%s, medalLevel=%d, level=%d\n",
				//	msg.Uname, msg.Text, msg.Uid, msg.MedalUid, msg.MedalUname, msg.MedalName, msg.MedalLevel, msg.LiveLevel)
			case *bilichat.GiftMessage:
				fmt.Printf("[gift]giftId:%d, giftName:%s, giftNum:%d, price=%d, medalName=%s, medalLevel=%d,  medalUid=%d\n",
					msg.GiftId, msg.GiftName, msg.Num, msg.Price, msg.MedalName, msg.MedalLevel, msg.MedalUid)
			case *bilichat.EntryMessage:
				fmt.Printf("[entry]uname=%s, uid=%d, medalUid=%d,medalName=%s, medalLevel=%d\n",
					msg.Uname, msg.Uid, msg.MedalUid, msg.MedalName, msg.MedalLevel)
			case *bilichat.HotRankMessage:
				fmt.Printf("[rank]rank:%d, area:%s\n", msg.Rank, msg.Area)
			case *bilichat.RankCountMessage:
				fmt.Printf("[count]count%d\n", msg.Count)
			case *bilichat.GuardMessage:
				fmt.Printf("[guard]price=%d, name=%s, uname=%s, uid=%d\n", msg.Price, msg.Name, msg.Uname, msg.Uid)
			case *bilichat.RoomFansMessage:
				fmt.Printf("[fans]fans=%d, fansClud=%d\n", msg.Fans, msg.FansClub)
			case *bilichat.SuperChatMessage:
				fmt.Printf("[SC]price=%d, text=%s\n", msg.Price, msg.Text)
			}
		}
	}
}
