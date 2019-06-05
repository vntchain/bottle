#include "./vntlib.h"

KEY address ens;
KEY string rootNode;

CALL void setSubnodeOwner(CallParams params, string node, string label, address owner);
CALL address owner(CallParams params, string node);


void keyhorarqeg(){
AddKeyInfo( &ens, 7, &ens, 9, false);
AddKeyInfo( &rootNode, 6, &rootNode, 9, false);
}
constructor FIFSRegistrar(address ensAddr, string node)
{
keyhorarqeg();
InitializeVariables();
    ens = ensAddr;
    rootNode = node;
}

void onlyOwner(string subnode)
{
    string node = SHA3(Concat(rootNode, subnode));
    CallParams params = {ens, U256(0), 100000};
    address currentOwner = owner(params, node);

    if (!Equal(currentOwner, Address("0x0")) && !Equal(currentOwner, GetSender()))
    {
        Revert("need owner");
    }
}

void registernode(string subnode, address owner)
{
    onlyOwner(subnode);
    CallParams params = {ens, U256(0), 100000};
    setSubnodeOwner(params, rootNode, subnode, owner);
}
