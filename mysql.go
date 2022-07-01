package bilichat

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDao struct {
	db *sql.DB
}

func newMysqlDao(user, password, address string, port int, dbname string) (dao, error) {
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&timeout=1s",
		user, password, address, port, dbname)
	db, err := sql.Open("mysql", sourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(15)
	return &mysqlDao{db: db}, nil
}

var (
	sqlApos      = []byte(`\'`)
	sqlQuot      = []byte(`\"`)
	sqlBackslash = []byte(`\\`)
)

func sqlEscape(src []byte) string {
	var sb strings.Builder
	//判断是否是特殊符号
	isSpecial := func(c byte) bool {
		switch c {
		case '\'', '"', '\\':
			return true
		}
		return false
	}
	last := 0
	for i := 0; i < len(src); i++ {
		c := src[i]
		if !isSpecial(c) {
			continue
		}
		sb.Write(src[last:i])

		switch c {
		case '\'':
			sb.Write(sqlApos)
		case '"':
			sb.Write(sqlQuot)
		case '\\':
			sb.Write(sqlBackslash)
		}
		last = i + 1
	}
	sb.Write(src[last:])
	return sb.String()
}
func (d *mysqlDao) insertDanMuMsg(room Room, dms []*DanMuMessage) error {
	sqlStr := `insert into danmu_msg(room_id, liver_uid, liver_uname, live_status,
                      cmd, time_stamp, medal_level, medal_uid, medal_name,
                      user_uid, user_name, live_level,
                      danmu_text, types, fontsize, color) values`
	values := `(%d, %d, '%s', %t, '%s', %d, %d, %d, '%s', %d, '%s', %d, '%s', %d, %d, %d)`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(dms)
	for i, dm := range dms {
		_, _ = fmt.Fprintf(sb, values, room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
			dm.Cmd, dm.Timestamp, dm.MedalLevel, dm.MedalUid, sqlEscape([]byte(dm.MedalName)),
			dm.Uid, dm.Uname, dm.LiveLevel, sqlEscape([]byte(dm.Text)), dm.Types, dm.FontSize, dm.Color)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		mainLogger.Info(sb.String())
		return err
	}
	_ = tx.Commit()
	return nil
}

func (d *mysqlDao) insertScMsg(room Room, sc *SuperChatMessage) error {
	stmt, err := d.db.Prepare(`insert into sc_msg(room_id, liver_uid, liver_uname, live_status,
                   cmd, time_stamp, medal_level, medal_uid, medal_name,
                   user_uid, user_name, live_level, sc_text, price)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
		sc.Cmd, sc.Timestamp, sc.MedalLevel, sc.MedalUid, sc.MedalName,
		sc.Uid, sc.Uname, sc.LiveLevel, sc.Text, sc.Price)
	if err != nil {
		return err
	}
	return nil
}

func (d *mysqlDao) insertGiftMsg(room Room, gm *GiftMessage) error {
	stmt, err := d.db.Prepare(`insert into gift_msg(room_id, liver_uid, liver_uname, live_status,
                     cmd, time_stamp, medal_level, medal_uid, medal_name,
                     user_uid, user_name, gift_id, gift_name, price, num)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
		gm.Cmd, gm.Timestamp, gm.MedalLevel, gm.MedalUid, gm.MedalName,
		gm.Uid, gm.Uname, gm.GiftId, gm.GiftName, gm.Price, gm.Num)
	if err != nil {
		return err
	}
	return nil
}

func (d *mysqlDao) insertGuardMsg(room Room, gm *GuardMessage) error {
	stmt, err := d.db.Prepare(`insert into guard_msg(room_id, liver_uid, liver_uname, live_status,
                      cmd, time_stamp, user_uid, user_name, name, price)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
		gm.Cmd, gm.Timestamp, gm.Uid, gm.Uname, gm.Name, gm.Price)
	if err != nil {
		return err
	}
	return nil
}

func (d *mysqlDao) insertEntryMsg(room Room, em *EntryMessage) error {
	stmt, err := d.db.Prepare(`insert into entry_msg(room_id, liver_uid, liver_uname, live_status,
                      cmd, time_stamp, user_uid, user_name,
                      medal_level, medal_uid, medal_name)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
		em.Cmd, em.Timestamp, em.Uid, em.Uname,
		em.MedalLevel, em.MedalUid, em.MedalName)
	if err != nil {
		return err
	}
	return nil
}

func (d *mysqlDao) insertFansMsg(room Room, rfm *RoomFansMessage) error {
	stmt, err := d.db.Prepare(`insert into fans_msg(room_id, liver_uid, liver_uname, live_status,
                     cmd, time_stamp, fans, fans_club)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
		rfm.Cmd, rfm.Timestamp, rfm.Fans, rfm.FansClub)
	if err != nil {
		return err
	}
	return nil
}

func (d *mysqlDao) insertRankCountMsg(room Room, rcm *RankCountMessage) error {
	stmt, err := d.db.Prepare(`insert into rank_count_msg(room_id, liver_uid, liver_uname, live_status,
                           cmd, time_stamp, count_num)
VALUES (?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
		rcm.Cmd, rcm.Timestamp, rcm.Count)
	if err != nil {
		return err
	}
	return nil
}

func (d *mysqlDao) insertHotRankMsg(room Room, hrm *HotRankMessage) error {
	stmt, err := d.db.Prepare(`insert into hot_rank_msg(room_id, liver_uid, liver_uname, live_status, 
                         cmd, time_stamp, rank_num, area_name)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
		hrm.Cmd, hrm.Timestamp, hrm.Rank, hrm.Area)
	if err != nil {
		return err
	}
	return nil
}

func (d *mysqlDao) insertRoomChangeMsg(room Room, rcm *RoomChangeMessage) error {
	stmt, err := d.db.Prepare(`insert into room_change_msg(room_id, liver_uid, liver_uname, live_status,
                            cmd, time_stamp, title, area_name, parent_area_name)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
		rcm.Cmd, rcm.Timestamp, rcm.Title, rcm.AreaName, rcm.ParentAreaName)
	if err != nil {
		return err
	}
	return nil
}

func (d *mysqlDao) insertWatchedChangeMsg(room Room, wcm *WatchedChangeMessage) error {
	stmt, err := d.db.Prepare(`insert into watched_change_msg(room_id, liver_uid, liver_uname, live_status,
                               cmd, time_stamp, watched_num)
VALUES (?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room.Id, room.Liver.Uid, room.Liver.Uname, room.IsLive,
		wcm.Cmd, wcm.Timestamp, wcm.Num)
	if err != nil {
		return err
	}
	return nil
}
