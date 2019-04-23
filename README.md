# Bottle

``Bottle``是开发``VNT``智能合约的命令行工具。
``Bottle``支持将c语言智能合约编译成``wasm``，提取``abi``文件，并将``wasm``和``abi``压缩及编码成``VNT``网络合约部署所需要的智能合约文件。

## 编译得到``bottle``命令

### Mac OS编译

将``bottle``代码clone到选择的目录

```
git clone https://github.com/vntchain/bottle
```

编译``bottle``需要[go](https://golang.org/)编译器

```
brew install go
```

然后使用以下命令编译得到``bottle``

```
cd bottle
make bottle
```

最后使用以下命令运行bottle

```
./build/bin/bottle
```

### Ubuntu 14.04及16.06编译

将``bottle``代码clone到选择的目录

```
git clone https://github.com/vntchain/bottle
```

编译``bottle``需要[go](https://golang.org/)编译器，``go``的安装请参考``go``的官方文档


然后使用以下命令编译得到``bottle``

```
cd bottle
make bottle
```

最后使用以下命令运行bottle

```
./build/bin/bottle
```

### 其他系统

目前``bottle``暂不支持除上述系统之外的系统，如果希望在不支持的系统上运行bottle，请使用``docker``的方式

## 使用docker运行

### 编译得到docker镜像

将``bottle``代码clone到选择的目录

```
git clone https://github.com/vntchain/bottle
```


然后使用以下命令编译得到镜像

```
cd bottle
make bottle-docker
```

### 通过dockerhub得到镜像

镜像拉取命名

```
docker pull vntchain/bottle:0.6.0
```

## 使用

### 编译智能合约代码

通过编译得到``bottle``

```
bottle compile -code <your contract path> -output <the path of your choosing to save the compiled contract file>
```

使用``docker``运行

```
docker run --rm -v <your contract directory>:/tmp vntchain/bottle:0.6.0 compile -code /tmp/<your contract file name> 
```

通过以上命令可以获得后缀名为``compress``的编译文件和后缀名为``abi``的abi文件，使用compress文件可以部署智能合约到``VNT``网络，通过abi文件可以访问``VNT``网络中的智能合约

### 智能合约代码纠错及提示

通过编译得到``bottle``

```
bottle hint -code <your contract path>
```

使用``docker``运行

```
docker run --rm -v <your contract directory>:/tmp vntchain/bottle:0.6.0  hint -code /tmp/<your contract file name> 
```

### Bottle命令
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
  --code value  Specific a contract code path
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
  --code value  Specific a contract code path
  
GLOBAL OPTIONS:
  --help, -h  show help
  

COPYRIGHT:
   Copyright 2018-2019 The bottle Authors
```

## 许可证

所有`bottle`仓库生成的二进制程序都采用GNU General Public License v3.0许可证, 具体请查看[COPYING](https://github.com/vntchain/bottle/blob/master/LICENSE)。
