package bilichat

import (
	_ "github.com/go-sql-driver/mysql"
)

type dao interface {
	insertDanMuMsg(room Room, dms []*DanMuMessage) error
	insertScMsg(room Room, sc *SuperChatMessage) error
	insertGiftMsg(room Room, gm *GiftMessage) error
	insertGuardMsg(room Room, gm *GuardMessage) error
	insertEntryMsg(room Room, em *EntryMessage) error
	insertFansMsg(room Room, rfm *RoomFansMessage) error
	insertRankCountMsg(room Room, rcm *RankCountMessage) error
	insertHotRankMsg(room Room, hrm *HotRankMessage) error
	insertRoomChangeMsg(room Room, rcm *RoomChangeMessage) error
	insertWatchedChangeMsg(room Room, wcm *WatchedChangeMessage) error
	Close() error
}
