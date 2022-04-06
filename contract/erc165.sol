pragma solidity ^0.6.0;

import "./ierc165.sol";

contract ERC165 is IERC165 {
bytes4 private  _INTERFACE_ID_ERC165;
mapping(bytes4 => bool) private _supportedInterfaces;

constructor () internal {
/*
   This is not necessarily true... as the constant should be sth related to the function
*/
_INTERFACE_ID_ERC165 = bytes4(keccak256(abi.encode(now,msg.sender)));
_registerInterface(_INTERFACE_ID_ERC165);
}

function _registerInterface(bytes4 interfaceId) internal virtual {
require(interfaceId != 0xffffffff, "ERC165 : invalid Interface id");
_supportedInterfaces[interfaceId] = true;
}

function supportsInterface(bytes4 interfaceId) external view override returns (bool) {
return _supportedInterfaces[interfaceId];
}
}