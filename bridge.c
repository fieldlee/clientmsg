#include "bridge.h"

ReturnInfo CSyncHandleData(ptfSyncCallBack pf, CStr data,  uint64_t service, CStr ip){
    return pf(data,service,ip);
}

ReturnInfo CAsyncHandleData(ptfAsyncCallBack pf, CStr data,  CStr uid,  uint64_t service, CStr ip){
    return pf(data,uid,service,ip);
}

int CHandleCall(ptfFuncCall pf, CStr data,int success){
    return pf(data,success);
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
    return(structInitialized);
};