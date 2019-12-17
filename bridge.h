#include <stdio.h>
#include <stdlib.h>
#ifndef POINT_HXX
#define POINT_HXX
typedef int (*ptfFuncAnswerData)(const char* data,int len);
typedef int (*ptfFuncReportData)(const char* data,int len);
typedef int (*ptfFuncMemory)(const char* data,int len , char* redata);
extern int CHandleData(ptfFuncReportData pf,const char* data,int len);
extern int CHandleReData(ptfFuncMemory pf, const char* data,int len , char* redata);
#endif