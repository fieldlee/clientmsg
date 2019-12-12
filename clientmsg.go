
package clientmsg

import "C"

import (
	"clientmsg/call"
	pb "clientmsg/proto"
	"clientmsg/utils"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"context"
	"net"
)

var (
	Host = utils.Address
	Port = fmt.Sprintf("%d",utils.Port)
)

type MsgHandle struct {}

//type HandleMidMsg struct {
//	Handle func([]byte)([]byte,error)
//	AsyncHandle func(*pb.SingleResultInfo)([]byte,error)
//}
//
//var HandleObj = HandleMidMsg{}

func (m *MsgHandle)Call(ctx context.Context, info *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}

	/////调用C函数

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
	// 调用C函数

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
func Register(seq,ip string){
	err := call.Register(seq,ip,"8989")
	if err != nil {
		return
	}
}

//export Publish
func Publish(service string){
	err := call.Publish(service)
	if err != nil {
		return
	}
}

//export Subscribe
func Subscribe(service,ip string){
	err := call.Subscribe(service,ip,"8989")
	if err != nil {
		return
	}
}

//export Sync
func Sync(body []byte){
	//调用C函数
	syncResult,err := call.CallSync(body)
	if err != nil {

	}
	//syncResult.
}

//export Async
func Async(body []byte){
	//调用C函数
	syncResult,err := call.CallAsync(body)
	if err != nil {

	}
}



