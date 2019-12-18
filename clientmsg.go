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
var callSyncBack      C.ptfFuncCallBack
var callAsyncBack     C.ptfFuncCallBack
var callAnswerBack    C.ptfFuncCall
//export SetSyncReturnBack
func SetSyncReturnBack(f C.ptfFuncCallBack) {
	callSyncBack = f
}

//export SetAsyncReturnBack
func SetAsyncReturnBack(f C.ptfFuncCallBack) {
	callAsyncBack = f
}

//export SetAnswerBack
func SetAnswerBack(f C.ptfFuncCall) {
	callAnswerBack = f
}

type MsgHandle struct {}

//export GenerateMemory
func GenerateMemory(n int)unsafe.Pointer{
	p := C.malloc(C.sizeof_char * C.uint(n))
	return unsafe.Pointer(p)
}

func (m *MsgHandle)Call(ctx context.Context, info *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	////remove uuid
	info.Uuid = ""
	////获取msg head


	////获取msg body
	rq := info.M_Body.M_Msg
	/////调用C函数
	p := C.CHandleData(callSyncBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)))
	defer C.free(unsafe.Pointer(p.content))
	resultByte := C.GoBytes(unsafe.Pointer(p.content), p.length)
	out.M_Net_Rsp = resultByte
	return &out,nil
}

func (m *MsgHandle)AsyncCall(ctx context.Context, resultInfo *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}

	if resultInfo.Uuid == ""{
		return &out,errors.New("the Async Uuid is empty")
	}
	/// remove Service info
	resultInfo.Service = ""
	//resultInfo.Uuid
	////获取msg head

	////获取msg body
	rq := resultInfo.M_Body.M_Msg
	/////调用C函数
	p := C.CHandleData(callAsyncBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)))
	defer C.free(unsafe.Pointer(p.content))
	resultByte := C.GoBytes(unsafe.Pointer(p.content), p.length)
	out.M_Net_Rsp = resultByte
	return &out,nil
}

func (m *MsgHandle)AsyncAnswer(ctx context.Context, resultInfo *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	resultInfo.Uuid = ""
	resultInfo.Service = ""

	////获取msg body
	rq := resultInfo.M_Body.M_Msg
	/////调用C函数
	result := C.CHandleCall(callAnswerBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)))

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

func Broadcast(body,service []byte,info C.BodyInfo)[]byte{
	infoBody , err := MarshalBody(body,info)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	//调用C函数
	service_name := string(service)

	broadResult,err := call.CallBroadcast(infoBody,service_name)
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
func Sync(body []byte,info C.BodyInfo)[]byte{
	infoBody , err := MarshalBody(body,info)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	//调用C函数
	syncResult,err := call.CallSync(infoBody)
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
func Async(body []byte,info C.BodyInfo)[]byte{
	infoBody , err := MarshalBody(body,info)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	//调用C函数
	asyncResult,err := call.CallAsync(infoBody)
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


