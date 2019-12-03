package call

import (
	"bufio"
	pb "clientmsg/proto"
	"clientmsg/utils"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"os"
)

func CallSync(mbody []byte)(map[uint32]*pb.SendResultInfo,error){
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return nil,err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)

	ctx := context.Background()

	rsp, err := c.Sync(ctx,&pb.NetReqInfo{M_Body:mbody,Service:""})
	//////////////////////异步处理 ， 调用客户端的接口，异步发送
	if err != nil {
		return nil,err
	}
	if rsp.M_Err != nil {
		return nil,errors.New(string(rsp.M_Err))
	}

	return rsp.M_Net_Rsp,nil

	//for k ,_ := range rsp.M_Net_Rsp{
	//	response := rsp.M_Net_Rsp[k]
	//	fmt.Println("response.SendCount:",response.SendCount)
	//	fmt.Println("response.SuccessCount:",response.SuccessCount)
	//	fmt.Println("response.FailCount:",response.FailCount)
	//	fmt.Println("response.DiscardCount:",response.DiscardCount)
	//	fmt.Println("response.ReSendCount:",response.ReSendCount)
	//	fmt.Println("response.Key:",response.Key)
	//	fmt.Println("response.CheckErr:",string(response.CheckErr))
	//	for kResult,_ := range response.ResultList {
	//		result := response.ResultList[kResult]
	//		fmt.Println("response.SyncType:",result.SyncType)
	//		fmt.Println("response.IsResend:",result.IsResend)
	//		fmt.Println("response.Errinfo:",string(result.Errinfo))
	//		fmt.Println("response.Result:",string(result.Result))
	//	}
	//}
}

func CallAsync(mbody []byte)(map[uint32]*pb.SendResultInfo,error){
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return nil,err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	ctx := context.Background()

	rsp, err := c.Async(ctx,&pb.NetReqInfo{M_Body:mbody,Service:""})

	//////////////////////异步处理 ， 调用客户端的接口，异步发送
	if err != nil {
		return nil,err
	}
	if rsp.M_Err != nil {
		return nil,errors.New(string(rsp.M_Err))
	}
	return rsp.M_Net_Rsp,nil
	//for k ,_ := range rsp.M_Net_Rsp{
	//	response := rsp.M_Net_Rsp[k]
	//	fmt.Println("response.SendCount:",response.SendCount)
	//	fmt.Println("response.SuccessCount:",response.SuccessCount)
	//	fmt.Println("response.FailCount:",response.FailCount)
	//	fmt.Println("response.DiscardCount:",response.DiscardCount)
	//	fmt.Println("response.ReSendCount:",response.ReSendCount)
	//	fmt.Println("response.Key:",response.Key)
	//	fmt.Println("response.CheckErr:",string(response.CheckErr))
	//	for kResult,_ := range response.ResultList {
	//		result := response.ResultList[kResult]
	//		fmt.Println("response.SyncType:",result.SyncType)
	//		fmt.Println("response.IsResend:",result.IsResend)
	//		fmt.Println("response.Errinfo:",string(result.Errinfo))
	//		fmt.Println("response.Result:",string(result.Result))
	//	}
	//}
}

func CallBroadcast(mbody []byte,service string)(map[uint32]*pb.SendResultInfo,error){
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return nil,err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	ctx := context.Background()
	rsp, err := c.Broadcast(ctx,&pb.NetReqInfo{M_Body:mbody,Service:service})
	if err != nil {
		return nil,err
	}
	if rsp.M_Err != nil {
		return nil,errors.New(string(rsp.M_Err))
	}
	return rsp.M_Net_Rsp,nil
}

func GetBody()[]byte{
	fileName := "./2.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return nil
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	bodyByte := make([]byte,110)
	_,err = buf.Read(bodyByte)
	if err != nil {
		return nil
	}
	return bodyByte
}

func Register(sequence , ip , port string)error{
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)

	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	var ctx context.Context
	ctx = context.Background()
	registerReturn , err := c.Register(ctx,&pb.RegisterInfo{
		Sequence:sequence,
		Ip:ip,
		Port:port,
	})
	if err != nil {
		return err
	}
	if registerReturn.Success != true {
		return errors.New(fmt.Sprintf("Register error"))
	}
	return nil
}

func Publish(service string)error{
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	var ctx context.Context
	ctx = context.Background()
	publishReturn , err :=c.Publish(ctx,&pb.PublishInfo{
		Service:service,
	})
	if err != nil {
		return err
	}
	if publishReturn.Success != true {
		return errors.New(fmt.Sprintf("Publish error"))
	}
	return nil
}

func Subscribe(service,ip,port string)error{
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	var ctx context.Context
	ctx = context.Background()
	subscribeInfo,err := c.Subscribe(ctx,&pb.SubscribeInfo{
		Service:service,
		Ip:ip,
		Port:port,
	})
	if err != nil {
		return err
	}
	if subscribeInfo.Success != true {
		return errors.New(fmt.Sprintf("Subscribe error"))
	}
	return nil
}