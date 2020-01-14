#include "bridge.h"

ReturnInfo CHandleData(ptfFuncCallBack pf,const char* data,int len, char* uid,int uidlen,uint64_t service,char* clientip, int iplen){
    return pf(data,len,uid,uidlen,service,clientip,iplen);
}

int CHandleCall(ptfFuncCall pf,const char* data,int len){
    return pf(data,len);
}

BodyInfo InitializeBody(){
    BodyInfo structInitialized;
    structInitialized.ServerSequence = 0;
    structInitialized.AskSequence = 0;
    structInitialized.MsgAckType  = 1;
    structInitialized.MsgType     = 0;
    structInitialized.SendCount   = 1;
    structInitialized.ExpireTime  = 0;
    structInitialized.SendTimeApp = 0;
    structInitialized.Result   = 0;
    structInitialized.Back     = 0;
    structInitialized.Discard  = 0;
    structInitialized.Encrypt  = 0;
    structInitialized.Compress = 0;
    structInitialized.UUID     = 0;
    structInitialized.LenID    = 0;
    return(structInitialized);
};