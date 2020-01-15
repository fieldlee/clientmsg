#include "bridge.h"

ReturnInfo CSyncHandleData(ptfSyncCallBack pf,const CStr data, const uint64_t service,const CStr ip){
    return pf(data,service,ip);
}

ReturnInfo CAsyncHandleData(ptfAsyncCallBack pf,const CStr data, const CStr uid, const uint64_t service,const CStr ip){
    return pf(data,uid,service,ip);
}

int CHandleCall(ptfFuncCall pf,const CStr data){
    return pf(data);
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