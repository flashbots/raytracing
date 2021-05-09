package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	miner = "0xd912aecb07e9f4e1ea8e6b4779e7fb6aa1c3e4d8"
	// SPDX-License-Identifier: UNLICENSED
	// pragma solidity ^0.7.0;
	// contract Bribe {
	//     function bribe() payable public {
	//         block.coinbase.transfer(msg.value);
	//     }
	// }

	lidoContract = `0x60806040523480156200001157600080fd5b5060405162000fc038038062000fc08339818101604052810190620000379190620000ca565b670de0b6b3a76400008111156200004d57600080fd5b816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508060018190555050506200017d565b600081519050620000ad8162000149565b92915050565b600081519050620000c48162000163565b92915050565b60008060408385031215620000de57600080fd5b6000620000ee858286016200009c565b92505060206200010185828601620000b3565b9150509250929050565b600062000118826200011f565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b62000154816200010b565b81146200016057600080fd5b50565b6200016e816200013f565b81146200017a57600080fd5b50565b610e33806200018d6000396000f3fe6080604052600436106100385760003560e01c80633c45a3d21461007b57806369571ba114610085578063afa724e8146100b057610076565b36610076577f0a6e8c861fe004b2c9d6d3dd59541d622fb4c6dbb805375df34b31b516c89d403460405161006c9190610ad8565b60405180910390a1005b600080fd5b6100836100db565b005b34801561009157600080fd5b5061009a610449565b6040516100a79190610a94565b60405180910390f35b3480156100bc57600080fd5b506100c561046d565b6040516100d29190610ad8565b60405180910390f35b600034111561011c577f0a6e8c861fe004b2c9d6d3dd59541d622fb4c6dbb805375df34b31b516c89d40346040516101139190610ad8565b60405180910390a15b60004790506000811161012e57600080fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508073ffffffffffffffffffffffffffffffffffffffff1663a1903eab83306040518363ffffffff1660e01b815260040161018e9190610a94565b6020604051808303818588803b1580156101a757600080fd5b505af11580156101bb573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906101e09190610a4d565b5060008173ffffffffffffffffffffffffffffffffffffffff166318160ddd6040518163ffffffff1660e01b815260040160206040518083038186803b15801561022957600080fd5b505afa15801561023d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102619190610a4d565b905060008273ffffffffffffffffffffffffffffffffffffffff1663d5002f2e6040518163ffffffff1660e01b815260040160206040518083038186803b1580156102ab57600080fd5b505afa1580156102bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102e39190610a4d565b905060008373ffffffffffffffffffffffffffffffffffffffff1663f5eb42dc306040518263ffffffff1660e01b81526004016103209190610a94565b60206040518083038186803b15801561033857600080fd5b505afa15801561034c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103709190610a4d565b905060008284836103819190610ba1565b61038b9190610b70565b90506000670de0b6b3a7640000600154836103a69190610ba1565b6103b09190610b70565b9050600081866103c09190610bfb565b84866103cc9190610bfb565b836103d79190610ba1565b6103e19190610b70565b9050600081856103f19190610bfb565b90506103fd8289610473565b610407818961077d565b7ffada74d188be3c06536ca96cc80440684fa5d1ecddf2989e65815ec63418bcce896040516104369190610ad8565b60405180910390a1505050505050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015481565b60008173ffffffffffffffffffffffffffffffffffffffff16637a28fb88846040518263ffffffff1660e01b81526004016104ae9190610ad8565b60206040518083038186803b1580156104c657600080fd5b505afa1580156104da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104fe9190610a4d565b90506000808373ffffffffffffffffffffffffffffffffffffffff166327a099d86040518163ffffffff1660e01b815260040160206040518083038186803b15801561054957600080fd5b505afa15801561055d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105819190610a24565b73ffffffffffffffffffffffffffffffffffffffff166362dcfda1846040518263ffffffff1660e01b81526004016105b99190610ad8565b60006040518083038186803b1580156105d157600080fd5b505afa1580156105e5573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525081019061060e919061098f565b91509150805182511461064a577f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b60005b8251811015610775578473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8483815181106106ab577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200260200101518484815181106106ec577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200260200101516040518363ffffffff1660e01b8152600401610711929190610aaf565b602060405180830381600087803b15801561072b57600080fd5b505af115801561073f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061076391906109fb565b508061076e90610cba565b905061064d565b505050505050565b8073ffffffffffffffffffffffffffffffffffffffff1663ee7a7c0430846040518363ffffffff1660e01b81526004016107b8929190610aaf565b602060405180830381600087803b1580156107d257600080fd5b505af11580156107e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061080a9190610a4d565b505050565b600061082261081d84610b18565b610af3565b9050808382526020820190508285602086028201111561084157600080fd5b60005b85811015610871578161085788826108e7565b845260208401935060208301925050600181019050610844565b5050509392505050565b600061088e61088984610b44565b610af3565b905080838252602082019050828560208602820111156108ad57600080fd5b60005b858110156108dd57816108c3888261097a565b8452602084019350602083019250506001810190506108b0565b5050509392505050565b6000815190506108f681610da1565b92915050565b600082601f83011261090d57600080fd5b815161091d84826020860161080f565b91505092915050565b600082601f83011261093757600080fd5b815161094784826020860161087b565b91505092915050565b60008151905061095f81610db8565b92915050565b60008151905061097481610dcf565b92915050565b60008151905061098981610de6565b92915050565b600080604083850312156109a257600080fd5b600083015167ffffffffffffffff8111156109bc57600080fd5b6109c8858286016108fc565b925050602083015167ffffffffffffffff8111156109e557600080fd5b6109f185828601610926565b9150509250929050565b600060208284031215610a0d57600080fd5b6000610a1b84828501610950565b91505092915050565b600060208284031215610a3657600080fd5b6000610a4484828501610965565b91505092915050565b600060208284031215610a5f57600080fd5b6000610a6d8482850161097a565b91505092915050565b610a7f81610c2f565b82525050565b610a8e81610c7f565b82525050565b6000602082019050610aa96000830184610a76565b92915050565b6000604082019050610ac46000830185610a76565b610ad16020830184610a85565b9392505050565b6000602082019050610aed6000830184610a85565b92915050565b6000610afd610b0e565b9050610b098282610c89565b919050565b6000604051905090565b600067ffffffffffffffff821115610b3357610b32610d61565b5b602082029050602081019050919050565b600067ffffffffffffffff821115610b5f57610b5e610d61565b5b602082029050602081019050919050565b6000610b7b82610c7f565b9150610b8683610c7f565b925082610b9657610b95610d32565b5b828204905092915050565b6000610bac82610c7f565b9150610bb783610c7f565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615610bf057610bef610d03565b5b828202905092915050565b6000610c0682610c7f565b9150610c1183610c7f565b925082821015610c2457610c23610d03565b5b828203905092915050565b6000610c3a82610c5f565b9050919050565b60008115159050919050565b6000610c5882610c2f565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b610c9282610d90565b810181811067ffffffffffffffff82111715610cb157610cb0610d61565b5b80604052505050565b6000610cc582610c7f565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610cf857610cf7610d03565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b610daa81610c2f565b8114610db557600080fd5b50565b610dc181610c41565b8114610dcc57600080fd5b50565b610dd881610c4d565b8114610de357600080fd5b50565b610def81610c7f565b8114610dfa57600080fd5b5056fea2646970667358221220d0bf9a90badca709b422e771109bafd4bbc60fe399f9d14c8a6af9ed5ff1d12564736f6c63430008040033`

	lidoABIRAW = `[
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "_lidoAddress",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "_validatorsMevShare",
        "type": "uint256"
      }
    ],
    "stateMutability": "nonpayable",
    "type": "constructor"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "LidoMevDistributed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "LidoMevReceived",
    "type": "event"
  },
  {
    "inputs": [],
    "name": "lidoAddress",
    "outputs": [
      {
        "internalType": "address",
        "name": "",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "validatorsMevShare",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "stateMutability": "payable",
    "type": "receive"
  },
  {
    "inputs": [],
    "name": "distribureMev",
    "outputs": [],
    "stateMutability": "payable",
    "type": "function"
  }
]
`
	bribeContractBin = `0x608060405234801561001057600080fd5b50610a00806100206000396000f3fe6080604052600436106100435760003560e01c8063132bd3a91461005757806320eda6d51461006157806335372a401461009e578063c4d3f9d9146100c757610052565b36610052576100506100d1565b005b600080fd5b61005f610369565b005b34801561006d57600080fd5b506100886004803603810190610083919061060a565b61036b565b60405161009591906107a2565b60405180910390f35b3480156100aa57600080fd5b506100c560048036038101906100c09190610633565b610383565b005b6100cf61040d565b005b60006100dc41610417565b9050600047905060006003826100f291906107d4565b9050600083604001519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146101375780610139565b415b9050600084602001519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461017e5780610180565b415b9050600085600001519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146101c557806101c7565b415b9050856040015173ffffffffffffffffffffffffffffffffffffffff166108fc859081150290604051600060405180830381858888f19350505050158015610213573d6000803e3d6000fd5b50856020015173ffffffffffffffffffffffffffffffffffffffff166108fc859081150290604051600060405180830381858888f1935050505015801561025e573d6000803e3d6000fd5b50856000015173ffffffffffffffffffffffffffffffffffffffff166108fc859081150290604051600060405180830381858888f193505050501580156102a9573d6000803e3d6000fd5b507f95b37ac100e4bf5cd43ff3a48e440813330509e42025bdc7e778b7fa4e2a0c18418433876040516102df9493929190610734565b60405180910390a17f95b37ac100e4bf5cd43ff3a48e440813330509e42025bdc7e778b7fa4e2a0c184183338760405161031c9493929190610734565b60405180910390a17f95b37ac100e4bf5cd43ff3a48e440813330509e42025bdc7e778b7fa4e2a0c18418233876040516103599493929190610734565b60405180910390a1505050505050565b565b610373610577565b61037c82610417565b9050919050565b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081816103ce91906109a5565b9050507fad6568b1bcddeef6bdbc1963e088ecd723c3202b1e9820204a5cbc6a277f5ca03382604051610402929190610779565b60405180910390a150565b6104156100d1565b565b61041f610577565b60008060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206040518060600160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681525050905080915050919050565b6040518060600160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff1681525090565b6000813590506105e9816109b3565b92915050565b60006060828403121561060157600080fd5b81905092915050565b60006020828403121561061c57600080fd5b600061062a848285016105da565b91505092915050565b60006060828403121561064557600080fd5b6000610653848285016105ef565b91505092915050565b61066581610841565b82525050565b61067481610805565b82525050565b61068381610805565b82525050565b6060820161069a60008301836107bd565b6106a7600085018261066b565b506106b560208301836107bd565b6106c2602085018261066b565b506106d060408301836107bd565b6106dd604085018261066b565b50505050565b6060820160008201516106f9600085018261066b565b50602082015161070c602085018261066b565b50604082015161071f604085018261066b565b50505050565b61072e81610837565b82525050565b6000608082019050610749600083018761065c565b610756602083018661067a565b610763604083018561067a565b6107706060830184610725565b95945050505050565b600060808201905061078e600083018561067a565b61079b6020830184610689565b9392505050565b60006060820190506107b760008301846106e3565b92915050565b60006107cc60208401846105da565b905092915050565b60006107df82610837565b91506107ea83610837565b9250826107fa576107f96108f0565b5b828204905092915050565b600061081082610817565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600061084c82610865565b9050919050565b600061085e82610865565b9050919050565b600061087082610877565b9050919050565b600061088282610817565b9050919050565b60008101600083018061089b81610929565b90506108a78184610982565b5050506001810160208301806108bc81610929565b90506108c88184610982565b5050506002810160408301806108dd81610929565b90506108e98184610982565b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000819050919050565b60008135610936816109b3565b80915050919050565b60008160001b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff61096c8461093f565b9350801983169250808416831791505092915050565b61098b82610853565b61099e6109978261091f565b835461094c565b8255505050565b6109af8282610889565b5050565b6109bc81610805565b81146109c757600080fd5b5056fea2646970667358221220bb46835ac5377f1aaa431faac44bb8c7e5072613f397f377a3e203216a2de29664736f6c63430008040033`
	bribeContractABI = `
[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "coinbase",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "receivingAddress",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "msgSender",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "FlashbotsPayment",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "coinbase",
          "type": "address"
        },
        {
          "components": [
            {
              "internalType": "address",
              "name": "executor",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "validator",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "stakingPool",
              "type": "address"
            }
          ],
          "indexed": false,
          "internalType": "struct Recipients",
          "name": "receivingAddress",
          "type": "tuple"
        }
      ],
      "name": "RecipientUpdate",
      "type": "event"
    },
    {
      "stateMutability": "payable",
      "type": "receive"
    },
    {
      "inputs": [
        {
          "components": [
            {
              "internalType": "address",
              "name": "executor",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "validator",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "stakingPool",
              "type": "address"
            }
          ],
          "internalType": "struct Recipients",
          "name": "_newReceivingAddress",
          "type": "tuple"
        }
      ],
      "name": "setRecipients",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_who",
          "type": "address"
        }
      ],
      "name": "getRecipients",
      "outputs": [
        {
          "components": [
            {
              "internalType": "address",
              "name": "executor",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "validator",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "stakingPool",
              "type": "address"
            }
          ],
          "internalType": "struct Recipients",
          "name": "",
          "type": "tuple"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "payMiner",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "queueEther",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    }
  ]
`
)

