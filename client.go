package main

import (
	"context"
	"github.com/rezamt/ping-cap/echo"
	"log"
	"net"
	"zombiezen.com/go/capnproto2/rpc"
)

type Callback_Server_Impl struct{}

func (h Callback_Server_Impl) Log(call echo.Callback_log) error {

	msg, _ := call.Params.Msg()

	log.Printf("From Server: %s", msg)

	return nil
}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		log.Fatal("Failed to obtain connection")
	}

	t2 := rpc.StreamTransport(conn)

	// Client-side
	ctx := context.Background()
	clientConn := rpc.NewConn(t2)
	defer clientConn.Close()

	p := echo.Echo{clientConn.Bootstrap(ctx)}

	result, err := p.Ping(ctx, func(params echo.Echo_ping_Params) error {
		params.SetMsg("Hello World Ping Me!")
		return nil
	}).Struct()

	if err != nil {
		log.Fatal(err)
	}

	msg, err := result.Reply()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received: %s", msg)

	_, err = p.Heartbeat(ctx, func(params echo.Echo_heartbeat_Params) error {

		cb := echo.Callback_ServerToClient(Callback_Server_Impl{})

		params.SetCallback(cb)

		params.SetMsg("Fuck you")

		return nil
	}).Struct()

	if err != nil {
		log.Fatal(err)
	}

}
