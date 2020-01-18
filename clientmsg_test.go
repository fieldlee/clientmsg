package main

import (
	"testing"
)

func TestTest(t *testing.T){
	//r := Test()
	//fmt.Println(r.success)
	//fmt.Println(r.error)
}

func TestB2(t *testing.T) {
	/***


	broadResult := &pb.NetRspInfo{
		M_Err:nil,
		M_Net_Rsp: map[uint32]*pb.SendResultInfo{
			12:&pb.SendResultInfo{
				Key:12,
				SendCount:2,
				SuccessCount:1,
				FailCount:1,
				DiscardCount:1,
				ReSendCount: 0,
				CheckErr:nil,
				ResultList: map[uint32]*pb.SingleResultInfo{
					0:&pb.SingleResultInfo{
						AskSequence:1,
						SendTimeApp:1,
						MsgAckType:1,
						MsgType:2,
						IsResend:true,
						IsDisCard:false,
						IsTimeOut:true,
						Errinfo:[]byte(errors.New("IS TIME OUT").Error()),
						Result:nil,
					},
					1:&pb.SingleResultInfo{
						AskSequence:2,
						SendTimeApp:1,
						MsgAckType:1,
						MsgType:2,
						IsResend:false,
						IsDisCard:false,
						IsTimeOut:false,
						Errinfo:nil,
						Result:[]byte("hello struct"),
					},
				},
			},
			13:&pb.SendResultInfo{
				Key:13,
				SendCount:2,
				SuccessCount:1,
				FailCount:1,
				DiscardCount:1,
				ReSendCount: 0,
				CheckErr:nil,
				ResultList: map[uint32]*pb.SingleResultInfo{
					0:&pb.SingleResultInfo{
						AskSequence:1,
						SendTimeApp:1,
						MsgAckType:1,
						MsgType:2,
						IsResend:true,
						IsDisCard:false,
						IsTimeOut:true,
						Errinfo:[]byte(errors.New("IS TIME OUT").Error()),
						Result:nil,
					},
					1:&pb.SingleResultInfo{
						AskSequence:2,
						SendTimeApp:1,
						MsgAckType:1,
						MsgType:2,
						IsResend:false,
						IsDisCard:false,
						IsTimeOut:false,
						Errinfo:nil,
						Result:[]byte("hello struct2"),
					},
				},
			},
		},
	}
	*/

}

