package main
// #include "bridge.h"
// #cgo LDFLAGS: -Wl,-unresolved-symbols=ignore-all
import "C"
import (
	"clientmsg/call"
	pb "clientmsg/proto"
	"clientmsg/utils"
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"unsafe"
)

var (
	Host = utils.Address
	//Port = fmt.Sprintf("%d",utils.Port)
)

var callBackSyncFunc C.ptfFuncReportData
var callBackAsyncFunc C.ptfFuncReportData
//export SetSyncCallBack
func SetSyncCallBack(f C.ptfFuncReportData) {
	callBackSyncFunc = f
}
//export SetAsyncCallBack
func SetAsyncCallBack(f C.ptfFuncReportData) {
	callBackAsyncFunc = f
}

func GoSyncHandleData(data []byte)[]byte {
	result := C.CHandleData(callBackSyncFunc, (*C.char)(unsafe.Pointer(&data[0])), C.int(len(data)))
	defer C.free(unsafe.Pointer(result))
	resultByte := C.GoBytes(unsafe.Pointer(result), C.int(unsafe.Sizeof(result)))
	return resultByte
}
func GoAsyncHandleData(data []byte)[]byte {
	result := C.CHandleData(callBackAsyncFunc, (*C.char)(unsafe.Pointer(&data[0])), C.int(len(data)))
	defer C.free(unsafe.Pointer(result))
	resultByte := C.GoBytes(unsafe.Pointer(result), C.int(unsafe.Sizeof(result)))
	return resultByte
}
type MsgHandle struct {}


func (m *MsgHandle)Call(ctx context.Context, info *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	rq,err := proto.Marshal(info)
	if err != nil {
		return &out,err
	}
	/////调用C函数
	result := GoSyncHandleData(rq)

	err = proto.Unmarshal(result,&out)
	if err != nil {
		return &out,err
	}

	//if HandleObj.Handle == nil {
	//	out.M_Net_Rsp = []byte("The Handle Call function not instance")
	//}else{
	//
	//	//info.Service
	//	//info.M_Body.M_MsgBody.MLBack
	//	//info.M_Body.M_MsgBody.MSSendCount
	//	//info.M_Body.M_MsgBody.MLServerSequence
	//	//info.M_Body.M_MsgBody.MLExpireTime
	//	//info.M_Body.M_MsgBody.MLAskSequence
	//	//info.M_Body.M_MsgBody.MISendTimeApp
	//	//info.M_Body.M_MsgBody.MCMsgType
	//	//info.M_Body.M_MsgBody.MCMsgAckType
	//	//info.M_Body.M_MsgBody.MIDiscard
	//	//info.M_Body.M_MsgBody.MLAsktype
	//	//info.M_Body.M_MsgBody.MLResult
	//
	//	reT,err := HandleObj.Handle(info.M_Body.M_Msg)
	//	if err != nil {
	//		out.M_Net_Rsp = []byte(err.Error())
	//	}
	//	out.M_Net_Rsp = reT
	//}
	return &out,nil
}

func (m *MsgHandle)AsyncCall(ctx context.Context, resultInfo *pb.SingleResultInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	rq,err := proto.Marshal(resultInfo)
	if err != nil {
		return &out,err
	}
	/////调用C函数
	result := GoAsyncHandleData(rq)

	err = proto.Unmarshal(result,&out)
	if err != nil {
		return &out,err
	}
	//if HandleObj.AsyncHandle == nil {
	//	out.M_Net_Rsp = []byte("The AsyncHandle Call function not instance")
	//}else{
	//	reT,err := HandleObj.AsyncHandle(resultInfo)
	//	if err != nil {
	//		out.M_Net_Rsp = []byte(err.Error())
	//	}
	//	out.M_Net_Rsp = reT
	//}
	return &out,nil
}

//export Run
func Run(port []byte)  {
	Port := string(port)
	listener, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		log.Fatalln("faile listen at: " + Host + ":" + Port)
	} else {
		log.Println("server is listening at: " + Host + ":" + Port)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterClientServiceServer(rpcServer,&MsgHandle{})
	reflection.Register(rpcServer)
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatalln("failed serve at: " + Host + ":" + Port)
	}
}

//export Register
func Register(seq,ip,port []byte)[]byte{
	strseq,strip,strport := string(seq),string(ip),string(port)
	err := call.Register(strseq,strip,strport)
	if err != nil {
		return []byte(err.Error())
	}
	return nil
}

//export Publish
func Publish(service []byte)[]byte{
	strservice := string(service)
	err := call.Publish(strservice)
	if err != nil {
		return []byte(err.Error())
	}
	return nil
}

//export Subscribe
func Subscribe(service,ip,port []byte)[]byte{
	strservice,strip,strport := string(service),string(ip),string(port)
	err := call.Subscribe(strservice,strip,strport)
	if err != nil {
		return []byte(err.Error())
	}
	return nil
}

func Broadcast(body,service []byte)[]byte{
	//调用C函数
	service_name := string(service)
	broadResult,err := call.CallBroadcast(body,service_name)
	if err != nil {
		fmt.Println("C++ call Broadcast err:",err.Error())
		return nil
	}
	broadInfo,err := proto.Marshal(broadResult)
	if err != nil {
		fmt.Println("C++ call Broadcast Proto Marshal err:",err.Error())
		return nil
	}
	return broadInfo
}

//export Sync
func Sync(body []byte)[]byte{
	//调用C函数
	syncResult,err := call.CallSync(body)
	if err != nil {
		fmt.Println("C++ call Sync err:",err.Error())
		return nil
	}
	syncInfo,err := proto.Marshal(syncResult)
	if err != nil {
		fmt.Println("C++ call Sync Proto Marshal err:",err.Error())
		return nil
	}
	return syncInfo
}

//export Async
func Async(body []byte)[]byte{
	//调用C函数
	asyncResult,err := call.CallAsync(body)
	if err != nil {
		fmt.Println("C++ call Async err:",err.Error())
		return nil
	}
	asyncInfo,err := proto.Marshal(asyncResult)
	if err != nil {
		fmt.Println("C++ call Async Proto Marshal err:",err.Error())
		return nil
	}
	return asyncInfo
}

func main(){

}


