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
	keys        = []*ecdsa.PrivateKey{faucet}
	bribeABI, _ = abi.JSON(strings.NewReader(string(bribeContractABI)))
)

func mbTxList(
	client *ethclient.Client,
	toAddr common.Address,
	chainID *big.Int,
) (types.Transactions, error) {

	packed, err := bribeABI.Methods["bribe"].Inputs.Pack()
	if err != nil {
		return nil, err
	}
	txs := make(types.Transactions, len(keys))

	for i, key := range keys {
		k := crypto.PubkeyToAddress(key.PublicKey)
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
			packed,
		)
		t, err = types.SignTx(t, types.NewEIP155Signer(chainID), key)
		if err != nil {
			return nil, err
		}
		txs[i] = t
	}
	return txs, nil
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
		newContractAddr common.Address
		usedTxs         types.Transactions
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
			if blockNumber == deployAt {
				t, err := deployBribeContract(client, chainID)
				if err != nil {
					return err
				}

				newContractAddr = crypto.CreateAddress(
					crypto.PubkeyToAddress(faucet.PublicKey),
					t.Nonce(),
				)
				fmt.Println("\tdeployed bribe contract ", newContractAddr.Hex(), blockNumber)
				continue
			}

			if blockNumber > *at {
				//
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
