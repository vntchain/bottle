// pragma solidity >=0.4.21 <0.6.0;

// contract Migrations {
//   address public owner;
//   uint public last_completed_migration;

//   constructor() public {
//     owner = msg.sender;
//   }

//   modifier restricted() {
//     if (msg.sender == owner) _;
//   }

//   function setCompleted(uint completed) public restricted {
//     last_completed_migration = completed;
//   }

//   function upgrade(address new_address) public restricted {
//     Migrations upgraded = Migrations(new_address);
//     upgraded.setCompleted(last_completed_migration);
//   }
// }

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
