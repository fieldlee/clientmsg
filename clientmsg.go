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
	"strconv"
	"syscall"
	"unsafe"
)

var (
	Host = utils.Address
	Port = fmt.Sprintf("%d",utils.Port)
)
var callSyncBack      C.ptfSyncCallBack
var callAsyncBack     C.ptfAsyncCallBack
var callAnswerBack    C.ptfFuncCall

type RInfo C.CallReturnInfo
type MsgInfo C.MsgReturnInfo
type StrInfo C.CString

//export SetSyncReturnBack
func SetSyncReturnBack(f C.ptfSyncCallBack) {
	callSyncBack = f
}

//export SetAsyncReturnBack
func SetAsyncReturnBack(f C.ptfAsyncCallBack) {
	callAsyncBack = f
}

//export SetAnswerBack
func SetAnswerBack(f C.ptfFuncCall) {
	callAnswerBack = f
}

type MsgHandle struct {}

func (m *MsgHandle)Call(ctx context.Context, info *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}

	////获取msg body
	rq := info.M_Body.M_Msg

	var seqno uint64
	var err  error
	if info.Service == "" {
		seqno = info.M_Body.M_MsgBody.MLServerSequence
	}else{
		seqno,err = strconv.ParseUint(info.Service,10,64)
		if err != nil {
			return &out,err
		}
	}

	ip_byte := []byte(info.Clientip)
	///同步请求数据
	data := StrInfo{
		content:(*C.char)(unsafe.Pointer(&rq[0])),
		length:C.int(len(rq)),
	}
	///客户端ip
	ip := StrInfo{
		content:(*C.char)(unsafe.Pointer(&ip_byte[0])),
		length:C.int(len(ip_byte)),
	}
	/////调用C函数
	p := C.CSyncHandleData(callSyncBack,data,C.ulonglong(seqno),ip)
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
	//resultInfo.Uuid
	//// uuid bytes
	uid_byte := []byte(resultInfo.Uuid)
	////获取msg body
	var seqno uint64
	var err  error
	if resultInfo.Service == "" {
		seqno = resultInfo.M_Body.M_MsgBody.MLServerSequence
	}else{
		seqno,err = strconv.ParseUint(resultInfo.Service,10,64)
		if err != nil {
			return &out,err
		}
	}

	ip_byte := []byte(resultInfo.Clientip)

	rq := resultInfo.M_Body.M_Msg
	////异步传送数据
	data := StrInfo{
		content:(*C.char)(unsafe.Pointer(&rq[0])),
		length:C.int(len(rq)),
	}
	////客户端ip
	ip := StrInfo{
		content:(*C.char)(unsafe.Pointer(&ip_byte[0])),
		length:C.int(len(ip_byte)),
	}
	///异步请求uid
	uid := StrInfo{
		content:(*C.char)(unsafe.Pointer(&uid_byte[0])),
		length:C.int(len(uid_byte)),
	}
	/////调用C函数
	p := C.CAsyncHandleData(callAsyncBack,data,uid,C.ulonglong(seqno),ip)
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
	////回调数据
	data := StrInfo{
		content:(*C.char)(unsafe.Pointer(&rq[0])),
		length:C.int(len(rq)),
	}
	/////调用C函数
	result := C.CHandleCall(callAnswerBack,data)
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
	serno,err := strconv.ParseUint(strseq,10,64)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}

	err = call.Register(strconv.FormatUint(serno,10))
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

	serno,err := strconv.ParseUint(strservice,10,64)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}

	err = call.Publish(strconv.FormatUint(serno,10))
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

	serno,err := strconv.ParseUint(strservice,10,64)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}

	err = call.Subscribe(strconv.FormatUint(serno,10))
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

	serno,err := strconv.ParseUint(service_name,10,64)
	if err != nil {
		r.success = C.int(1)
		r.error = C.CString(err.Error())
		return r
	}

	broadResult,err := call.CallBroadcast(infoBody,strconv.FormatUint(serno,10))
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