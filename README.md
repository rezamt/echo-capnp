# Ping Server in CAP'n Proto

# A basic echo service

<p align='center'>
  <img src="./diagrams/ping.svg"/>
</p>

### A basic echo service

Start by writing a [Cap'n Proto schema file][schema].
For example, here is a very simple echo service:

```capnp
interface Echo {
  ping @0 (msg :Text) -> (reply :Text);
}
```

### Passing Capabilities
This version of the protocol adds a heartbeat method. Instead of returning the text directly, it will send it to a 
callback at regular intervals.

```capnp

interface Callback {
  log @0 (msg :Text) -> ();
}

interface Echo {
  ping      @0 (msg :Text) -> (reply :Text);
  heartbeat @1 (msg :Text, callback :Callback) -> ();
}
```

Step 1: The client creates the callback:
<p align='center'>
  <img src="./diagrams/callback1.svg"/>
</p>
Step 2: The client calls the heartbeat method, passing the callback as an argument:
<p align='center'>
  <img src="./diagrams/callback2.svg"/>
</p>
Step 3: The service receives the callback and calls the log method on it:
<p align='center'>
  <img src="./diagrams/callback3.svg"/>
</p>


### Pipelining
Let's say the server also offers a logging service, which the client can get from the main echo service:

```capnp

interface Echo {
  ping      @0 (msg :Text) -> (reply :Text);
  heartbeat @1 (msg :Text, callback :Callback) -> ();
  getLogger @2 () -> (callback :Callback);
}

```
The implementation of the new method in the service is simple - we export the callback in the response in the same way 
we previously exported the client's callback in the request:


### Networking

### How to generate schema?
```bash
cd echo
capnp compile -I$GOPATH/src/zombiezen.com/go/capnproto2/std -ogo * -ocapnp
```