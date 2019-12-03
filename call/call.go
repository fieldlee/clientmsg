package call

import (
	"bufio"
	pb "clientmsg/proto"
	"clientmsg/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"os"
)

func CallSync(m_body []byte, num int){
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	fmt.Println(caddr)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	var ctx context.Context

	ctx = context.Background()

	for i := 0 ;i< num;i++{
		//fmt.Println(string(body))
		//for j :=0;j<1;j++{
		//	go func(ctx context.Context,client pb.MidServiceClient) {
				rsp, err := c.Sync(ctx,&pb.NetReqInfo{M_Body:m_body,Service:""})
				//////////////////////异步处理 ， 调用客户端的接口，异步发送
				if err != nil {
					fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",err.Error())
					return
				}
				if rsp.M_Err != nil {
					fmt.Println("==============================",string(rsp.M_Err))
				}

				for k ,_ := range rsp.M_Net_Rsp{
					response := rsp.M_Net_Rsp[k]
					fmt.Println("response.SendCount:",response.SendCount)
					fmt.Println("response.SuccessCount:",response.SuccessCount)
					fmt.Println("response.FailCount:",response.FailCount)
					fmt.Println("response.DiscardCount:",response.DiscardCount)
					fmt.Println("response.ReSendCount:",response.ReSendCount)
					fmt.Println("response.Key:",response.Key)
					fmt.Println("response.CheckErr:",string(response.CheckErr))
					for kResult,_ := range response.ResultList {
						result := response.ResultList[kResult]
						fmt.Println("response.SyncType:",result.SyncType)
						fmt.Println("response.IsResend:",result.IsResend)
						fmt.Println("response.Errinfo:",string(result.Errinfo))
						fmt.Println("response.Result:",string(result.Result))
					}
				}

			//	<- time.After(time.Microsecond*5)
			//}(ctx,c)
		//	<- time.After(time.Microsecond*50)
		//}

	}


}


func CallAsync(m_body []byte, num int){
	caddr := fmt.Sprintf("%v:%v",utils.ServerAddr,utils.ServerPort)
	fmt.Println(caddr)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewMidServiceClient(conn)
	var ctx context.Context

	ctx = context.Background()
	for i := 0 ;i< num;i++{
		//fmt.Println(string(body))
		rsp, err := c.Async(ctx,&pb.NetReqInfo{M_Body:m_body,Service:""})

		//////////////////////异步处理 ， 调用客户端的接口，异步发送
		if err != nil {
			fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",err.Error())
			return
		}
		if rsp.M_Err != nil {
			fmt.Println("==============================",string(rsp.M_Err))
		}

		for k ,_ := range rsp.M_Net_Rsp{
			response := rsp.M_Net_Rsp[k]
			fmt.Println("response.SendCount:",response.SendCount)
			fmt.Println("response.SuccessCount:",response.SuccessCount)
			fmt.Println("response.FailCount:",response.FailCount)
			fmt.Println("response.DiscardCount:",response.DiscardCount)
			fmt.Println("response.ReSendCount:",response.ReSendCount)
			fmt.Println("response.Key:",response.Key)
			fmt.Println("response.CheckErr:",string(response.CheckErr))
			for kResult,_ := range response.ResultList {
				result := response.ResultList[kResult]
				fmt.Println("response.SyncType:",result.SyncType)
				fmt.Println("response.IsResend:",result.IsResend)
				fmt.Println("response.Errinfo:",string(result.Errinfo))
				fmt.Println("response.Result:",string(result.Result))
			}
		}

	}

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