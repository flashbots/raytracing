// SPDX-FileCopyrightText: 2020 Lido <info@lido.fi>

// SPDX-License-Identifier: GPL-3.0

pragma solidity 0.4.24;

import "./interfaces/IDepositContract.sol";

contract DepositContractMock is IDepositContract {
    event Deposit(
        bytes pubkey,
        bytes withdrawal_credentials,
        bytes signature,
        bytes32 deposit_data_root,
        uint256 value
    );

    function deposit(
        bytes /* 48 */ pubkey,
        bytes /* 32 */ withdrawal_credentials,
        bytes /* 96 */ signature,
        bytes32 deposit_data_root
    )
        external
        payable
    {
        emit Deposit(pubkey, withdrawal_credentials, signature, deposit_data_root, msg.value);
    }
}
