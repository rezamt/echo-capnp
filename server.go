package main

import (
	"github.com/rezamt/ping-cap/echo"
	"log"
	"net"
	"time"
	"zombiezen.com/go/capnproto2/rpc"
)

type Echo_Server_Impl struct{}

func (e Echo_Server_Impl) Ping(call echo.Echo_ping) error {
	result, err := call.Params.Msg()
	if err != nil {
		return err
	}

	return call.Results.SetReply(e.reverse(result))
}

func (e Echo_Server_Impl) reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func (e Echo_Server_Impl) Heartbeat(call echo.Echo_heartbeat) error {

	if call.Params.HasCallback() {
		log.Println("Call back has been passed")

		callback := call.Params.Callback()

		for i := 0; i < 2; i++ {
			s, _ := callback.Log(call.Ctx, func(params echo.Callback_log_Params) error {
				return params.SetMsg("Fuck You AAS HOLE ")
			}).Struct()

			time.Sleep(20 * time.Second)

			log.Println(s.IsValid())
		}

		/*
			s, _ := callback.Log(call.Ctx, func(params echo.Callback_log_Params) error {
				return params.SetMsg("Fuck You AAS HOLE ")
			}).Struct()

			fmt.Println(s.IsValid())
		*/
	} else {
		log.Println("No Callback")
	}

	msg, _ := call.Params.Msg()

	log.Printf("Client Msg is: %s", msg)

	return nil
}

func handleRPCRequest(conn net.Conn) error {

	srv := echo.Echo_ServerToClient(Echo_Server_Impl{})

	transport := rpc.StreamTransport(conn)

	serverConn := rpc.NewConn(transport, rpc.MainInterface(srv.Client))

	log.Println("SERVER: Goodbye :)")

	err := serverConn.Wait()

	return err
}

func main() {
	ln, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatal("Failed to bind to local port 7777")
	}

	for {

		log.Println("Waiting for a request")
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal("Failed to accept incoming connection")
		}

		go handleRPCRequest(conn)
	}
}
