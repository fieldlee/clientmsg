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
	Port = fmt.Sprintf("%d",utils.Port)
)
var callSyncBack      C.ptfFuncCallBack
var callAsyncBack     C.ptfFuncCallBack
var callAnswerBack    C.ptfFuncCall

type RInfo C.CallReturnInfo
type MsgInfo C.MsgReturnInfo
type SingleInfo C.MsgSingleInfo

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
	p := C.malloc(C.sizeof_char * C.ulong(n))
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
func Run()  {
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
func Register(seq []byte)RInfo{
	r := RInfo{}
	r.success = C.int(1)

	strseq := string(seq)
	err := call.Register(strseq)
	if err != nil {
		r.success = C.int(0)
		r.error = C.CString(err.Error())
	}
	return r
}


//export Publish
func Publish(service []byte)RInfo{
	r := RInfo{}
	r.success = C.int(1)

	strservice := string(service)
	err := call.Publish(strservice)
	if err != nil {
		r.success = C.int(0)
		r.error = C.CString(err.Error())
	}
	return r
}

//export Subscribe
func Subscribe(service[]byte)RInfo{
	r := RInfo{}
	r.success = C.int(1)

	strservice := string(service)
	err := call.Subscribe(strservice)
	if err != nil {
		r.success = C.int(0)
		r.error = C.CString(err.Error())
	}
	return r
}

//export Broadcast
func Broadcast(body,service []byte,info C.BodyInfo)RInfo{
	r := RInfo{}
	r.success = C.int(1)

	infoBody , err := MarshalBody(body,info)
	if err != nil {
		r.success = C.int(0)
		r.error = C.CString(err.Error())
		return r
	}
	//调用C函数
	service_name := string(service)

	broadResult,err := call.CallBroadcast(infoBody,service_name)
	if err != nil {
		r.success = C.int(0)
		r.error = C.CString(err.Error())
		return r
	}

	if broadResult.M_Err != nil {
		r.success = C.int(0)
		r.error = C.CString(string(broadResult.M_Err ))
		return r
	}

	list := make([]MsgInfo,0)
	for _ ,v := range broadResult.M_Net_Rsp{
		m := MsgInfo{}
		m.key = C.int(v.Key)
		m.sendcount = C.int(v.SendCount)
		m.successcount = C.int(v.SuccessCount)
		m.failurecount = C.int(v.FailCount)
		m.discardcount = C.int(v.DiscardCount)
		m.resendcount = C.int(v.ReSendCount)
		m.error       = C.CString(string(v.CheckErr))

		slist := make([]SingleInfo,0)
		for _,sv := range v.ResultList{
			sinfo := SingleInfo{}
			sinfo.sequence = C.int(sv.AskSequence)
			sinfo.sendtimeapp = C.int(sv.SendTimeApp)
			sinfo.msgtype     = C.int(sv.MsgType)
			sinfo.msgacktype  = C.int(sv.MsgAckType)
			if sv.IsTimeOut {
				sinfo.istimeout = C.int(1)
			}else{
				sinfo.istimeout = C.int(0)
			}

			if sv.IsDisCard {
				sinfo.isdiscard = C.int(1)
			}else{
				sinfo.isdiscard = C.int(0)
			}

			if sv.IsResend {
				sinfo.isresend  = C.int(1)
			}else{
				sinfo.isresend  = C.int(0)
			}
			if sv.Errinfo != nil {
				sinfo.error = C.CString(string(sv.Errinfo))
			}
			if sv.Result != nil {
				sinfo.result =  C.CString(string(sv.Result))
			}

			slist = append(slist,sinfo)
		}

		m.resultlist = (*C.char)(unsafe.Pointer(&slist[0]))

		list = append(list,m)
	}

	r.result = (*C.char)(unsafe.Pointer(&list[0]))
	return r
}

//export B
func B(broadResult *pb.NetRspInfo)RInfo{
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
		m.sendcount = C.int(v.SendCount)
		m.successcount = C.int(v.SuccessCount)
		m.failurecount = C.int(v.FailCount)
		m.discardcount = C.int(v.DiscardCount)
		m.resendcount = C.int(v.ReSendCount)
		m.error       = C.CString(string(v.CheckErr))

		slist := make([]SingleInfo,0)
		for _,sv := range v.ResultList{
			sinfo := SingleInfo{}
			sinfo.sequence = C.int(sv.AskSequence)
			sinfo.sendtimeapp = C.int(sv.SendTimeApp)
			sinfo.msgtype     = C.int(sv.MsgType)
			sinfo.msgacktype  = C.int(sv.MsgAckType)
			if sv.IsTimeOut {
				sinfo.istimeout = C.int(1)
			}else{
				sinfo.istimeout = C.int(0)
			}

			if sv.IsDisCard {
				sinfo.isdiscard = C.int(1)
			}else{
				sinfo.isdiscard = C.int(0)
			}

			if sv.IsResend {
				sinfo.isresend  = C.int(1)
			}else{
				sinfo.isresend  = C.int(0)
			}
			if sv.Errinfo != nil {
				sinfo.error = C.CString(string(sv.Errinfo))
			}
			if sv.Result != nil {
				sinfo.result =  C.CString(string(sv.Result))
			}

			slist = append(slist,sinfo)
		}

		m.resultlist = (*C.char)(unsafe.Pointer(&slist[0]))

		list = append(list,m)
	}
	fmt.Println(list)
	r.result = (*C.char)(unsafe.Pointer(&list[0]))
	return r
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


