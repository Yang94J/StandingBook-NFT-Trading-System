pragma solidity ^0.6.0;

abstract contract IERC721Receiver {
    function onERC721Received(address operator, address from, uint256 tokenId, bytes memory data) public virtual returns (bytes4);
}