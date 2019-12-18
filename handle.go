package main
// #include "bridge.h"
import "C"
import (
	pb "clientmsg/proto"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"clientmsg/model"
	"clientmsg/utils"
	"time"
	"math/rand"
)

func randInt(min int , max int) uint32 {
	rand.Seed( time.Now().UTC().UnixNano())
	return uint32(min + rand.Intn(max-min))
}

func MarshalBody(body []byte,info C.BodyInfo)([]byte,error){
	net_msgbody := &pb.Net_Pack_Min_Net_MsgBody{
		MLAsktype:uint64(info.Asktype),
		MLServerSequence:uint64(info.ServerSequence),
		MLAskSequence:uint64(info.AskSequence),
		MCMsgAckType:int32(info.MsgAckType),
		MCMsgType:int32(info.MsgType),
		MSSendCount:int32(info.SendCount),
		MLExpireTime:uint32(info.ExpireTime),
		MISendTimeApp:uint64(info.SendTimeApp),
		MLResult:int32(info.Result),
		MLBack:uint64(info.Back),
		MIDiscard:int32(info.Discard),
	}

	net_pack := &pb.Net_Pack{
		M_Msg:body,
		M_MsgBody:net_msgbody,
	}


	gj_net_pack := &pb.GJ_Net_Pack{
		M_Net_Pack: map[uint32]*pb.Net_Pack{randInt(10000,1000000):net_pack},
	}

	gjbody,err := proto.Marshal(gj_net_pack)

	if err != nil {
		return nil,err
	}

	if model.COMPRESS_TYPE(info.Compress) < model.Compression_no || model.COMPRESS_TYPE(info.Compress) >= model.CompressionWayMax{
		return nil,errors.New("compress way error")
	}

	if model.ENCRPTION_TYPE(info.Encrypt) < model.Encryption_No || model.ENCRPTION_TYPE(info.Encrypt) >= model.Encryption_Max{
		return nil,errors.New("encrypt way error")
	}

	fullbody := FullHead(gjbody,info.Compress,info.Encrypt)

	return fullbody,nil
}


func FullHead(inbody []byte,compress ,encryptType int)[]byte{
	headINfo := model.HeadInfo{
		Tag:model.HeadTag,
		Version:int16(model.HeadVersion),
		ClientType:int16(model.HeadClientType),
		HeadLength:int16(model.HeadLength),
		CompressWay:uint8(compress),
		Encryption:uint8(encryptType),
		Sig:uint8(model.HeadSig),
		Format:uint8(model.HeadFormat),
		NetFlag:uint8(model.HeadNetFlag),
		Back1:0,
		BufSize:int32(len(inbody)),
		UncompressedSize:int32(len(inbody)),
		Back2:0,
	}
	encryptByte := inbody
	//////加密
	switch model.ENCRPTION_TYPE(headINfo.Encryption) {
	case model.Encryption_AES:
		encryptByte,_ = utils.EncryptAes(inbody,[]byte(model.PassPass))
	case model.Encryption_RSA:
		pubkey := utils.BytesToPublicKey(inbody)
		encryptByte = utils.EncryptWithPublicKey(inbody,pubkey)
	case model.Encryption_Des:
		encryptByte = utils.Encrypt3DES(inbody,[]byte(model.PassPass))
	}
	inbody = encryptByte
	/////压缩body bytes
	if model.COMPRESS_TYPE(headINfo.CompressWay) == model.Compression_zip {
		if zipbody,err := utils.ZipByte(inbody);err != nil {
			fmt.Println(err.Error())
		}else{
			inbody = zipbody
		}
	}
	/////修改压缩后的buffer长度
	headINfo.BufSize = int32(len(inbody))

	inbody = utils.JoinHeadAndBody(headINfo,inbody)

	return nil
}
