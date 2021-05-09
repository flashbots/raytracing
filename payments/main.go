package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	m "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	clientDial = flag.String(
		"client_dial", "ws://127.0.0.1:8545", "could be websocket or IPC",
	)
	// meh - these came from ganache - so whatever
	deployerKey, _ = crypto.HexToECDSA(
		"57bae31ef9370140e635ba1fd70707f7a33076827490e4ad5d3eb87784710ce7",
	)
	treasuryKey, _ = crypto.HexToECDSA(
		"0252d6c2476583794ac844385f63040d76f1904978a700e59a37c4e9f68c2f30",
	)
	lidoABI, _            = abi.JSON(strings.NewReader(string("")))
	nodeRegistryABI, _    = abi.JSON(strings.NewReader(string("")))
	depositContractABI, _ = abi.JSON(strings.NewReader(string("")))
)

// func invokeDistributeMEV(
// 	client *ethclient.Client,
// 	toAddr common.Address,
// 	chainID *big.Int,
// ) (*types.Transaction, error) {

// 	k := crypto.PubkeyToAddress(someoneElse.PublicKey)
// 	non, err := client.NonceAt(
// 		context.Background(), k, nil,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	balance, err := client.BalanceAt(context.Background(), k, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if balance.Cmp(common.Big0) == 0 {
// 		return nil, errors.New("need non-zero balance")
// 	}

// 	packed, err := lidoABI.Pack("distribureMev")

// 	t := types.NewTransaction(
// 		non,
// 		toAddr,
// 		big.NewInt(1e17),
// 		100_000,
// 		big.NewInt(3e9),
// 		packed,
// 	)
// 	return types.SignTx(t, types.NewEIP155Signer(chainID), someoneElse)
// }

// func invokeMinerPayment(
// 	client *ethclient.Client,
// 	toAddr common.Address,
// 	chainID *big.Int,
// ) (*types.Transaction, error) {

// 	k := crypto.PubkeyToAddress(someoneElse.PublicKey)
// 	non, err := client.NonceAt(
// 		context.Background(), k, nil,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	balance, err := client.BalanceAt(context.Background(), k, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if balance.Cmp(common.Big0) == 0 {
// 		return nil, errors.New("need non-zero balance")
// 	}
// 	t := types.NewTransaction(
// 		non,
// 		toAddr,
// 		new(big.Int),
// 		100_000,
// 		big.NewInt(3e9),
// 		nil,
// 	)
// 	return types.SignTx(t, types.NewEIP155Signer(chainID), someoneElse)
// }

// func deployBribeContract(
// 	client *ethclient.Client,
// 	chainID *big.Int,
// ) (*types.Transaction, error) {
// 	t := types.NewContractCreation(
// 		0, new(big.Int), 400_000, big.NewInt(10e9),
// 		common.Hex2Bytes(bribeContractBin),
// 	)

// 	t, err := types.SignTx(t, types.NewEIP155Signer(chainID), faucet)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return t, client.SendTransaction(context.Background(), t)
// }

func deployLidoContract(
	from common.Address,
	client *ethclient.Client,
	chainID *big.Int,
	args ...[]byte,
) (*types.Transaction, error) {

	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return nil, err
	}

	payload := common.Hex2Bytes(lidoByteCode)

	for _, k := range args {
		payload = append(payload, k...)
	}

	t := types.NewContractCreation(
		nonce, new(big.Int), 400_000, big.NewInt(10e9), payload,
	)

	t, err = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
	if err != nil {
		return nil, err
	}

	return t, client.SendTransaction(context.Background(), t)
}

func contractDeploy(
	from common.Address,
	client *ethclient.Client,
	chainID *big.Int,
	rawByteCode []byte,
) (*types.Transaction, error) {
	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return nil, err
	}

	t := types.NewContractCreation(
		nonce, new(big.Int), 400_000, big.NewInt(10e9), rawByteCode,
	)

	t, err = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
	if err != nil {
		return nil, err
	}
	return t, client.SendTransaction(context.Background(), t)
}

const (
	blockDeployMock             = 3
	blockDeployOperatorRegistry = 4
	blockDeployLido             = 5
	blockDeployLightPrism       = 6
	blockDeployMEVDistributor   = 7
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func addOperators(
	client *ethclient.Client,
	registry common.Address,
	operators []common.Address,
	deployer common.Address,
	chainID *big.Int,
) error {
	for i, oper := range operators {
		nonce, err := client.NonceAt(context.Background(), deployer, nil)
		if err != nil {
			return err
		}

		packed, err := nodeRegistryABI.Pack("addNodeOperator", fmt.Sprintf("%d", i), oper, 20)
		if err != nil {
			return err
		}
		t := types.NewTransaction(
			nonce, registry, new(big.Int), 200_000, big.NewInt(1e9), packed,
		)
		t, err = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
		if err != nil {
			return err
		}

		if err := client.SendTransaction(context.Background(), t); err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 50)

		pubKeys := make([]byte, 48*20)
		signatures := make([]byte, 96*20)
		rand.Read(pubKeys)
		rand.Read(signatures)

		packed, err = nodeRegistryABI.Pack(
			"addSigningKeys", fmt.Sprintf("%d", i), 20, pubKeys, signatures,
		)

		nonce, err = client.NonceAt(context.Background(), deployer, nil)
		if err != nil {
			return err
		}

		t = types.NewTransaction(
			nonce, registry, new(big.Int), 200_000, big.NewInt(1e9), packed,
		)
		t, err = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
		if err != nil {
			return err
		}

		if err := client.SendTransaction(context.Background(), t); err != nil {
			return err
		}

	}

	return nil
}

