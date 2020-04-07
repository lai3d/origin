package GateService

import (
	"fmt"
	"github.com/duanhf2012/origin/event"
	"github.com/duanhf2012/origin/network"
	"github.com/duanhf2012/origin/network/processor"
	"github.com/duanhf2012/origin/node"
	"github.com/duanhf2012/origin/service"
	"github.com/duanhf2012/origin/sysservice"
)

type GateService struct {
	service.Service
	processor network.Processor
	httpRouter sysservice.IHttpRouter
}

func (slf *GateService) OnInit() error{
	tcpervice := node.GetService("TcpService").(*sysservice.TcpService)
	slf.processor = &processor.PBProcessor{}
	tcpervice.SetProcessor(slf.processor)

	httpervice := node.GetService("HttpService").(*sysservice.HttpService)
	slf.httpRouter = sysservice.NewHttpHttpRouter(slf)
	httpervice.SetHttpRouter(slf.httpRouter)

	slf.httpRouter.RegRouter(sysservice.METHOD_GET,"/get/query",slf.HttpTest)
	slf.httpRouter.SetServeFile(sysservice.METHOD_GET,"/img/head/","d:/img")
	return nil
}

func (slf *GateService) HttpTest(session *sysservice.HttpSession) {
	session.SetHeader("a","b")
	session.Write([]byte("this is a test"))
	v,_:=session.Query("a")
	v2,_:=session.Query("b")
	fmt.Print(string(session.GetBody()),"\n",v,"\n",v2)
}

func (slf *GateService) OnEventHandler(ev *event.Event) error{
	if ev.Type == event.Sys_Event_Tcp_RecvPack {
		pPack := ev.Data.(*sysservice.TcpPack)
		slf.processor.Route(ev.Data,pPack.ClientId)
	}else if ev.Type == event.Sys_Event_Tcp_Connected {
		pPack := ev.Data.(*sysservice.TcpPack)
		slf.OnConnected(pPack.ClientId)
	}else if ev.Type == event.Sys_Event_Tcp_DisConnected {
		pPack := ev.Data.(*sysservice.TcpPack)
		slf.OnDisconnected(pPack.ClientId)
	}
	return nil
}

func (slf *GateService) OnConnected(clientid uint64){
	fmt.Printf("client id %d connected",clientid)
}


func (slf *GateService) OnDisconnected(clientid uint64){
	fmt.Printf("client id %d disconnected",clientid)
}
