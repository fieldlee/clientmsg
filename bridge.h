#ifndef POINT_HXX
#define POINT_HXX
typedef char* (*ptfFuncReportData)(const char* data,int len);
extern char* CHandleData(ptfFuncReportData pf,const char* data,int len);
#endif