package main

import (
	"fmt"
	"github.com/rezamt/ping-cap/echo"
	"golang.org/x/net/context"
	"net"
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

func server(c net.Conn) error {
	// Create a new locally implemented HashFactory.
	main := echo.Echo_ServerToClient(Echo_Server_Impl{})
	// Listen for calls, using the HashFactory as the bootstrap interface.
	conn := rpc.NewConn(rpc.StreamTransport(c), rpc.MainInterface(main.Client))
	// Wait for connection to abort.
	err := conn.Wait()
	return err
}

func client(ctx context.Context, c net.Conn) error {
	// Create a connection that we can use to get the Ping
	conn := rpc.NewConn(rpc.StreamTransport(c))
	defer conn.Close()
	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	ec := echo.Echo{Client: conn.Bootstrap(ctx)}

	// Now we can call ping methods
	promise := ec.Ping(ctx, func(p echo.Echo_ping_Params) error {
		err := p.SetMsg("Hello, YO YO !")
		return err
	})

	result, err := promise.Struct()
	if err != nil {
		return err
	}
	reply, err := result.Reply()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", reply)
	return nil
}

func main() {
	c1, c2 := net.Pipe()
	go server(c1)
	client(context.Background(), c2)

}