var (
	clientDial = flag.String(
		"client_dial", "ws://127.0.0.1:8546", "could be websocket or IPC",
	)
	at        = flag.Uint64("kickoff", 2, "what number to kick off at")
	faucet, _ = crypto.HexToECDSA(
		"133be114715e5fe528a1b8adf36792160601a2d63ab59d1fd454275b31328791",
	)
	lidoContractAddr = flag.String(
		"lido_addr", "", "needs to be the 0x addr for lido contract",
	)
	someoneElse, _ = crypto.HexToECDSA("")
	keys           = []*ecdsa.PrivateKey{faucet}
	bribeABI, _    = abi.JSON(strings.NewReader(string(bribeContractABI)))
	lidoABI, _     = abi.JSON(strings.NewReader(string(lidoABIRAW)))
)

func invokeDistributeMEV(
	client *ethclient.Client,
	toAddr common.Address,
	chainID *big.Int,
) (*types.Transaction, error) {

	k := crypto.PubkeyToAddress(someoneElse.PublicKey)
	non, err := client.NonceAt(
		context.Background(), k, nil,
	)
	if err != nil {
		return nil, err
	}

	balance, err := client.BalanceAt(context.Background(), k, nil)
	if err != nil {
		return nil, err
	}
	if balance.Cmp(common.Big0) == 0 {
		return nil, errors.New("need non-zero balance")
	}

	packed, err := lidoABI.Pack("distribureMev")

	t := types.NewTransaction(
		non,
		toAddr,
		big.NewInt(1e17),
		100_000,
		big.NewInt(3e9),
		packed,
	)
	return types.SignTx(t, types.NewEIP155Signer(chainID), someoneElse)
}