type staker struct {
	Key *ecdsa.PrivateKey
	Amt *big.Int
}

func stake(
	client *ethclient.Client,
	lido common.Address,
	stakers []staker,
	oracle common.Address,
	chainID *big.Int,
	deployer common.Address,
) error {

	for _, s := range stakers {
		nonce, err := client.NonceAt(context.Background(), deployer, nil)
		if err != nil {
			return err
		}

		packed, err := lidoABI.Pack("submit", common.Address{})
		if err != nil {
			return err
		}
		t := types.NewTransaction(
			nonce, lido, s.Amt, 200_000, big.NewInt(1e9), packed,
		)

		t, err = types.SignTx(t, types.NewEIP155Signer(chainID), s.Key)
		if err != nil {
			return err
		}

		if err := client.SendTransaction(context.Background(), t); err != nil {
			return err
		}
	}

	nonce, err := client.NonceAt(context.Background(), deployer, nil)
	if err != nil {
		return err
	}

	packed, err := lidoABI.Pack("depositBufferedEther")
	if err != nil {
		return err
	}
	t := types.NewTransaction(
		nonce, lido, new(big.Int), 200_000, big.NewInt(1e9), packed,
	)

	t, err = types.SignTx(t, types.NewEIP155Signer(chainID), stakers[0].Key)
	if err != nil {
		return err
	}

	if err := client.SendTransaction(context.Background(), t); err != nil {
		return err
	}

	// skip the getBeaconStat call, just hardcode call for pushBeacon

	raise := m.Exp(big.NewInt(10), big.NewInt(18))
	stakerCount := big.NewInt(int64(len(stakers)))
	plusOne := new(big.Int).Add(
		new(big.Int).Mul(stakerCount, big.NewInt(32)),
		common.Big1,
	)
	p := new(big.Int).Mul(plusOne, raise)
	packed, err = lidoABI.Pack(
		"pushBeacon",
		stakerCount,
		p,
	)

	if err != nil {
		return err
	}

	t = types.NewTransaction(
		nonce, lido, new(big.Int), 200_000, big.NewInt(1e9), packed,
	)

	t, err = types.SignTx(t, types.NewEIP155Signer(chainID), stakers[0].Key)
	if err != nil {
		return err
	}

	if err := client.SendTransaction(context.Background(), t); err != nil {
		return err
	}

	return nil
}

func distribureMev() {
	//
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
		lightPrismAddr            common.Address
		lidoContractAddr          common.Address
		depositContractAddr       common.Address
		nodeOperatorsRegistryAddr common.Address
		deployerAddr              = crypto.PubkeyToAddress(deployerKey.PublicKey)
		oracleAddr                = deployerAddr
		treasuryAddr              = crypto.PubkeyToAddress(treasuryKey.PublicKey)
	)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}

	for {
		select {
		case e := <-sub.Err():
			return e
		case incoming := <-ch:
			blockNumber := incoming.Number.Uint64()

			if blockNumber == blockDeployMock {
				t, err := contractDeploy(
					deployerAddr,
					client, chainID, common.Hex2Bytes(depositContractByteCode),
				)
				if err != nil {
					log.Fatal(err)
				}
				depositContractAddr = crypto.CreateAddress(deployerAddr, t.Nonce())
				continue
			}

			if blockNumber == blockDeployOperatorRegistry {
				t, err := contractDeploy(
					deployerAddr,
					client, chainID, common.Hex2Bytes(nodeOperatorsRegistryByteCode),
				)
				if err != nil {
					log.Fatal(err)
				}
				nodeOperatorsRegistryAddr = crypto.CreateAddress(
					deployerAddr,
					t.Nonce()+1,
				)
				fmt.Println("deployed node operator registry at ", nodeOperatorsRegistryAddr.Hex())
				continue
			}

			if blockNumber == blockDeployLido {
				t, err := deployLidoContract(
					deployerAddr,
					client, chainID,
					// constructor args
					depositContractAddr.Bytes(),
					oracleAddr.Bytes(),
					nodeOperatorsRegistryAddr.Bytes(),
					treasuryAddr.Bytes(),
					treasuryAddr.Bytes(),
				)

				if err != nil {
					log.Fatal(err)
				}

				lidoContractAddr = crypto.CreateAddress(
					deployerAddr,
					t.Nonce(),
				)

				fmt.Println("deployed lido contract addr ", lidoContractAddr.Hex())
				continue
			}

			if blockNumber == blockDeployLightPrism {
				t, err := deployLidoContract(deployerAddr, client, chainID)
				if err != nil {
					return err
				}

				lightPrismAddr = crypto.CreateAddress(
					crypto.PubkeyToAddress(deployerKey.PublicKey),
					t.Nonce(),
				)

				fmt.Println("\tdeployed light prism contract ", lightPrismAddr.Hex())
				continue
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
