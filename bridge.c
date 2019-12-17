#include "bridge.h"

char* CHandleData(ptfFuncCallBack pf,const char* data,int len,int* length){
    return pf(data,len,length);
}

int CHandleCall(ptfFuncCall pf,const char* data,int len){
    return pf(data,len);
}


struct BodyInfo InitializeBody(){

    struct BodyInfo structInitialized;

    structInitialized.ServerSequence = 0;
    structInitialized.AskSequence = 0;
    structInitialized.MsgAckType = 1;
    structInitialized.MsgType   = 0;
    structInitialized.SendCount = 1;
    structInitialized.ExpireTime = 0;
    structInitialized.SendTimeApp = 0;
    structInitialized.Result = 0;
    structInitialized.Back = 0;
    structInitialized.Discard = 0;

    return(structInitialized);
};