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

同时需要安装依赖包libxml2-dev

```
sudo apt-get install -y libxml2-dev
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

## 使用

### 编译智能合约代码

通过编译得到``bottle``

```
bottle compile -code <your contract path> -output <the path of your choosing to save the compiled contract file>
```

使用``docker``运行

```
docker run --rm -v <your contract directory>:/tmp bottle:0.6.0 compile -code /tmp/<your contract file name> 
```

通过以上命令可以获得后缀名为``compress``的编译文件和后缀名为``abi``的abi文件，使用compress文件可以部署智能合约到``VNT``网络，通过abi文件可以访问``VNT``网络中的智能合约

### 智能合约代码纠错及提示

通过编译得到``bottle``

```
bottle hint -code <your contract path>
```

使用``docker``运行

```
docker run --rm -v <your contract directory>:/tmp bottle:0.6.0  hint -code /tmp/<your contract file name> 
```

### 合约创建与部署

通过以下命令在空文件夹中创建智能合约与部署配置

```
bottle init
```

生成的文件结构

```
|____migrations
| |____1_initial_migration.js
|____contracts
| |____Migrations.c
| |____vntlib.h
| |____Erc20.c
|____bottle.js
```

创建完成后可以通过``bottle build``,``bottle migrate``来编译智能合约并将他们部署到``VNTChain``网络.


#### Build命令

``bottle build``命令需要在``bottle.js``所在的目录下执行，会把``contracts``文件夹下的智能合约文件进行编译并输出到``build/contracts``文件夹下

#### Migrate命令

``bottle migrate``命令需要在``bottle.js``所在的目录下执行，用于执行``migrations``文件夹下的``js``文件，``js``文件用于将``contracts``文件夹下的智能合约部署到``VNTChain``网络

#### Migrate文件

``migrate``文件是位于``migrations``文件夹下的``js``文件，一个简单的``migrate``文件如下所示

文件名: 4_example_migration.js

```js
var MyContract = artifacts.require("./contracts/MyContract.c");

module.exports = function(deployer) {
  // deployment steps
  deployer.deploy(MyContract);
};
```

``migrate``文件名必须以数字为前缀，后缀为描述。数字前缀用于按顺序执行``migrate``文件，以及记录文件是否已被执行，已被执行的``migrate``文件会被再一次运行的``bottle migrate``命令所忽略，如果想要重新执行之前的``migrate``文件，请参考``migrate``命令的参数。后缀用于描述``migrate``文件，方便识别和理解文件的作用。

##### artifacts.require(<contract_path>)

``migrate``文件通过``artifacts.require``来引入需要操作的智能合约，参数是``contracts``文件夹内的智能合约的绝对地址或者相对地址，通过引用，将返回一个智能合约对象，智能合约对象用于智能合约部署以及访问智能合约的方法

假设在``contracts``文件夹下有如下智能合约文件

智能合约1: `./contracts/Contract1.c`
智能合约2: `./contracts/Contract2.c`

为了与以上两个智能合约交互，``artifacts.require``可以这样使用

```js
var contract1 = artifacts.require("./contracts/Contract1.c");
var contract2 = artifacts.require("./contracts/Contract2.c");
```

##### module.exports

所有的``migrate``文件都必须通过``module.exports``语法导出函数，导出的函数的第一个参数为``deployer``，``deployer``对象提供了部署和访问智能合约的方法，第二个参数为``network``，``network``用于通过不同的网络部署智能合约，参考以下例子：

不使用``network``

```js
module.exports = function(deployer) {

}
```

使用``network``

```js
module.exports = function(deployer, network) {
  if (network == "live") {
    // Do something specific to the network named "live".
  } else {
    // Perform a different step otherwise.
  }
}
```

指定``network``通过``bottle migrate``的``network``参数

```
bottle migrate --network xxx
```

#### 初始化migrate智能合约

``bottle``需要有一个``migrate``智能合约才能使用``bottle migrate``功能，该智能合约包含特定的接口，会在第一次执行``botlte migrate``时部署，此后将不会更新。在使用`bottle init`创建新项目时，会默认创建该智能合约。

文件名: `contracts/Migrations.c`

```c
#include "vntlib.h"

KEY address owner;
KEY uint32 last_completed_migration;
constructor Migrations()
{
  owner = GetSender();
}

void onlyOwner()
{
  Require(Equal(owner, GetSender()), "is not owner");
}

MUTABLE
void setCompleted(uint32 completed)
{
  onlyOwner();
  last_completed_migration = completed;
}

UNMUTABLE
uint32 get_last_completed_migration()
{
  return last_completed_migration;
}

```

``migrate``智能合约必须在第一次执行``bottle migrate``的时候进行部署，因此，需要创建如下的``migrate``文件

文件名: `migrations/1_initial_migration.js`

```javascript
var Migrations = artifacts.require("../contracts/Migrations.c");

