package handle

import (
	pb "clientmsg/proto"
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
)

type MsgHandle struct {}

func (m *MsgHandle)Call(ctx context.Context, info *pb.NetReqInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}

	/////////////////////////////////////////////
	netPack := pb.Net_Pack{}
	err  := proto.Unmarshal(info.M_Body,&netPack)

	if err != nil {
		return nil , err
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
	out.M_Net_Rsp = []byte("call info to return")
	return &out,nil
}

func (m *MsgHandle)AsyncCall(ctx context.Context, resultInfo *pb.SingleResultInfo) (*pb.CallRspInfo, error) {
	out := pb.CallRspInfo{}
	/////////////////////////////////////////////

	fmt.Println("AsyncCall info:",resultInfo)

	//fmt.Println("resultInfo.IsResend:",resultInfo.IsResend)
	//fmt.Println("resultInfo.IsDisCard:",resultInfo.IsDisCard)
	//fmt.Println("resultInfo.IsTimeOut:",resultInfo.IsTimeOut)
	//fmt.Println("resultInfo.Result:",string(resultInfo.Result))
	//fmt.Println("resultInfo.Errinfo:",string(resultInfo.Errinfo))
	//fmt.Println("resultInfo.SyncType:",resultInfo.SyncType)
	//fmt.Println("resultInfo.MsgAckType:",resultInfo.MsgAckType)
	//fmt.Println("resultInfo.SendTimeApp:",resultInfo.SendTimeApp)
	//fmt.Println("resultInfo.AskSequence:",resultInfo.AskSequence)

	return &out,nil
}
