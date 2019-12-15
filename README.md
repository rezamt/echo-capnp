# Ping Server in CAP'n Proto

# A basic echo service


### A basic echo service

Start by writing a [Cap'n Proto schema file][schema].
For example, here is a very simple echo service:

```capnp
interface Echo {
  ping @0 (msg :Text) -> (reply :Text);
}
```

### Passing Capabilities

```ignorelang

```


### How to generate schema?
```bash
cd echo
capnp compile -I$GOPATH/src/zombiezen.com/go/capnproto2/std -ogo * -ocapnp
```