module.exports = function (deployer) {
  // Deploy the Migrations contract as our only task
  deployer.deploy(Migrations)
};
```

之后，可以增加编号前缀来创建新的``migrate``文件，以部署其他智能合约。

`bottle init`创建新项目时，会默认创建该``migrate``文件。

#### Deployer对象

``migrate``文件需要使用``deployer``对象来执行部署任务，同时可以同步编写部署任务，它们将以正确的顺序执行：

```js
// Stage deploying A before B
deployer.deploy(A);
deployer.deploy(B);
```

或者，部署程序上的每个函数都可以用作Promise，按顺序依赖于上一个任务执行的部署任务：

```js
// Deploy A, then deploy B, passing in A's newly deployed address
deployer.deploy(A).then(function() {
  return deployer.deploy(B, A.address);
});
```

#### Deployer API

``deployer``对象提供了方法用于简化智能合约的部署。
##### deployer.deploy(contract, args..., options)

参数``contract``为使用``artifacts.require``引用的智能合约对象。
参数``args...``为智能合约的构造函数的参数，用于初始化智能合约。
参数``options``用于指定``from``，``gas``及``overwrite``等信息，``overwrite``用于重新部署某个已经完成部署的智能合约，默认的``options``参数在``bottle.js``文件中配置

例子:

```js
// Deploy a single contract without constructor arguments
deployer.deploy(A);

// Deploy a single contract with constructor arguments
deployer.deploy(A, arg1, arg2, ...);

// Don't deploy this contract if it has already been deployed
deployer.deploy(A, {overwrite: false});

// Set a maximum amount of gas and `from` address for the deployment
deployer.deploy(A, {gas: 4612388, from: "0x...."});

// External dependency example:
//
// For this example, our dependency provides an address when we're deploying to the
// live network, but not for any other networks like testing and development.
// When we're deploying to the live network we want it to use that address, but in
// testing and development we need to deploy a version of our own. Instead of writing
// a bunch of conditionals, we can simply use the `overwrite` key.
deployer.deploy(SomeDependency, {overwrite: false});
```

##### deployer.then(function() {...})

通过promise对象可以运行任意的部署步骤并调用指定的智能合约内部方法来进行交互

例子:

```js
var ERC20 = artifacts.require("../contracts/Erc20.c")

module.exports = function (deployer, a) {
    deployer.deploy(ERC20, "1000000", "bitcoin", "BTC").then(function (instance) {
        deploy = instance;
        return deploy.GetTotalSupply()
    }).then(function (totalSupply) {
        console.log("totalSupply", totalSupply.toString());
        return deploy.GetDecimals();
    }).then(function (decimals) {
        console.log("decimals", decimals.toString());
        return deploy.GetTokenName();
    }).then(function (tokenName) {
        console.log("tokenName", tokenName);
        return deploy.GetAmount("0x122369f04f32269598789998de33e3d56e2c507a")
    }).then(function (balance) {
        console.log("balance", balance.toString());
    })
};
```

### Bottle命令
```
NAME:
   bottle - the bottle command line interface

   Copyright 2018-2019 The bottle Authors

USAGE:
   bottle [global options] command [command options] [arguments11...]
   
VERSION:
   0.6.0-beta-1e52fa2f
   
COMMANDS:
   build       Build contracts
   compile     Compile contract source file
   compress    Compress wasm and abi file
   decompress  Deompress file into wasm and abi file
   hint        Contract hint
   init        Initialize dapp project
   migrate     Run migrations to deploy contracts
   help, h     Shows a list of commands or help for one command
   
PATH OPTIONS:
  --code value  Specific a contract code path
  --include     Specific the head file directory need by contract
  --output      Specific a output directory path
  --file value  Specific a compress file path to decompress
  --abi value   Specific a abi path needed by contract
  --wasm value  Specific a wasm path
  
MIGRATE OPTIONS:
  --reset          Run all migrations from the beginning, instead of running from the last completed migration
  -f value         Run contracts from a specific migration. The number refers to the prefix of the migration file (default: 0)
  -t value         Run contracts to a specific migration. The number refers to the prefix of the migration file (default: 0)
  --network value  Specify the network to use, saving artifacts specific to that network. Network name must exist in the configuration
  --verbose-rpc    Log communication between bottle and the VNTChain client
  
GLOBAL OPTIONS:
  --help, -h  show help
  

COPYRIGHT:
   Copyright 2018-2019 The bottle Authors
```

## 许可证

所有`bottle`仓库生成的二进制程序都采用GNU General Public License v3.0许可证, 具体请查看[COPYING](https://github.com/vntchain/bottle/blob/master/LICENSE)。
