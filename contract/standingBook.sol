pragma solidity ^0.6.0;

// -------------------------------------------------------------------------------------------------------------------
import "./openzeppelin/Ownable.sol";
import "./openzeppelin/Roles.sol";

import "./ierc721.sol";
import "./erc165.sol";
import "./ierc721receiver.sol";




// ---StandingBook contract (main contract)--------------------------------------------------------------------------------------------------------

contract StandingBook is ERC165, IERC721, Ownable {

    // ------variables-----------------------------------------------------------------------------------------------------------------------------

    // Use this constant to illustrate the payable(NFT Token)
    bytes4 private _ERC721_RECEIVED = 0x150b7a02;
    // Use this constant to illustrate the existance of all functions required by the IERC721
    bytes4 private constant _INTERFACE_ID_ERC721 = 0x80ac58cd;

    mapping (uint256 => address) private _tokenOwner;
    mapping (uint256 => address) private _tokenApprovals;
    mapping (address => uint256) private _ownedTokenCount;
    mapping (address => mapping (address => bool)) private _operatorApprovals;

    using Roles for Roles.Role;
    Roles.Role private _minters;
    // ------constructor function -----------------------------------------------------------------------------------------------------------
    constructor() public {
        _registerInterface(_INTERFACE_ID_ERC721);
    }


    // -----ERC721 modifier Implementation-------------------------------------------------------------------------------------------------
    modifier _ifTokenExists(uint256 tokenId) {
        require(_tokenOwner[tokenId] != address(0x0),"Invalid Standing_book Id");
        _;
    }

    modifier _ifTokenOccupied(uint256 tokenId) {
        require(_tokenOwner[tokenId] == address(0x0),"Current Token Already Exists");
        _;
    }

    modifier _isUserValid(address user){
        require(user != address(0x0),"Invalid User Address");
        _;
    }

    modifier _isActualUserOfToken(address user, uint256 tokenId){
        require (ownerOf(tokenId) == user,"User Is Not The Actual Owner Of the Standing_book");
        _;
    }

    modifier _senderIsActualUserOfToken(uint256 tokenId){
        require (ownerOf(tokenId) == msg.sender,"User Is Not The Actual Owner Of the Standing_book");
        _;
    }

    modifier _isSelfOperation(address to){
        require (msg.sender != to, "Already Acquired");
        _;
    }

    modifier _hasRightOverToken(uint256 tokenId){
        require ((ownerOf(tokenId) == msg.sender) || (getApproved(tokenId)!=address(0x0)) || (isApprovedForAll(ownerOf(tokenId),msg.sender)),"User Has No Rights Over The Token");
        _;
    }

    // -----RBAC modifier Implementation-------------------------------------------------------------------------------------------------
    modifier _isMinter(){
        require(_minters.has(msg.sender),"Minter Denied, No Rights.");
        _;
    }


    // -----ERC721 function Implementation-------------------------------------------------------------------------------------------------
    function ownerOf(uint256 tokenId) public view override _ifTokenExists(tokenId) returns (address) {
        return _tokenOwner[tokenId];
    }

    function balanceOf(address owner) public view override returns (uint256) {
        return _ownedTokenCount[owner];
    }

    function _approve(address to, uint256 tokenId) internal {
        _tokenApprovals[tokenId] = to;
        emit Approval (ownerOf(tokenId), to, tokenId);
    }

    function approve(address to, uint256 tokenId) public override _senderIsActualUserOfToken(tokenId) _isSelfOperation(to) {
        _approve(to,tokenId);
    }


    function getApproved(uint256 tokenId) public override view _ifTokenExists(tokenId) returns (address){
        return _tokenApprovals[tokenId];
    }


    function setApprovalForAll(address operator, bool _approved) public override _isSelfOperation(operator){
        _operatorApprovals[msg.sender][operator] = _approved;
        emit ApprovalForAll(msg.sender, operator, _approved);
    }

    function isApprovedForAll(address owner, address operator) public view  override _isUserValid(owner) returns (bool) {
        return _operatorApprovals[owner][operator];
    }

    function _transferFrom(address from, address to, uint256 tokenId) internal virtual {
        _approve(address(0x0),tokenId);
        _ownedTokenCount[to]++;
        _ownedTokenCount[from]--;
        _tokenOwner[tokenId] = to;
        emit Transfer(from,to,tokenId);
    }

    function transferFrom(address from, address to, uint256 tokenId) public _isActualUserOfToken(from,tokenId) {
        _transferFrom(from,to,tokenId);
    }

    function isContract(address addr) internal view returns (bool) {
        uint256 size;
        assembly { size := extcodesize(addr)}
        return size > 0;
    }

    function _checkOnERC721Received(address from, address to, uint256 tokenId, bytes memory _data) private returns (bool) {
        if (!isContract(to)){
            return true;
        }
        (bool success, bytes memory returndata) = to.call(abi.encodeWithSelector(IERC721Receiver(to).onERC721Received.selector,msg.sender,from,tokenId,_data));
        if (!success) {
            revert("ERC721: transfer to non ERC721Receiver implementer");
        } else {
            bytes4 retval = abi.decode(returndata,(bytes4));
            return (retval == _ERC721_RECEIVED);
        }
    }

    function _safeTransferFrom(address from, address to, uint256 tokenId, bytes memory data) internal virtual {
        _transferFrom(from,to,tokenId);
        require(_checkOnERC721Received(from,to,tokenId,data),"The Target Address DOESNT Support Current Token");
    }

    function safeTransferFrom(address from, address to, uint256 tokenId, bytes memory data) public virtual override _hasRightOverToken(tokenId) _isActualUserOfToken(from,tokenId){
        _safeTransferFrom(from,to,tokenId,data);
    }

    function safeTransferFrom(address from, address to, uint256 tokenId) public virtual override _hasRightOverToken(tokenId) _isActualUserOfToken(from, tokenId){
        _safeTransferFrom(from,to,tokenId,"");
    }

    // -----RBAC Implementation-------------------------------------------------------------------------------------------------
    function addMinter(address minter) public onlyOwner {
        _minters.add(minter);
    }

    function removeMinter(address minter) public onlyOwner{
        _minters.remove(minter);
    }

    function isMinter() public view returns (bool){
        return _minters.has(msg.sender);
    }

    function _mint(address to, uint256 tokenId) internal virtual _isUserValid(to) _ifTokenOccupied(tokenId){
        _tokenOwner[tokenId] = to;
        _ownedTokenCount[to] ++;
        emit Transfer(address(0x0),to,tokenId);
    }

    function uploadAndMint(address to, uint256 tokenId) external _isMinter(){
        _mint(to,tokenId);
    }

}