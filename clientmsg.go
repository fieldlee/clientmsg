package main
// #include "bridge.h"
import "C"
import (
	"clientmsg/call"
	pb "clientmsg/proto"
	"clientmsg/utils"
	"context"
	"fmt"
	"github.com/pborman/uuid"
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
var callSyncBack    C.ptfSyncCallBack
var notifyBack   C.ptfFuncCall

type RInfo C.CallReturnInfo

//export SetSyncReturnBack
func SetSyncReturnBack(f C.ptfSyncCallBack) {
	callSyncBack = f
}

//export SetNotifyBack
func SetNotifyBack(f C.ptfFuncCall) {
	notifyBack = f
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
	data := C.CStr{
		content:(*C.char)(unsafe.Pointer(&rq[0])),
		length:C.int(len(rq)),
	}
	///客户端ip
	ip := C.CStr{
		content:(*C.char)(unsafe.Pointer(&ip_byte[0])),
		length:C.int(len(ip_byte)),
	}

	if info.Uuid == ""{ // 同步请求
		/////调用C函数
		p := C.CSyncHandleData(callSyncBack,data,C.ulonglong(seqno),ip)
		defer C.free(unsafe.Pointer(p.content))

		resultByte := C.GoBytes(unsafe.Pointer(p.content), p.length)
		out.M_Net_Rsp = resultByte

		return &out,nil
	}else{  // 异步请求

		go func(uid string,data C.CStr,seqno uint64,ip C.CStr) {

			p := C.CSyncHandleData(callSyncBack,data,C.ulonglong(seqno),ip)

			defer C.free(unsafe.Pointer(p.content))

			resultByte := C.GoBytes(unsafe.Pointer(p.content), p.length)

			infoBody , err := MarshalBody(resultByte,C.BodyInfo{},false)
			if err != nil {
				fmt.Println(fmt.Sprintf("MarshalBody error : %s",err.Error()))
				return
			}

			syncResult,err := call.CallSync(infoBody,uid)

			if err != nil {
				fmt.Println(fmt.Sprintf("call.CallSync error : %s",err.Error()))
				return
			}

			if syncResult.M_Err != nil {
				fmt.Println(fmt.Sprintf("syncResult error : %s",syncResult.M_Err))
				return
			}
			return
		}(info.Uuid,data,seqno,ip)

		out.M_Net_Rsp = []byte{0}

		return &out,nil
	}

}

func (m *MsgHandle)AsyncAnswer(ctx context.Context, resultInfo *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	////获取msg body
	//body := utils.AsyncBody{
	//	Uid:resultInfo.Uuid,
	//	Body:resultInfo.M_Body.M_Msg,
	//}

	data := C.CStr{
		content:(*C.char)(unsafe.Pointer(&resultInfo.M_Body.M_Msg[0])),
		length:C.int(len(resultInfo.M_Body.M_Msg)),
	}

	result := C.CHandleCall(notifyBack,data,0)

	out.M_Net_Rsp = utils.Int32ToBytes(int32(result))

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

func errTransRinfo(err error)RInfo{
	errByte := []byte(err.Error())
	r := RInfo{}
	r.success = C.int(1)
	r.error = (*C.char)(unsafe.Pointer(&errByte[0]))
	r.length = C.int(len(errByte))
	return r
}

func errByteTransRinfo(err []byte)RInfo{
	r := RInfo{}
	r.success = C.int(1)
	r.error = (*C.char)(unsafe.Pointer(&err[0]))
	r.length = C.int(len(err))
	return r
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
		return errTransRinfo(err)
	}

	err = call.Register(strconv.FormatUint(serno,10))
	if err != nil {
		return errTransRinfo(err)
	}
	body_byte := []byte("register success")
	r.result = (*C.char)(unsafe.Pointer(&body_byte[0]))
	r.length = C.int(len(body_byte))
	return r
}

//export Publish
func Publish(service []byte)RInfo{
	r := RInfo{}
	r.success = C.int(0)
	strservice := string(service)

	serno,err := strconv.ParseUint(strservice,10,64)
	if err != nil {
		return errTransRinfo(err)
	}

	err = call.Publish(strconv.FormatUint(serno,10))
	if err != nil {
		return errTransRinfo(err)
	}
	body_byte := []byte("public success")
	r.result = (*C.char)(unsafe.Pointer(&body_byte[0]))
	r.length = C.int(len(body_byte))
	return r
}

//export Subscribe
func Subscribe(service[]byte)RInfo{
	r := RInfo{}
	r.success = C.int(0)
	strservice := string(service)

	serno,err := strconv.ParseUint(strservice,10,64)
	if err != nil {
		return errTransRinfo(err)
	}

	err = call.Subscribe(strconv.FormatUint(serno,10))
	if err != nil {
		return errTransRinfo(err)
	}

	body_byte := []byte("subscribe success")
	r.result = (*C.char)(unsafe.Pointer(&body_byte[0]))
	r.length = C.int(len(body_byte))
	return r
}

//export Broadcast
func Broadcast(body,service []byte,info C.BodyInfo)RInfo{
	r := RInfo{}
	r.success = C.int(0)

	infoBody , err := MarshalBody(body,info,true)
	if err != nil {
		return errTransRinfo(err)
	}
	//调用C函数
	service_name := string(service)

	serno,err := strconv.ParseUint(service_name,10,64)
	if err != nil {
		return errTransRinfo(err)
	}

	broadResult,err := call.CallBroadcast(infoBody,strconv.FormatUint(serno,10))
	if err != nil {
		return errTransRinfo(err)
	}
	if broadResult != nil {
		if broadResult.M_Err != nil {
			return errByteTransRinfo(broadResult.M_Err)
		}

		var body_byte []byte
		for _ ,v := range broadResult.M_Net_Rsp{
			body_byte = v.Result
		}
		r.result = (*C.char)(unsafe.Pointer(&body_byte[0]))
		r.length = C.int(len(body_byte))
	}
	return r
}

//export Send
func Send(body []byte,info C.BodyInfo)RInfo{
	r := RInfo{}
	r.success = C.int(0)

	uuid := ""

	infoBody , err := MarshalBody(body,info,true)
	if err != nil {
		return errTransRinfo(err)
	}

	//调用C函数
	syncResult,err := call.CallSync(infoBody,uuid)
	if err != nil {
		return errTransRinfo(err)
	}

	if syncResult != nil {
		if syncResult.M_Err != nil {
			return errByteTransRinfo(syncResult.M_Err)
		}

		var body_result []byte
		for _ , v := range  syncResult.M_Net_Rsp{
			body_result = v.Result
		}
		r.result = (*C.char)(unsafe.Pointer(&body_result[0]))
		r.length = C.int(len(body_result))
	}
	return r
}

//export AsyncSend
func AsyncSend(body []byte,info C.BodyInfo)RInfo{
	r := RInfo{}
	r.success = C.int(0)

	uid := uuid.New()

	infoBody , err := MarshalBody(body,info,true)
	if err != nil {
		return errTransRinfo(err)
	}
	//调用C函数
	asyncResult,err := call.CallAsync(infoBody,uid)
	if err != nil {
		return errTransRinfo(err)
	}

	if asyncResult != nil {
		if asyncResult.M_Err != nil {
			return errByteTransRinfo(asyncResult.M_Err)
		}

		var body_result []byte
		for _,v := range asyncResult.M_Net_Rsp{
			body_result = v.Result
		}
		r.result = (*C.char)(unsafe.Pointer(&body_result[0]))
		r.length = C.int(len(body_result))
	}

	return r
}

//export AsyncSendCallback
func AsyncSendCallback(body []byte,info C.BodyInfo,cb C.ptfFuncCall){
	uid := ""
	infoBody , err := MarshalBody(body,info,true)
	if err != nil {
		errByte := []byte(err.Error())
		errdata := C.CStr{
			content:(*C.char)(unsafe.Pointer(&errByte[0])),
			length:C.int(len(errByte)),
		}
		C.CHandleCall(cb,errdata,1)
	}
	//调用C函数
	asyncResult,err := call.CallSync(infoBody,uid)
	if err != nil {
		errByte := []byte(err.Error())
		errdata := C.CStr{
			content:(*C.char)(unsafe.Pointer(&errByte[0])),
			length:C.int(len(errByte)),
		}
		C.CHandleCall(cb,errdata,1)
	}

	if asyncResult != nil {
		if asyncResult.M_Err != nil {
			errdata := C.CStr{
				content:(*C.char)(unsafe.Pointer(&asyncResult.M_Err[0])),
				length:C.int(len(asyncResult.M_Err)),
			}
			C.CHandleCall(cb,errdata,1)
		}

		var body_result []byte
		for _,v := range asyncResult.M_Net_Rsp{
			body_result = v.Result
		}
		resultdata := C.CStr{
			content:(*C.char)(unsafe.Pointer(&body_result[0])),
			length:C.int(len(body_result)),
		}
		C.CHandleCall(cb,resultdata,0)
	}
}

func main(){

}

//env GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -buildmode=c-shared -o clientmsg.dll