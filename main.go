
package main

import (
	"clientmsg/call"
	"clientmsg/handle"
	pb "clientmsg/proto"
	"clientmsg/utils"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

var (
	Host = utils.Address
	Port = fmt.Sprintf("%d",utils.Port)
)

func main()  {
	listener, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		log.Fatalln("faile listen at: " + Host + ":" + Port)
	} else {
		log.Println("server is listening at: " + Host + ":" + Port)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterClientServiceServer(rpcServer,&handle.MsgHandle{})
	reflection.Register(rpcServer)
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatalln("failed serve at: " + Host + ":" + Port)
	}
}

func testSync(){
	body := call.GetBody()
	time.Sleep(time.Second)
	call.CallSync(body)
}

func testRegister(){
	call.Register("100002","192.168.1.31","8989")
}
func testPublish(){
	//call.Publish("test.service3")
}
func testSubscribe(){
	call.Publish("test.service3")
	call.Subscribe("test.service3","192.168.1.31","8989")
}
func testAsync(){
	body := call.GetBody()
	call.CallAsync(body)
}