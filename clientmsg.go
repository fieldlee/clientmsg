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
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

var (
	Host = utils.Address
	Port = fmt.Sprintf("%d",utils.Port)
)
var callSyncBack      C.ptfFuncCallBack
var callAsyncBack     C.ptfFuncCallBack
var callAnswerBack    C.ptfFuncCall

type RInfo C.CallReturnInfo
type MsgInfo C.MsgReturnInfo

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

func (m *MsgHandle)Call(ctx context.Context, info *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	////remove uuid
	uid := []byte("")
	////获取msg body
	rq := info.M_Body.M_Msg
	/////调用C函数
	p := C.CHandleData(callSyncBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)),(*C.char)(unsafe.Pointer(&uid[0])),C.int(0))
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
	//// uuid bytes
	uid := []byte(resultInfo.Uuid)
	////获取msg body
	rq := resultInfo.M_Body.M_Msg
	/////调用C函数
	p := C.CHandleData(callAsyncBack, (*C.char)(unsafe.Pointer(&rq[0])), C.int(len(rq)),(*C.char)(unsafe.Pointer(&uid[0])),C.int(len(uid)))
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
func Run()  {
	listener, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		log.Fatalln("faile listen at: " + Host + ":" + Port)
	} else {
		log.Println("server is listening at: " + Host + ":" + Port)
	}
	rpcServer := grpc.NewServer()

	go func() {
		log.Println("go routine waiting shutdown...")
		waitForShutdown(rpcServer)
	}()

	pb.RegisterClientServiceServer(rpcServer,&MsgHandle{})
	reflection.Register(rpcServer)
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatalln("failed serve at: " + Host + ":" + Port)
	}
}

func waitForShutdown(srv *grpc.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Block until we receive our signal.
	<-interruptChan
	// graceful shutdown grpc service
	srv.GracefulStop()

	log.Println("grpc stop...")
	// shutdown service
	os.Exit(0)
}
//export Register
func Register(seq []byte)RInfo{
	r := RInfo{}
	r.success = C.int(0)
	strseq := string(seq)
	err := call.Register(strseq)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
	}
	return r
}


//export Publish
func Publish(service []byte)RInfo{
	r := RInfo{}
	r.success = C.int(0)
	strservice := string(service)
	err := call.Publish(strservice)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
	}
	return r
}

//export Subscribe
func Subscribe(service[]byte)RInfo{
	r := RInfo{}
	r.success = C.int(0)
	strservice := string(service)
	err := call.Subscribe(strservice)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
	}
	return r
}

//export Broadcast
func Broadcast(body,service []byte,info C.BodyInfo)RInfo{
	r := RInfo{}
	r.success = C.int(0)

	infoBody , err := MarshalBody(body,info)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}
	//调用C函数
	service_name := string(service)
	broadResult,err := call.CallBroadcast(infoBody,service_name)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}
	if broadResult.M_Err != nil {
		r.success = C.int(1)
		r.error = C.CString(string(broadResult.M_Err ))
		return r
	}
	list := make([]MsgInfo,0)
	for _ ,v := range broadResult.M_Net_Rsp{
		m := MsgInfo{}
		m.key 		= C.int(v.Key)
		m.error     = C.CString(string(v.CheckErr))
		m.result	= (*C.char)(unsafe.Pointer(&v.Result[0]))
		list = append(list,m)
	}
	r.resultlist = (*C.char)(unsafe.Pointer(&list[0]))
	return r
}

func B()RInfo{
	broadResult := &pb.NetRspInfo{
		M_Err:nil,
		M_Net_Rsp: map[uint32]*pb.SendResultInfo{
			12:&pb.SendResultInfo{
				Key:12,
				CheckErr:nil,
				Result:[]byte("hello struct"),
			},
			13:&pb.SendResultInfo{
				Key:13,
				CheckErr:nil,
				Result:[]byte("hello struct2"),
			},
		},
	}
	r := RInfo{}
	r.success = C.int(1)
	if broadResult.M_Err != nil {
		r.success = C.int(0)
		r.error = C.CString(string(broadResult.M_Err ))
		return r
	}

	list := make([]MsgInfo,0)
	for _ ,v := range broadResult.M_Net_Rsp{
		m := MsgInfo{}
		m.key = C.int(v.Key)
		m.error       = C.CString(string(v.CheckErr))
		m.result = (*C.char)(unsafe.Pointer(&v.Result[0]))

		list = append(list,m)
	}
	fmt.Println(list)
	r.resultlist = (*C.char)(unsafe.Pointer(&list[0]))
	return r
}

//export Send
func Send(body []byte,info C.BodyInfo)RInfo{
	r := RInfo{}
	r.success = C.int(0)

	uuid := ""
	if int32(info.LenID) > 0 {
		uuid = string(C.GoBytes(unsafe.Pointer(info.UUID), info.LenID))
	}

	infoBody , err := MarshalBody(body,info)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}

	//调用C函数
	syncResult,err := call.CallSync(infoBody,uuid)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}
	syncInfo,err := proto.Marshal(syncResult)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}

	r.resultlist = (*C.char)(unsafe.Pointer(&syncInfo[0]))

	return r
}

//export AsyncSend
func AsyncSend(body []byte,info C.BodyInfo)RInfo{
	r := RInfo{}
	r.success = C.int(0)

	infoBody , err := MarshalBody(body,info)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}
	//调用C函数
	asyncResult,err := call.CallAsync(infoBody)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}

	asyncInfo,err := proto.Marshal(asyncResult)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}

	r.resultlist = (*C.char)(unsafe.Pointer(&asyncInfo[0]))
	return r
}

func main(){

}


//env GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -buildmode=c-shared -o clientmsg.dll