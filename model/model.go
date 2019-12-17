package model

var (
	HeadTag = "ent2015"
	HeadVersion = 1000
	HeadClientType = Mid_Msg_Type
	HeadLength = 32
	HeadCompressWay = 0
	HeadEncryption = 0
	HeadSig = 0
	HeadFormat = 0
	HeadNetFlag = 0

)

type HeadInfo struct {
	Tag 		string 	//包头标示
	Version 	int16	//数据版本
	ClientType 	int16	//客户端类型
	HeadLength 	int16	//包头长度
	CompressWay uint8	//0表示没有压缩 1 zip
	Encryption  uint8	//加密方   0表示没有加密 1 DES
	Sig			uint8	//协议标记 0 TCP ,1 UDP
	Format		uint8	//数据格式 dataFormat
	NetFlag		uint8	//网络标记
	Back1		uint8	//占位符
	BufSize		int32	//数据长度 ,该数据的长度和MsgBody的m_lBufSize是一样
	UncompressedSize  int32	//压缩前的长度
	Back2		int32	//备份数据

}



type ASK_TYPE int32
var (
	GJ_PUBLIC_START        		ASK_TYPE 	=   0							//公共部分的请求
	GJ_PUBLIC_NET_OPERATION		ASK_TYPE 	=   GJ_PUBLIC_START+10000		//操作部分
	ETN_ASK_LOAIN_SERVER     	ASK_TYPE 	=	GJ_PUBLIC_NET_OPERATION+1 	// 登录
	ETN_SERVER_NET_CONNET    	ASK_TYPE 	=	GJ_PUBLIC_NET_OPERATION+2 	//有客户端连接成功
	ETN_SERVER_NET_CLOSE     	ASK_TYPE 	=	GJ_PUBLIC_NET_OPERATION+3	//服务端网络层关闭
	ETN_ASK_USER_LEAVE       	ASK_TYPE 	=	GJ_PUBLIC_NET_OPERATION+4	//用户登录退出
	ETN_SERVER_PUSH_NOTICE_MSG  ASK_TYPE 	=   GJ_PUBLIC_NET_OPERATION+5	//服务器推送通知
	ETN_HEARTBEAT_PACK          ASK_TYPE 	=   GJ_PUBLIC_NET_OPERATION+6	//心跳包
	ETN_SERVER_SUBSRCTIBE_MSG   ASK_TYPE 	=   GJ_PUBLIC_NET_OPERATION+7	//广播消息
)

type MSG_TYPE int32
var (
	MSG_TYPE_ACK				MSG_TYPE	= 0				//普通请求类型
	MSG_TYPE_LOGIN_REQ			MSG_TYPE	= 1				//注册
	MSG_TYPE_LOGIN_ACK			MSG_TYPE	= 2				//注册响应
	MSG_TYPE_KEEPALIVE_REQ		MSG_TYPE	= 3				//心跳检测
	MSG_TYPE_KEEPALIVE_ACK		MSG_TYPE	= 4				//心跳检测响应
	MSG_TYPE_PUSHMSG_REQ		MSG_TYPE	= 5				//下发消息
	MSG_TYPE_PUSHMSG_ACK		MSG_TYPE	= 6				//下发消息响应
	MSG_TYPE_UPLOADMSG_REQ		MSG_TYPE	= 7				//上传消息
	MSG_TYPE_UPLOADMSG_ACK		MSG_TYPE	= 8				//上传消息响应
	MSG_TYPE_BROADCAST			MSG_TYPE	= 9				//广播消息
	MSG_TYPE_SUBSCRIBE_REQ		MSG_TYPE	= 10			//订阅消息
	MSG_TYPE_SUBSCRIBE_ACK		MSG_TYPE	= 11			//订阅消息响应
	MSG_TYPE_ERROR				MSG_TYPE	= 12			//错误信息应答
	MSG_TYPE_NOTICE             MSG_TYPE    = 13			//通知消息
	MSG_TYPEMAX					MSG_TYPE    = 14			//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
	MSG_TYPE_CONNECT_COUNT_MSG	MSG_TYPE	= 0xFFFF		//查询连接数
)

type CLIENT_TYPE int16
var (
	Window_pc				CLIENT_TYPE		= 0				//请求来自PC
	IOS_mobile				CLIENT_TYPE		= 1				//请求来自苹果手机
	Android_mobilewindow_pc	CLIENT_TYPE		= 3				//请求来自安卓手机
	Web_side				CLIENT_TYPE		= 4				//请求来自WEB端
	Mid_Msg_Type			CLIENT_TYPE		= 5				//中间件封装
	ClientTypeMax			CLIENT_TYPE		= 6				//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
)
type COMPRESS_TYPE int16
var (
	Compression_no			COMPRESS_TYPE		= 0				//表示数据未压缩
	Compression_zip			COMPRESS_TYPE		= 1				//表示数据被压缩zlib压缩格式
	CompressionWayMax		COMPRESS_TYPE		= 2				//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
)

type ENCRPTION_TYPE int16
var (
	Encryption_No			ENCRPTION_TYPE		= 0				//请求数据流未被加密
	Encryption_Des			ENCRPTION_TYPE		= 1				//请求数据是通过Des算法加密
	Encryption_AES			ENCRPTION_TYPE		= 2				//请求数据是通过AES算法加密
	Encryption_RSA			ENCRPTION_TYPE		= 3				//请求数据是通过RSA算法加密
	Encryption_Max			ENCRPTION_TYPE		= 4				//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
)

type DATAFORMAT_TYPE int16

var (
	DataFormat_Probufo		DATAFORMAT_TYPE		= 0			//数据是通过protobuf进行格式化
	DataFormatMax			DATAFORMAT_TYPE		= 1			//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
)

type CALL_CLIENT_TYPE uint32
var (
	CALL_CLIENT_SYNC 	  CALL_CLIENT_TYPE = 0
	CALL_CLIENT_ASYNC     CALL_CLIENT_TYPE = 1
	CALL_CLIENT_PUBLISH   CALL_CLIENT_TYPE = 2
	CALL_CLIENT_ANSWER    CALL_CLIENT_TYPE = 3
)