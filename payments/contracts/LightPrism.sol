// SPDX-License-Identifier: MIT
pragma solidity 0.8.4;

struct Recipients {
  address executor;
  address stakingPool;
}

contract LightPrism {
  mapping(address => Recipients) private _recipients;

  event FlashbotsPayment(
    address coinbase,
    address receivingAddress,
    address msgSender,
    uint256 amount
  );

  event RecipientUpdate(
    address coinbase,
    address executor,
    address stakingPool
  );

  receive() external payable {
    _payMiner();
  }

  function setRecipients(address executor, address stakingPool) external {
    // just simplify for now <- this should only be valid of msg.sender is coinbase
    _recipients[block.coinbase] = Recipients(executor, stakingPool);
    emit RecipientUpdate(msg.sender, executor, stakingPool);
  }

  function _getRecipients(address _who)
    private
    view
    returns (Recipients memory)
  {
    Recipients memory recipients = _recipients[_who];
    return recipients;
  }

  function getRecipients(address _who)
    external
    view
    returns (Recipients memory)
  {
    return _getRecipients(_who);
  }

  function _payMiner() private {
    Recipients memory recipients = _getRecipients(block.coinbase);
    uint256 amount = address(this).balance;
    uint256 poolShare = (amount * 2) / 3;
    uint256 executorShare = amount / 3;

    address stakingPool = recipients.stakingPool;
    stakingPool = (stakingPool == address(0)) ? block.coinbase : stakingPool;

    address executor = recipients.executor;
    executor = (executor == address(0)) ? block.coinbase : executor;

    // here 2/3 and 1/3 split for simplicity
    payable(recipients.stakingPool).transfer(poolShare);
    payable(recipients.executor).transfer(executorShare);
    emit FlashbotsPayment(block.coinbase, stakingPool, msg.sender, poolShare);
    emit FlashbotsPayment(block.coinbase, executor, msg.sender, executorShare);
  }

  function payMiner() external payable {
    _payMiner();
  }

  function queueEther() external payable {}
}
