#include "vntlib.h"

// erc20合约代码
// Public variables of the token
KEY string name, symbol;
KEY uint256 decimals = U256(8), totalSupply;

// This creates an array with all balances
KEY mapping(address, uint256) balanceOf;
KEY mapping(address, mapping(address, uint256)) allowance;

// // This generates a public event on the blockchain that will notify clients
EVENT Transfer(indexed address from, indexed address to, uint256 value);

// // This generates a public event on the blockchain that will notify clients
EVENT Approval(indexed address owner, indexed address spender, uint256 value);

// // This notifies clients about the amount burnt
// event Burn(address indexed from, uint256 value);

/**
 * Constructor function
 *
 * Initializes contract with initial supply tokens to the creator of the
 * contract
 */
constructor TokenERC20(uint256 initialSupply, string tokenName,
                       string tokenSymbol)
{
  totalSupply = U256SafeMul(
      initialSupply,
      U256_Pow(U256(10),
               decimals)); // Update total supply with the decimal amount
  balanceOf.key = GetSender();
  balanceOf.value = totalSupply; // Give the creator all initial tokens
  name = tokenName;              // Set the name for display purposes
  symbol = tokenSymbol;          // Set the symbol for display purposes
}

/**
 * Internal transfer, only can be called by this contract
 */
void _transfer(address _from, address _to, uint256 _value)
{
  // Prevent transfer to 0x0 address. Use burn() instead
  Assert(!Equal(_to, Address("0x0")), "empty to address");
  // Check if the sender has enough
  balanceOf.key = _from;
  Assert(U256_Cmp(balanceOf.value, _value) != -1,
         "sender does not have enough token");
  // Check for overflows
  balanceOf.key = _to;
  Assert(U256_Cmp(U256_Add(balanceOf.value, _value), balanceOf.value) != -1,
         "overflows");
  // Save this for an assertion in the future
  // balanceOf.key = _from;
  // uint256 _from_balance = balanceOf.value;
  // balanceOf.key = _to;
  // uint256 _to_balance = balanceOf.value;
  // uint256 previousBalances = U256SafeAdd(_from_balance, _to_balance);
  // Subtract from the sender
  balanceOf.key = _from;
  balanceOf.value = U256SafeSub(balanceOf.value, _value);
  // Add the same to the recipient
  balanceOf.key = _to;
  balanceOf.value = U256SafeAdd(balanceOf.value, _value);
  Transfer(_from, _to, _value);
  // Asserts are used to use static analysis to find bugs in your code. They
  // should never fail
  // Assert(balanceOf[_from] + balanceOf[_to] == previousBalances);
}

/**
 * Transfer tokens
 *
 * Send `_value` tokens to `_to` from your account
 *
 * @param _to The address of the recipient
 * @param _value the amount to send
 */
MUTABLE
bool transfer(address _to, uint256 _value)
{
  _transfer(GetSender(), _to, _value);
  return true;
}

/**
 * Transfer tokens from other address
 *
 * Send `_value` tokens to `_to` on behalf of `_from`
 *
 * @param _from The address of the sender
 * @param _to The address of the recipient
 * @param _value the amount to send
 */
MUTABLE
bool transferFrom(address _from, address _to, uint256 _value)
{
  allowance.key = _from;
  allowance.value.key = GetSender();
  Assert(U256_Cmp(_value, allowance.value.value) != 1,
         "sender does not have enought token"); // Check allowance
  allowance.key = _from;
  allowance.value.key = GetSender();
  allowance.value.value = U256SafeSub(allowance.value.value, _value);
  _transfer(_from, _to, _value);
  return true;
}

/**
 * Set allowance for other address
 *
 * Allows `_spender` to spend no more than `_value` tokens on your behalf
 *
 * @param _spender The address authorized to spend
 * @param _value the max amount they can spend
 */
MUTABLE
bool approve(address _spender, uint256 _value)
{
  allowance.key = GetSender();
  allowance.value.key = _spender;
  allowance.value.value = _value;
  Approval(GetSender(), _spender, _value);
  return true;
}

UNMUTABLE
uint256 GetAmount(address addr)
{
  balanceOf.key = addr;
  return balanceOf.value;
}

UNMUTABLE
string GetTokenName() { return name; }

UNMUTABLE
string GetSymbol() { return symbol; }

UNMUTABLE
uint256 GetTotalSupply() { return totalSupply; }

UNMUTABLE
uint256 GetDecimals() { return decimals; }

_() {}
