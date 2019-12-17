package main
// #include "bridge.h"
import "C"
import (
	"clientmsg/call"
	pb "clientmsg/proto"
	"clientmsg/utils"
	"context"
	"errors"
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
var callAnswerBack      C.ptfFuncAnswerData
var callSyncMemoryBack  C.ptfFuncReportData
var callAsyncMemoryBack C.ptfFuncReportData
var callSyncReturnBack  C.ptfFuncMemory
var callAsyncReturnBack C.ptfFuncMemory

//export SetSyncMemoryBack
func SetSyncMemoryBack(f C.ptfFuncReportData) {
	callSyncMemoryBack = f
}
//export SetSyncReturnBack
func SetSyncReturnBack(f C.ptfFuncMemory) {
	callSyncReturnBack = f
}

//export SetAsyncMemoryBack
func SetAsyncMemoryBack(f C.ptfFuncReportData) {
	callAsyncMemoryBack = f
}

//export SetAsyncReturnBack
func SetAsyncReturnBack(f C.ptfFuncMemory) {
	callAsyncReturnBack = f
}

//export SetAnswerBack
func SetAnswerBack(f C.ptfFuncAnswerData) {
	callAnswerBack = f
}

type MsgHandle struct {}

func ApplyMemory(n C.int)[]byte{
	p := C.malloc(C.size_t(n))
	return ((*[1 << 16]byte)(p))[0:n:n]
}

func (m *MsgHandle)Call(ctx context.Context, info *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	////remove uuid
	info.Uuid = ""

	rq,err := proto.Marshal(info)
	if err != nil {
		return &out,err
	}
	/////调用C函数
	length := C.CHandleData(callSyncMemoryBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)))
	/////申请内存空间
	reData := ApplyMemory(length)
	defer C.free(unsafe.Pointer(&reData[0]))

	result := C.CHandleReData(callSyncReturnBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)) , (*C.char)(unsafe.Pointer(&reData[0])))

	if result == 0 {
		return &out,errors.New("call c function handle error")
	}

	resultByte :=  C.GoBytes(unsafe.Pointer(&reData[0]), C.int(result))

	out.M_Net_Rsp = resultByte

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

func (m *MsgHandle)AsyncCall(ctx context.Context, resultInfo *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}

	if resultInfo.Uuid == ""{
		return &out,errors.New("the Async Uuid is empty")
	}
	/// remove Service info
	resultInfo.Service = ""

	rq,err := proto.Marshal(resultInfo)
	if err != nil {
		return &out,err
	}
	/////调用C函数
	length := C.CHandleData(callAsyncMemoryBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)))
	/////申请内存空间
	reData := ApplyMemory(length)
	defer C.free(unsafe.Pointer(&reData[0]))

	result := C.CHandleReData(callAsyncReturnBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)), (*C.char)(unsafe.Pointer(&reData[0])))

	if result == 0 {
		return &out,errors.New("call c function handle error")
	}

	resultByte :=  C.GoBytes(unsafe.Pointer(&reData[0]), C.int(result))

	out.M_Net_Rsp = resultByte
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

func (m *MsgHandle)AsyncAnswer(ctx context.Context, resultInfo *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	resultInfo.Uuid = ""
	resultInfo.Service = ""
	rq,err := proto.Marshal(resultInfo)
	if err != nil {
		return &out,err
	}
	/////调用C函数
	result := C.CHandleData(callAnswerBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)))

	if result == 0 {
		return &out,errors.New("call c function error")
	}

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


