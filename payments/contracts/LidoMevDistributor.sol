// SPDX-License-Identifier: MIT

pragma solidity 0.8.4;


interface INodeOperatorsRegistry {
    function getRewardsDistribution(uint256 _totalRewardShares)
        external view returns (
            address[] memory recipients,
            uint256[] memory shares
        );
}


interface ILido {
    function totalSupply() external view returns (uint256);
    function getTotalShares() external view returns (uint256);
    function submit(address _referral) external payable returns (uint256 sharesMinted);
    function sharesOf(address _account) external view returns (uint256);
    function getPooledEthByShares(uint256 _sharesAmount) external view returns (uint256);
    function burnShares(address _account, uint256 _amount) external returns (uint256 newTotalShares);
    function transfer(address _recipient, uint256 _amount) external returns (bool);
    function getOperators() external view returns (INodeOperatorsRegistry);
}


contract LidoMevDistributor {
    event LidoMevReceived(uint256 amount);
    event LidoMevDistributed(uint256 amount);

    address public lidoAddress;
    uint256 public validatorsMevShare;

    constructor(address _lidoAddress, uint256 _validatorsMevShare) {
        require(_validatorsMevShare <= 10**18);
        lidoAddress = _lidoAddress;
        validatorsMevShare = _validatorsMevShare;
    }

    receive() external payable {
        emit LidoMevReceived(msg.value);
    }

    function distribureMev() external payable {
        if (msg.value > 0) {
            emit LidoMevReceived(msg.value);
        }

        uint256 totalEth = address(this).balance;
        require(totalEth > 0);

        ILido lido = ILido(lidoAddress);
        lido.submit{value: totalEth}(address(this));

        uint256 stEthTotalSupply = lido.totalSupply();
        uint256 prevTotalShares = lido.getTotalShares();

        uint256 sharesToDistribute = lido.sharesOf(address(this));
        uint256 stEthToDistribute = sharesToDistribute * stEthTotalSupply / prevTotalShares;
        uint256 validatorsStEth = (stEthToDistribute * validatorsMevShare) / 10**18;

        // Since we're burning part of the shares, each non-burnt share will become more expensive,
        // including the shares we're going to transfer to validators:
        //
        // newStEthByShare = stEthTotalSupply / (prevTotalShares - stakersShares)
        // validatorsStEth = validatorsShares * newStEthByShare
        // validatorsShares + stakersShares = sharesToDistribute
        //
        // We need to account for this in order to distribure the received stETH between validators
        // and stakers in a given proportion:
        //
        // validatorsShares * newStEthByShare = validatorsStEth
        // validatorsShares * stEthTotalSupply / (prevTotalShares - sharesToDistribute + validatorsShares) = validatorsStEth
        // validatorsShares * stEthTotalSupply = validatorsStEth * (prevTotalShares - sharesToDistribute + validatorsShares)
        // validatorsShares * (stEthTotalSupply - validatorsStEth) = validatorsStEth * (prevTotalShares - sharesToDistribute)
        // validatorsShares = validatorsStEth * (prevTotalShares - sharesToDistribute) / (stEthTotalSupply - validatorsStEth)

        uint256 validatorsShares = validatorsStEth *
            (prevTotalShares - sharesToDistribute) /
            (stEthTotalSupply - validatorsStEth);

        uint256 stakersShares = sharesToDistribute - validatorsShares;

        _distributeValidatorsMev(validatorsShares, lido);
        _distributeStakersMev(stakersShares, lido);

        emit LidoMevDistributed(totalEth);
    }

    function _distributeValidatorsMev(uint256 sharesAmount, ILido lido) internal {
        uint256 stEthAmount = lido.getPooledEthByShares(sharesAmount);

        (address[] memory recipients, uint256[] memory amounts) =
            lido.getOperators().getRewardsDistribution(stEthAmount);

        assert(recipients.length == amounts.length);

        for (uint256 i = 0; i < recipients.length; ++i) {
            lido.transfer(recipients[i], amounts[i]);
        }
    }

    function _distributeStakersMev(uint256 sharesAmount, ILido lido) internal {
        // Burn the pool share. This would increase all other stakers' shares price,
        // effectively distributing the received ETH between the stakers.
        lido.burnShares(address(this), sharesAmount);
    }
}
