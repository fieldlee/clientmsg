package handle

import (
	pb "clientmsg/proto"
	"context"
)

type MsgHandle struct {}

type HandleMidMsg struct {
	Handle func([]byte)([]byte,error)
	AsyncHandle func(*pb.SingleResultInfo)([]byte,error)
}

var HandleObj = HandleMidMsg{}

func (m *MsgHandle)Call(ctx context.Context, info *pb.CallReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	if HandleObj.Handle == nil {
		out.M_Net_Rsp = []byte("The Handle Call function not instance")
	}else{

		//info.Service
		//info.M_Body.M_MsgBody.MLBack
		//info.M_Body.M_MsgBody.MSSendCount
		//info.M_Body.M_MsgBody.MLServerSequence
		//info.M_Body.M_MsgBody.MLExpireTime
		//info.M_Body.M_MsgBody.MLAskSequence
		//info.M_Body.M_MsgBody.MISendTimeApp
		//info.M_Body.M_MsgBody.MCMsgType
		//info.M_Body.M_MsgBody.MCMsgAckType
		//info.M_Body.M_MsgBody.MIDiscard
		//info.M_Body.M_MsgBody.MLAsktype
		//info.M_Body.M_MsgBody.MLResult

		reT,err := HandleObj.Handle(info.M_Body.M_Msg)
		if err != nil {
			out.M_Net_Rsp = []byte(err.Error())
		}
		out.M_Net_Rsp = reT
	}
	return &out,nil
}

func (m *MsgHandle)AsyncCall(ctx context.Context, resultInfo *pb.SingleResultInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	if HandleObj.AsyncHandle == nil {
		out.M_Net_Rsp = []byte("The AsyncHandle Call function not instance")
	}else{
		reT,err := HandleObj.AsyncHandle(resultInfo)
		if err != nil {
			out.M_Net_Rsp = []byte(err.Error())
		}
		out.M_Net_Rsp = reT
	}
	return &out,nil
}
