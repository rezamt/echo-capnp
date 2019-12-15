using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("echo");
$Go.import("github.com/rezamt/echo");

interface Echo {
  ping @0 (msg :Text) -> (reply :Text);
}
