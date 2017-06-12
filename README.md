# elements-cli

elements scripts or utility or some other clients.

### Quick start

#### docker
```
➜  elements-cli git:(master) docker run -it --rm inhandnet/elements-cli elements-cli fix migrate-online-stats -h
Unable to find image 'inhandnet/elements-cli:latest' locally
latest: Pulling from inhandnet/elements-cli
ef0380f84d05: Already exists
24c170465c65: Already exists
4f38f9d5c3c0: Already exists
d36744f83dc1: Already exists
107ef86d710c: Already exists
b789525fd509: Already exists
633896968756: Already exists
f451e6d11376: Pull complete
Digest: sha256:8c83c2f5afe5c8124e9d9ad2ee594b15884ea77edeef09f9db9678af2ee81f89
Status: Downloaded newer image for inhandnet/elements-cli:latest
Incorrect Usage.

NAME:
   migrate-online-stats - migrate device_oniline_stat to device.online.stats

USAGE:
   command migrate-online-stats [command options] [arguments...]

OPTIONS:
   --url "mongodb://admin:admin@localhost:27017/"	mongodb connect uri

```

#### build and install
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
