
package clientmsg

import "C"

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

//export Run
func Run()  {
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

//export Register
func Register(seq,ip string){
	call.Register(seq,ip,"8989")
}

//export Publish
func Publish(service string){
	call.Publish(service)
}

//export Subscribe
func Subscribe(service,ip string){
	call.Subscribe(service,ip,"8989")
}

//export Sync
func Sync(body []byte){
	syncResult,err := call.CallSync(body)
	if err != nil {

	}
	//syncResult.
}

//export Async
func Async(body []byte){
	syncResult,err := call.CallAsync(body)
	if err != nil {

	}
}