func invokeMinerPayment(
	client *ethclient.Client,
	toAddr common.Address,
	chainID *big.Int,
) (*types.Transaction, error) {

	k := crypto.PubkeyToAddress(someoneElse.PublicKey)
	non, err := client.NonceAt(
		context.Background(), k, nil,
	)
	if err != nil {
		return nil, err
	}

	balance, err := client.BalanceAt(context.Background(), k, nil)
	if err != nil {
		return nil, err
	}
	if balance.Cmp(common.Big0) == 0 {
		return nil, errors.New("need non-zero balance")
	}
	t := types.NewTransaction(
		non,
		toAddr,
		new(big.Int),
		100_000,
		big.NewInt(3e9),
		nil,
	)
	return types.SignTx(t, types.NewEIP155Signer(chainID), someoneElse)
}

func deployBribeContract(
	client *ethclient.Client,
	chainID *big.Int,
) (*types.Transaction, error) {
	t := types.NewContractCreation(
		0, new(big.Int), 400_000, big.NewInt(10e9),
		common.Hex2Bytes(bribeContractBin),
	)

	t, err := types.SignTx(t, types.NewEIP155Signer(chainID), faucet)
	if err != nil {
		return nil, err
	}

	return t, client.SendTransaction(context.Background(), t)
}

