# Bottle

``Bottle``是开发``VNT``智能合约的命令行工具。
``Bottle``支持将c语言智能合约编译成``wasm``，提取``abi``文件，并将``wasm``和``abi``压缩及编码成``VNT``网络合约部署所需要的智能合约文件。

## 编译与安装

### 安装依赖clang 5.0 llvm 5.0

#### ``ubuntu``

```
wget https://raw.githubusercontent.com/go-clang/gen/master/scripts/switch-clang-version.sh

/bin/bash ./switch-clang-version.sh 5.0

//使用绝对路径进行软链接
sudo ln -s /usr/lib/llvm-5.0/lib/libclang*so /usr/lib/

```

#### ``mac``

```
brew install llvm@5

sudo ln -s /usr/local/opt/llvm@5/lib/libclang*dylib /usr/local/lib

```
#### ``centos``

##### 添加yum源

```
[alonid-llvm-5.0.0]
name=Copr repo for llvm-5.0.0 owned by alonid
baseurl=https://copr-be.cloud.fedoraproject.org/results/alonid/llvm-5.0.0/epel-7-$basearch/
type=rpm-md
skip_if_unavailable=True
gpgcheck=1
gpgkey=https://copr-be.cloud.fedoraproject.org/results/alonid/llvm-5.0.0/pubkey.gpg
repo_gpgcheck=0
enabled=1
enabled_metadata=1
```

##### 更新源
``yum makecache``

##### 安装llvm
``yum install llvm``


### 安装依赖wasmception

#### ```mac```

[wasmception] https://github.com/ooozws/clang-heroku-slug/blob/master/precomp/wasmception-darwin-bin.tar.gz

#### ```linux```

[wasmception] https://github.com/ooozws/clang-heroku-slug/blob/master/precomp/wasmception-linux-bin.tar.gz

下载wasmception并解压wasmception,设置wasmception的环境变量

```
echo export VNT_WASMCEPTION="/[PATH]/wasmception-[XXX]-bin" >> ~/.bash_profile
source ~/.bash_profile
``` 

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
