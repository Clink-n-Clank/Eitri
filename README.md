# Eitri

> Eitri Command-line toolsets that help maintain code base and be productive.

## Install

```bash
go install github.com/Clink-n-Clank/eitri
```

## Use
You may need to install depended on tool since some Eitri commands build on top of that.

To install dependent on tools just call: `eitri tools` 

```
Usage:
  eitri [command]

Available Commands:
  build       Generate binary: Generates Google Wire DI code and compile it to binary
  completion  Generate the autocompletion script for the specified shell
  envoy       Build Envoy: transcoder: gRPC-JSON transcoder in Base 64
  help        Help about any command
  mock        Generate Mock: Generated mock code from config file: gomockhandler.json
  tools       Install tools: Will install depended tools for work

Flags:
  -h, --help   help for eitri

Use "eitri [command] --help" for more information about a command.
```