func deployLidoContract(
	client *ethclient.Client,
	chainID *big.Int,
) (*types.Transaction, error) {
	t := types.NewContractCreation(
		0, new(big.Int), 400_000, big.NewInt(10e9),
		common.Hex2Bytes(lidoContract),
	)

	t, err := types.SignTx(t, types.NewEIP155Signer(chainID), faucet)
	if err != nil {
		return nil, err
	}

	return t, client.SendTransaction(context.Background(), t)
}

func program() error {
	client, err := ethclient.Dial(*clientDial)
	if err != nil {
		return err
	}

	ch := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(
		context.Background(), ch,
	)

	if err != nil {
		return err
	}

	var (
		lightPrismAddr common.Address
		lidoMockAddr   common.Address
	)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}
	deployAt := *at

	for {
		select {
		case e := <-sub.Err():
			return e
		case incoming := <-ch:
			blockNumber := incoming.Number.Uint64()

			if blockNumber == (deployAt - 1) {
				t, err := deployLidoContract(client, chainID)
				if err != nil {
					return err
				}

				lidoMockAddr = crypto.CreateAddress(
					crypto.PubkeyToAddress(faucet.PublicKey),
					t.Nonce(),
				)
				fmt.Println("\tdeployed lido contract ",
					lightPrismAddr.Hex(), "at block number",
					blockNumber)
				continue
			}

			if blockNumber == deployAt {
				t, err := deployBribeContract(client, chainID)
				if err != nil {
					return err
				}

				lightPrismAddr = crypto.CreateAddress(
					crypto.PubkeyToAddress(faucet.PublicKey),
					t.Nonce(),
				)
				fmt.Println("\tdeployed lightPrism contract ",
					lightPrismAddr.Hex(), "at block number",
					blockNumber)
				continue
			}

			if blockNumber > *at {
				// probably ought to just make one contract do this?
				t, err := invokeMinerPayment(client, lightPrismAddr, chainID)
				if err != nil {
					log.Fatal(err)
				}
				if err := client.SendTransaction(context.Background(), t); err != nil {
					log.Fatal(err)
				}

				lidoCall, err := invokeDistributeMEV(client, lidoMockAddr, chainID)
				if err != nil {
					log.Fatal(err)
				}
				if err := client.SendTransaction(context.Background(), lidoCall); err != nil {
					log.Fatal(err)
				}

			}
		}
	}
}

func main() {
	flag.Parse()
	if err := program(); err != nil {
		log.Fatal(err)
	}
}
