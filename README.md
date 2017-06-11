# elements-cli

elements scripts or utility or some other clients.

### Quick start

1. Install go and set `GOPATH` and `$PATH=$PATH:$GOPATH/bin`
2. get package
```shell
go get github.com/inhandnet/elements-cli
```
3. try `elements-cli help`  
```shell
NAME:
   elements util - elements scripts utility.

USAGE:
   elements util [global options] command [command options] [arguments...]

VERSION:
   0.0.0

AUTHOR:
  Author - <unknown@email>

COMMANDS:
   fix		fix mongodb documents
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help		show help
   --version, -v	print the version
```

### support:
  1. migrate device_online_stats to device.online.events
