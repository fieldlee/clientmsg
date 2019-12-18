#include <stdio.h>
#include <stdlib.h>
#ifndef POINT_HXX
#define POINT_HXX

typedef struct _ReturnInfo{
    char* content;
    int   length;
}ReturnInfo;

typedef struct _BodyInfo {
     long unsigned int  Asktype;         //请求的服务类型
     long unsigned int  ServerSequence;  //服务端响应序列号
     long unsigned int  AskSequence;     //客户端请求序列号
     int                MsgAckType;      //0 无需回复, 1 回复到发送方, 2 回复到离线服务器
     int                MsgType;         //消息类型
     int                SendCount;       //同一请求，请求服务端的次数
     unsigned int       ExpireTime;      //过期时间，0表示永不过期
     long unsigned int  SendTimeApp;     //请求发送的时间，单位秒
     unsigned int       Result;          //0: SUCCESS  !0:FAILURE
     long unsigned int  Back;            //备份数据，默认为0
     int                Discard;         //消息是否可以丢弃 0 表示可
}BodyInfo;

extern _BodyInfo InitializeBody(); // 初始化BodyInfo

typedef _ReturnInfo (*ptfFuncCallBack)(const char* data,int len);
typedef int (*ptfFuncCall)(const char* data,int len);
extern _ReturnInfo CHandleData(ptfFuncCallBack pf,const char* data,int len);
extern int CHandleCall(ptfFuncCall pf,const char* data,int len);
#endif