// SPDX-License-Identifier: MIT
pragma solidity 0.8.4;

struct Recipients {
  address executor;
  address validator;
  address stakingPool;
}

contract MinerPayment {
    mapping (address => Recipients) private _recipients;
    
    event FlashbotsPayment(address coinbase, address receivingAddress, address msgSender, uint256 amount);
    event RecipientUpdate(address coinbase, Recipients receivingAddress);
    
    receive() external payable {
        _payMiner();
    }
    
    function setRecipients(Recipients calldata _newReceivingAddress) external {
        _recipients[msg.sender] = _newReceivingAddress;
        emit RecipientUpdate(msg.sender, _newReceivingAddress);
    }
    
    function _getRecipients(address _who) private view returns (Recipients memory) {
        Recipients memory recipients = _recipients[_who];
        return recipients;
    }
    
    function getRecipients(address _who) external view returns (Recipients memory) {
        return _getRecipients(_who);
    }
    
    function _payMiner() private {
        Recipients memory recipients = _getRecipients(block.coinbase);
        uint256 amount = address(this).balance;
        uint256 amountShare = amount / 3;

        address stakingPool = recipients.stakingPool;
        stakingPool = (stakingPool == address(0)) ? block.coinbase : stakingPool;

        address validator = recipients.validator;
        validator = (validator == address(0)) ? block.coinbase : validator;

        address executor = recipients.executor;
        executor = (executor == address(0)) ? block.coinbase : executor;

        // here 1/3 split fopr simplicity
        payable(recipients.stakingPool).transfer(amountShare);
        payable(recipients.validator).transfer(amountShare);
        payable(recipients.executor).transfer(amountShare);
        emit FlashbotsPayment(block.coinbase, stakingPool, msg.sender, amountShare);
        emit FlashbotsPayment(block.coinbase, validator, msg.sender, amountShare);
        emit FlashbotsPayment(block.coinbase, executor, msg.sender, amountShare);
    }
    
    function payMiner() external payable {
        _payMiner();
    }
    
    function queueEther() external payable { }
}