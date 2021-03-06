pragma solidity ^0.6.0;

import "./ierc165.sol";

abstract contract IERC721 is IERC165{
    event Transfer (address indexed from, address indexed to, uint256 indexed tokenId);
    event Approval (address indexed owner, address indexed approved, uint256 indexed tokenId);

    event ApprovalForAll (address indexed owner, address indexed operator, bool approved);

    function balanceOf(address owner) public view virtual returns (uint256);
    function ownerOf(uint256 tokenId) public view virtual returns (address);

    function approve(address to, uint256 tokenId) public virtual;
    function getApproved(uint256 tokenId) public view virtual returns (address);
    function setApprovalForAll(address operator, bool _approved) public virtual;
    function isApprovedForAll(address owner, address operator) public view virtual returns (bool);
    function safeTransferFrom(address from, address to, uint256 tokenId) public virtual;
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes memory data) public virtual;

}