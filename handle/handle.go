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

func (m *MsgHandle)Call(ctx context.Context, info *pb.NetReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	if HandleObj.Handle == nil {
		out.M_Net_Rsp = []byte("The Handle Call function not instance")
	}else{
		reT,err := HandleObj.Handle(info.M_Body)
		if err != nil {
			out.M_Net_Rsp = []byte(err.Error())
		}
		out.M_Net_Rsp = reT
	}

	//fmt.Println("netPack.M_MsgBody.MLAsktype:",netPack.M_MsgBody.MLAsktype)
	//fmt.Println("netPack.M_MsgBody.MIDiscard:",netPack.M_MsgBody.MIDiscard)
	//fmt.Println("netPack.M_MsgBody.MCMsgAckType:",netPack.M_MsgBody.MCMsgAckType)
	//fmt.Println("netPack.M_MsgBody.MCMsgType:",netPack.M_MsgBody.MCMsgType)
	//fmt.Println("netPack.M_MsgBody.MISendTimeApp:",netPack.M_MsgBody.MISendTimeApp)
	//fmt.Println("netPack.M_MsgBody.MLAskSequence:",netPack.M_MsgBody.MLAskSequence)
	//fmt.Println("netPack.M_MsgBody.MLExpireTime:",netPack.M_MsgBody.MLExpireTime)
	//fmt.Println("netPack.M_MsgBody.MLServerSequence:",netPack.M_MsgBody.MLServerSequence)
	//fmt.Println("netPack.M_MsgBody.MSSendCount:",netPack.M_MsgBody.MSSendCount)
	//fmt.Println("netPack.M_MsgBody.MLBack:",netPack.M_MsgBody.MLBack)
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
