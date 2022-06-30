create database if not exists live_info;
use live_info;

# 弹幕消息
drop table if exists danmu_msg;
create table danmu_msg
(
    id          int primary key auto_increment, -- 自增长的主键
    room_id     int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid   int,                            -- 主播uid
    liver_uname varchar(64),                    -- 主播昵称
    live_status bool,                           -- 是否开播
    cmd         varchar(64),                    -- websocket消息中的cmd字段
    time_stamp  bigint,                         -- 该消息的时间戳
    medal_level int         default 0,          -- 粉丝牌等级
    medal_uid   bigint      default 0,          -- 粉丝牌对应的账号uid
    medal_name  varchar(64) default '',         -- 粉丝牌名称
    user_uid    bigint,                         -- 该弹幕发送者的uid
    user_name   varchar(64),                    -- 该弹幕发送者的昵称
    live_level  int,                            -- 该弹幕发送者的直播等级
    danmu_text  text,                           -- 弹幕内容
    types       int,                            -- 弹幕类型，1：滚动弹幕，4：底部弹幕，5：顶部弹幕
    fontsize    int         default 25,         -- 弹幕字体大小，一般为25
    color       int                             -- 弹幕颜色，十进制的rgb值
);

# sc 消息
drop table if exists sc_msg;
create table sc_msg
(
    id          int primary key auto_increment, -- 自增长的主键
    room_id     int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid   int,                            -- 主播uid
    liver_uname varchar(64),                    -- 主播昵称
    live_status bool,                           -- 是否开播
    cmd         varchar(64),                    -- websocket消息中的cmd字段
    time_stamp  bigint,                         -- 该消息的时间戳
    medal_level int         default 0,          -- 粉丝牌等级
    medal_uid   bigint      default 0,          -- 粉丝牌对应的账号uid
    medal_name  varchar(64) default '',         -- 粉丝牌名称
    user_uid    bigint,                         -- 该sc发送者的uid
    user_name   varchar(64),                    -- 该sc发送者的昵称
    live_level  int,                            -- 直播等级
    sc_text     text,                           -- sc的内容
    price       int                             -- sc的价格
);

# 礼物消息
drop table if exists gift_msg;
create table gift_msg
(
    id          int primary key auto_increment, -- 自增长的主键
    room_id     int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid   int,                            -- 主播uid
    liver_uname varchar(64),                    -- 主播昵称
    live_status bool,                           -- 是否开播
    cmd         varchar(64),                    -- websocket消息中的cmd字段
    time_stamp  bigint,                         -- 该消息的时间戳
    medal_level int         default 0,          -- 粉丝牌等级
    medal_uid   bigint      default 0,          -- 粉丝牌对应的账号uid
    medal_name  varchar(64) default '',         -- 粉丝牌名称
    user_uid    bigint,                         -- 该礼物发送者的uid
    user_name   varchar(64),                    -- 该礼物发送者的昵称
    gift_id     int,                            -- 礼物id
    gift_name   varchar(64),                    -- 礼物名称
    price       int,                            -- 礼物总价格
    num         int                             -- 礼物数量
);

# 舰长购买消息
drop table if exists guard_msg;
create table guard_msg
(
    id          int primary key auto_increment, -- 自增长的主键
    room_id     int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid   int,                            -- 主播uid
    liver_uname varchar(64),                    -- 主播昵称
    live_status bool,                           -- 是否开播
    cmd         varchar(64),                    -- websocket消息中的cmd字段
    time_stamp  bigint,                         -- 该消息的时间戳
    user_uid    bigint,                         -- uid
    user_name   varchar(64),                    -- 昵称
    name        varchar(64),                    -- 类型：舰长，提督，总督
    price       int                             -- 价格
);

# 进场消息
drop table if exists entry_msg;
create table entry_msg
(
    id          int primary key auto_increment, -- 自增长的主键
    room_id     int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid   int,                            -- 主播uid
    liver_uname varchar(64),                    -- 主播昵称
    live_status bool,                           -- 是否开播
    cmd         varchar(64),                    -- websocket消息中的cmd字段
    time_stamp  bigint,                         -- 该消息的时间戳
    user_uid    bigint,                         -- uid
    user_name   varchar(64),                    -- 昵称
    medal_level int         default 0,
    -- 粉丝牌等级，舰长的进场消息中不含有粉丝牌信息，
    -- 所以如果是舰长进场，等级为21，其他粉丝牌相关字段为默认值
    medal_uid   bigint      default 0,          -- 粉丝牌对应的账号uid
    medal_name  varchar(64) default ''          -- 粉丝牌名称
);

# 粉丝数和粉丝团数量变化消息
drop table if exists fans_msg;
create table fans_msg
(
    id          int primary key auto_increment, -- 自增长的主键
    room_id     int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid   int,                            -- 主播uid
    liver_uname varchar(64),                    -- 主播昵称
    live_status bool,                           -- 是否开播
    cmd         varchar(64),                    -- websocket消息中的cmd字段
    time_stamp  bigint,                         -- 该消息的时间戳
    fans        int,                            -- 变化后的粉丝数
    fans_club   int                             -- 变化后的粉丝团数量
);

# 高能榜人数变化消息
drop table if exists rank_count_msg;
create table rank_count_msg
(
    id          int primary key auto_increment, -- 自增长的主键
    room_id     int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid   int,                            -- 主播uid
    liver_uname varchar(64),                    -- 主播昵称
    live_status bool,                           -- 是否开播
    cmd         varchar(64),                    -- websocket消息中的cmd字段
    time_stamp  bigint,                         -- 该消息的时间戳
    count_num   int                             -- 变化后的数量
);

# 直播间排名变化消息
drop table if exists hot_rank_msg;
create table hot_rank_msg
(
    id          int primary key auto_increment, -- 自增长的主键
    room_id     int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid   int,                            -- 主播uid
    liver_uname varchar(64),                    -- 主播昵称
    live_status bool,                           -- 是否开播
    cmd         varchar(64),                    -- websocket消息中的cmd字段
    time_stamp  bigint,                         -- 该消息的时间戳
    rank_num    int,                            -- 变化后的排名
    area_name   varchar(64)                     -- 所在分区
);

# 直播间信息改变消息
drop table if exists room_change_msg;
create table room_change_msg
(
    id               int primary key auto_increment, -- 自增长的主键
    room_id          int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid        int,                            -- 主播uid
    liver_uname      varchar(64),                    -- 主播昵称
    live_status      bool,                           -- 是否开播
    cmd              varchar(64),                    -- websocket消息中的cmd字段
    time_stamp       bigint,                         -- 该消息的时间戳
    title            varchar(64),                    -- 直播间标题
    area_name        varchar(64),                    -- 直播间分区
    parent_area_name varchar(64)                     -- 直播间父分区
);

# 直播间看过人数变化消息
drop table if exists watched_change_msg;
create table watched_change_msg
(
    id          int primary key auto_increment, -- 自增长的主键
    room_id     int,                            -- 外显的房间号，不一定是真实房间号
    liver_uid   int,                            -- 主播uid
    liver_uname varchar(64),                    -- 主播昵称
    live_status bool,                           -- 是否开播
    cmd         varchar(64),                    -- websocket消息中的cmd字段
    time_stamp  bigint,                         -- 该消息的时间戳
    watched_num int                             -- 变化后的看过人数
)
