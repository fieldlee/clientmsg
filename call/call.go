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

func CallSync(){
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

	body := GetBody()
	//fmt.Println(string(body))
	r, err := c.Sync(ctx,&pb.NetReqInfo{M_Body:body})

	//////////////////////异步处理 ， 调用客户端的接口，异步发送
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("M_Err:",string(r.M_Err))

	fmt.Println("r.M_Net_Rsp:",r.M_Net_Rsp)
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