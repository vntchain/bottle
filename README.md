# Bottle

``Bottle``是开发``VNT``智能合约的命令行工具。
``Bottle``支持将c语言智能合约编译成``wasm``，提取``abi``文件，并将``wasm``和``abi``压缩及编码成``VNT``网络合约部署所需要的智能合约文件。

### 编译得到bottle命令

```
git clone git@github.com:vntchain/bottle.git
cd bottle
make all
```

## 使用

```
NAME:
   bottle - the bottle command line interface

   Copyright 2018-2019 The bottle Authors

USAGE:
   bottle [global options] command [command options] [arguments...]
   
VERSION:
   0.6.0-beta
   
COMMANDS:
   compile     Compile contract code to wasm and compress
   compress    Compress wasm and abi
   decompress  Deompress file into wasm and abi
   hint        Contract hint
   help, h     Shows a list of commands or help for one command
   
COMPILE OPTIONS:
  --code value  Specific a contract code path, - for STDIN
  --include     Specific the head file directory need by contract
  --output      Specific a output directory path
  
COMPRESS OPTIONS:
  --wasm value  Specific a wasm path
  --abi value   Specific a abi path need by contract
  --output      Specific a output directory path
  
DECOMPRESS OPTIONS:
  --file value  Specific a compress file path to decompress
  --output      Specific a output directory path
  
HINT OPTIONS:
  --code value  Specific a contract code path, - for STDIN
  
GLOBAL OPTIONS:
  --help, -h  show help
  

COPYRIGHT:
   Copyright 2018-2019 The bottle Authors
```

## 许可证

所有`bottle`仓库生成的二进制程序都采用GNU General Public License v3.0许可证, 具体请查看[COPYING](https://github.com/vntchain/bottle/blob/master/LICENSE)。
