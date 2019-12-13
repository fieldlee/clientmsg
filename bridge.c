#include "bridge.h"

char* CHandleData(ptfFuncReportData pf,const char* data,int len){
    return pf(data,len);
}
