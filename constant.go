package bilichat

//数据类型
const (
	verPlain  = 0 //普通文本，utf-8编码
	verInt    = 1 //控制通信消息
	verZlib   = 2 // zlib格式压缩，已弃用
	verBrotli = 3 //brotli压缩
)

//操作码
const (
	opHeartbeat      = 2 //发送心跳包
	opHeartbeatReply = 3 //服务端回应心跳包
	opMessage        = 5 //弹幕消息等
	opEnterRoom      = 7 //进入直播间
	opEnterRoomReply = 8 //进入直播间成功
)

// Cmd
const (
	CmdInteractWord              = "INTERACT_WORD"                 //进场消息
	CmdEntryEffect               = "ENTRY_EFFECT"                  //舰长进场消息
	CmdSendGift                  = "SEND_GIFT"                     //投喂礼物
	CmdComboSend                 = "COMBO_SEND"                    //礼物连击
	CmdDanMuMSG                  = "DANMU_MSG"                     //弹幕
	CmdWatchedChange             = "WATCHED_CHANGE"                //看过人数变化
	CmdOnlineRankCount           = "ONLINE_RANK_COUNT"             //高能榜人数
	CmdUserToastMsg              = "USER_TOAST_MSG"                //续费舰长
	CmdSuperChatMessage          = "SUPER_CHAT_MESSAGE"            //sc
	CmdRoomRealTimeMessageUpdate = "ROOM_REAL_TIME_MESSAGE_UPDATE" //粉丝数，粉丝团变化
	CmdLive                      = "LIVE"                          //开播了
	CmdPreparing                 = "PREPARING"                     //下播了
	CmdRoomChange                = "ROOM_CHANGE"                   //直播间信息变化
	CmdRoomBlackMsg              = "ROOM_BLACK_MSG"                //用户被禁言
	CmdCutOff                    = "CUT_OFF"                       //被超管切断
	CmdHotRankChanged            = "HOT_RANK_CHANGED_V2"           //直播间分区排名变化
)
