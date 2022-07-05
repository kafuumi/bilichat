package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func connectMysqlDao(user, password, address string, port int, dbname string) (*sql.DB, error) {
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&timeout=1s",
		user, password, address, port, dbname)
	db, err := sql.Open("mysql", sourceName)
	if err != nil {
		return nil, errors.Wrap(err, "sql.open fail")
	}
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "db.ping fail")
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	return db, nil
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
func insertDanMuMsg(db *sql.DB, dms []danMu) error {
	sqlStr := `insert into danmu_msg(room_id, liver_uid, liver_uname, live_status,
                      cmd, time_stamp, medal_level, medal_uid, medal_name,
                      user_uid, user_name, live_level,
                      danmu_text, types, fontsize, color) values`
	values := `(%d, %d, '%s', %t, '%s', %d, %d, %d, '%s', %d, '%s', %d, '%s', %d, %d, %d)`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(dms)
	for i, dm := range dms {
		r := dm.BaseMsg.Room
		m := dm.Medal
		u := dm.User
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			dm.BaseMsg.Cmd, dm.BaseMsg.Timestamp, m.MedalLevel, m.MedalUid, sqlEscape([]byte(m.MedalName)),
			u.UserUid, u.UserName, u.LiveLevel, sqlEscape([]byte(dm.DanMuText)), dm.Types, dm.Fontsize, dm.Color)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}

func insertScMsg(db *sql.DB, scs []sc) error {
	sqlStr := `insert into sc_msg(room_id, liver_uid, liver_uname, live_status,
                   cmd, time_stamp, medal_level, medal_uid, medal_name,
                   user_uid, user_name, live_level, sc_text, price) values`
	values := `(%d, %d, '%s', %t, '%s', %d, %d, %d, '%s', %d, '%s', %d, '%s', %.2f)`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(scs)
	for i, s := range scs {
		r := s.BaseMsg.Room
		m := s.Medal
		u := s.User
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			s.BaseMsg.Cmd, s.BaseMsg.Timestamp, m.MedalLevel, m.MedalUid, sqlEscape([]byte(m.MedalName)),
			u.UserUid, u.UserName, u.LiveLevel, sqlEscape([]byte(s.ScText)), s.Price)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}

func insertGiftMsg(db *sql.DB, gms []gift) error {
	sqlStr := `insert into gift_msg(room_id, liver_uid, liver_uname, live_status,
                     cmd, time_stamp, medal_level, medal_uid, medal_name,
                     user_uid, user_name, gift_id, gift_name, price, num) values`
	values := `(%d, %d, '%s', %t, '%s', %d, %d, %d, '%s', %d, '%s', %d, '%s', %.2f, %d)`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(gms)
	for i, gm := range gms {
		r := gm.BaseMsg.Room
		m := gm.Medal
		u := gm.User
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			gm.BaseMsg.Cmd, gm.BaseMsg.Timestamp, m.MedalLevel, m.MedalUid, sqlEscape([]byte(m.MedalName)),
			u.UserUid, u.UserName, gm.GiftId, gm.GiftName, gm.Price, gm.Num)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}

func insertGuardMsg(db *sql.DB, gms []guard) error {
	sqlStr := `insert into guard_msg(room_id, liver_uid, liver_uname, live_status,
                      cmd, time_stamp, user_uid, user_name, name, price) VALUES`
	values := `(%d, %d, '%s', %t, '%s', %d, %d, '%s', '%s', %.2f)`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(gms)
	for i, gm := range gms {
		r := gm.BaseMsg.Room
		u := gm.User
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			gm.BaseMsg.Cmd, gm.BaseMsg.Timestamp,
			u.UserUid, u.UserName, gm.RoleName, gm.Price)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}

func insertEntryMsg(db *sql.DB, ems []entry) error {
	sqlStr := `insert into entry_msg(room_id, liver_uid, liver_uname, live_status,
                      cmd, time_stamp, user_uid, user_name,
                      medal_level, medal_uid, medal_name) VALUES`
	values := `(%d, %d, '%s', %t, '%s', %d, %d, '%s', %d, %d, '%s')`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(ems)
	for i, em := range ems {
		r := em.BaseMsg.Room
		m := em.Medal
		u := em.User
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			em.BaseMsg.Cmd, em.BaseMsg.Timestamp, u.UserUid, u.UserName,
			m.MedalLevel, m.MedalUid, sqlEscape([]byte(m.MedalName)))
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}

func insertFansMsg(db *sql.DB, rfm []fans) error {
	sqlStr := `insert into fans_msg(room_id, liver_uid, liver_uname, live_status,
                     cmd, time_stamp, fans, fans_club) VALUES`
	values := `(%d, %d, '%s', %t, '%s', %d, %d, %d)`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(rfm)
	for i, rf := range rfm {
		r := rf.BaseMsg.Room
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			rf.BaseMsg.Cmd, rf.BaseMsg.Timestamp, rf.Fans, rf.FansClub)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}

func insertRankCountMsg(db *sql.DB, rcm []rankCount) error {
	sqlStr := `insert into rank_count_msg(room_id, liver_uid, liver_uname, live_status,
                           cmd, time_stamp, count_num) VALUES`
	values := `(%d, %d, '%s', %t, '%s', %d, %d)`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(rcm)
	for i, rc := range rcm {
		r := rc.BaseMsg.Room
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			rc.BaseMsg.Cmd, rc.BaseMsg.Timestamp, rc.CountNum)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}

func insertHotRankMsg(db *sql.DB, hrm []hotRank) error {
	sqlStr := `insert into hot_rank_msg(room_id, liver_uid, liver_uname, live_status, 
				cmd, time_stamp, rank_num, area_name) VALUES`
	values := `(%d, %d, '%s', %t, '%s', %d, %d, '%s')`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(hrm)
	for i, hr := range hrm {
		r := hr.BaseMsg.Room
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			hr.BaseMsg.Cmd, hr.BaseMsg.Timestamp, hr.RankNum, hr.AreaName)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}

func insertRoomChangeMsg(db *sql.DB, rcm []roomChanged) error {
	sqlStr := `insert into room_change_msg(room_id, liver_uid, liver_uname, live_status,
                            cmd, time_stamp, title, area_name, parent_area_name) VALUES`
	values := `(%d, %d, '%s', %t, '%s', %d, '%s', '%s', '%s')`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(rcm)
	for i, rc := range rcm {
		r := rc.BaseMsg.Room
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			rc.BaseMsg.Cmd, rc.BaseMsg.Timestamp, rc.Title, rc.AreaName, rc.ParentAreaName)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}

func insertWatchedChangeMsg(db *sql.DB, wcm []watchedChange) error {
	sqlStr := `insert into watched_change_msg(room_id, liver_uid, liver_uname, live_status,
                               cmd, time_stamp, watched_num) VALUES`
	values := `(%d, %d, '%s', %t, '%s', %d, %d)`
	sb := &strings.Builder{}
	sb.WriteString(sqlStr)
	lens := len(wcm)
	for i, wc := range wcm {
		r := wc.BaseMsg.Room
		_, _ = fmt.Fprintf(sb, values, r.RoomId, r.LiverUid, r.LiverUname, r.LiveStatus,
			wc.BaseMsg.Cmd, wc.BaseMsg.Timestamp, wc.WatchedNum)
		if i != lens-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteByte(';')
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "开启事务失败")
	}
	_, err = tx.Exec(sb.String())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "插入数据失败")
	}
	_ = tx.Commit()
	return nil
}
