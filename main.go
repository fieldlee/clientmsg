
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
)

var (
	Host = utils.Address
	Port = fmt.Sprintf("%d",utils.Port)
)

func main()  {
	fmt.Println(Port)
	listener, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		log.Fatalln("faile listen at: " + Host + ":" + Port)
	} else {
		log.Println("server is listening at: " + Host + ":" + Port)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterClientServiceServer(rpcServer,&handle.MsgHandle{})
	reflection.Register(rpcServer)
	go test()
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatalln("faile serve at: " + Host + ":" + Port)
	}
}

func test(){
	for i:= 0;i<10000;i++{
		fmt.Println("iiiii::",i)
		call.CallSync()
	}
}