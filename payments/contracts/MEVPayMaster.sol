// SPDX-License-Identifier: MIT
pragma solidity >=0.4.22 <0.9.0;

struct Payee {
  address entity;
  uint256 already_paid_out;
}

contract MEVPayMaster is Ownable {
  uint256 private total_ever_received;
  uint256 private total_ever_paid_out;

  event ValidatorPaid(uint256 amount);
  event StakingPoolPaid(uint256 amount);

  Payee public validator;
  Payee public staking_pool;
  Payee public coinbase;

  fallback() external payable {
    require(_msgSender() == owner());
    total_ever_received += msg.value;
  }

  function validator_register(address entity)
    external
    onlyOwner
    returns (bool)
  {
    require(msg.sender == block.coinbase && entity != address(0));
    validator.entity = entity;
    return true;
  }

  function staking_pool_register(address entity)
    external
    onlyOwner
    returns (bool)
  {
    require(msg.sender == block.coinbase && entity != address(0));
    staking_pool.entity = entity;
    return true;
  }

  function coinbase_cut(uint256 gross_owed) private {
    uint256 net_owed = gross_owed - coinbase.already_paid_out;
    require(_amount < net_owed);
    coinbase.already_paid_out += _amount;
    address(this).transfer(block.coinbase, _amount);
  }

  function validator_collect(uint256 _amount) external {
    require(msg.sender == validator.entity);
    uint256 gross_owed = (100 * (total_ever_received) * 33) / 100;
    uint256 net_owed = gross_owed - validator.already_paid_out;
    require(_amount < net_owed);
    validator.already_paid_out += _amount;
    address(this).transfer(msg.sender, _amount);
    coinbase_cut(gross_owed);
    emit ValidatorPaid(_amount);
  }

  function staking_pool_collect(uint256 _amount) external {
    require(msg.sender == staking_pool.entity);
    uint256 gross_owed = (100 * (total_ever_received) * 33) / 100;
    uint256 net_owed = gross_owed - staking_pool.already_paid_out;
    require(_amount < net_owed);
    staking_pool.already_paid_out += _amount;
    address(this).transfer(msg.sender, _amount);
    coinbase_cut(gross_owed);
    emit StakingPoolPaid(_amount);
  }
}