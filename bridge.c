#include "bridge.h"

int CHandleData(ptfFuncReportData pf,const char* data,int len){
    return pf(data,len);
}

int CHandleReData(ptfFuncMemory pf, char* redata){
    return pf(redata);
}