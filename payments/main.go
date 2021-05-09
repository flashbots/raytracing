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
	abiLido, _            = abi.JSON(strings.NewReader(string(lidoABI)))
	abiRegistry, _        = abi.JSON(strings.NewReader(string(nodeRegistryABI)))
	abiDepositContract, _ = abi.JSON(strings.NewReader(string(depositABI)))
	abiMEVLido, _         = abi.JSON(strings.NewReader(string(lidoMEVABI)))
	abiLightPrism, _      = abi.JSON(strings.NewReader(string(lightPrismABI)))
)

func deployLidoContract(
	from common.Address,
	client *ethclient.Client,
	chainID *big.Int,
	args []byte,
) (*types.Transaction, error) {

	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return nil, err
	}

	payload := append(common.Hex2Bytes(lidoByteCode), args...)

	t := types.NewContractCreation(
		nonce, new(big.Int), 400_000, big.NewInt(10e9), payload,
	)

	t, err = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
	if err != nil {
		return nil, err
	}

	return t, client.SendTransaction(context.Background(), t)
}

func deployLightPrism(
	from common.Address,
	client *ethclient.Client,
	chainID *big.Int,
) (*types.Transaction, error) {

	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return nil, err
	}

	t := types.NewContractCreation(
		nonce, new(big.Int), 400_000, big.NewInt(10e9), common.Hex2Bytes(lightPrismByteCode),
	)

	t, err = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
	if err != nil {
		return nil, err
	}

	return t, client.SendTransaction(context.Background(), t)
}

func deployMEVDistributor(
	from common.Address,
	client *ethclient.Client,
	chainID *big.Int,
	args []byte,
) (*types.Transaction, error) {

	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return nil, err
	}

	payload := append(common.Hex2Bytes(lidoMEVdistribByteCode), args...)

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

var (
	stakers   []staker
	operators []common.Address
)

func init() {
	rand.Seed(time.Now().UnixNano())
	s1, _ := crypto.HexToECDSA(
		"98f6682f8e413e40bfe68e7551ea8c6d005f40aa520681f1dfccdf7238784113",
	)
	s2, _ := crypto.HexToECDSA(
		"cc72af8c0acadaaabcbc1a485a4ea35380bb50633d36dd7dae57d99de9a5f639",
	)
	s3, _ := crypto.HexToECDSA(
		"9e3dcc452413270fe82ac4e64e8cc989859b44a0caf00b336d6fcf7ff8d6dbd2",
	)
	operators = []common.Address{
		crypto.PubkeyToAddress(s1.PublicKey),
		crypto.PubkeyToAddress(s2.PublicKey),
		crypto.PubkeyToAddress(s3.PublicKey),
	}

	stakers = []staker{
		{s1, big.NewInt(1e18)},
		{s2, big.NewInt(1e18)},
		{s3, big.NewInt(1e18)},
	}
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

		packed, err := abiRegistry.Pack("addNodeOperator", fmt.Sprintf("%d", i), oper, uint64(20))
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

		packed, err = abiRegistry.Pack(
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
		nonce, err := client.NonceAt(
			context.Background(),
			crypto.PubkeyToAddress(s.Key.PublicKey),
			nil,
		)
		if err != nil {
			return err
		}

		packed, err := abiLido.Pack("submit", common.Address{})
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
		time.Sleep(time.Millisecond * 100)
		if err := client.SendTransaction(context.Background(), t); err != nil {
			return err
		}
	}

	nonce, err := client.NonceAt(
		context.Background(),
		crypto.PubkeyToAddress(stakers[0].Key.PublicKey), nil,
	)
	if err != nil {
		return err
	}

	packed, err := abiLido.Pack("depositBufferedEther", big.NewInt(9e18))
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
	time.Sleep(time.Millisecond * 100)

	if err := client.SendTransaction(context.Background(), t); err != nil {
		fmt.Println("died here C")
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
	packed, err = abiLido.Pack(
		"pushBeacon",
		stakerCount,
		p,
	)

	if err != nil {
		return err
	}

	nonce, err = client.NonceAt(
		context.Background(),
		crypto.PubkeyToAddress(stakers[0].Key.PublicKey),
		nil,
	)
	if err != nil {
		return err
	}

	t = types.NewTransaction(
		nonce, lido, new(big.Int), 200_000, big.NewInt(1e9), packed,
	)

	time.Sleep(time.Millisecond * 100)
	t, err = types.SignTx(t, types.NewEIP155Signer(chainID), stakers[0].Key)

	if err != nil {
		return err
	}

	return client.SendTransaction(context.Background(), t)
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
		lidoMEVContractAddr       common.Address
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
				fmt.Println("Deploying Lido Contract")
				packed, err := abiLido.Pack("", depositContractAddr,
					oracleAddr,
					nodeOperatorsRegistryAddr,
					treasuryAddr,
					treasuryAddr,
				)
				if err != nil {
					log.Fatal(err)
				}

				t, err := deployLidoContract(
					deployerAddr, client, chainID, packed,
				)

				if err != nil {
					log.Fatal(err)
				}

				lidoContractAddr = crypto.CreateAddress(deployerAddr, t.Nonce())
				fmt.Println("deployed lido contract addr ", lidoContractAddr.Hex())

				packed, err = abiRegistry.Pack("setLido", lidoContractAddr)

				if err != nil {
					fmt.Println("died a")
					return err
				}

				nonce, err := client.NonceAt(
					context.Background(), deployerAddr, nil,
				)
				t = types.NewTransaction(
					nonce, nodeOperatorsRegistryAddr, new(big.Int),
					200_000, big.NewInt(1e9), packed,
				)
				t, _ = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
				if err := client.SendTransaction(context.Background(), t); err != nil {
					log.Fatal(err)
				}

				cred := [32]byte{}
				rand.Read(cred[:])

				packed, err = abiLido.Pack("setWithdrawalCredentials", cred)

				if err != nil {
					log.Fatal(err)
				}
				nonce, err = client.NonceAt(
					context.Background(), deployerAddr, nil,
				)
				t = types.NewTransaction(
					nonce, lidoContractAddr, new(big.Int),
					200_000, big.NewInt(1e9), packed,
				)
				t, _ = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
				if err := client.SendTransaction(context.Background(), t); err != nil {
					log.Fatal(err)
				}

				packed, err = abiLido.Pack("resume")
				if err != nil {
					log.Fatal(err)
				}
				nonce, err = client.NonceAt(
					context.Background(), deployerAddr, nil,
				)
				t = types.NewTransaction(
					nonce, lidoContractAddr, new(big.Int),
					200_000, big.NewInt(1e9), packed,
				)
				t, _ = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
				if err := client.SendTransaction(context.Background(), t); err != nil {
					return err
				}

				packed, err = abiMEVLido.Pack(
					"", lidoContractAddr, new(big.Int).Mul(big.NewInt(5), big.NewInt(10e17)),
				)

				if err != nil {
					return err
				}

				t, err = deployMEVDistributor(
					deployerAddr, client, chainID, packed,
				)

				if err != nil {
					return err
				}

				lidoMEVContractAddr = crypto.CreateAddress(deployerAddr, t.Nonce())
				fmt.Println("lidoMEV contract addr deployed at", lidoMEVContractAddr.Hex())
				// # 50% of the received MEV goes to validators, 50% to stakers
				// distributor = LidoMevDistributor.deploy(lido, 5 * 10**17, {'from': deployer})

				if err := addOperators(client, nodeOperatorsRegistryAddr,
					operators, deployerAddr, chainID,
				); err != nil {
					log.Fatal(err)
				}

				if err := stake(
					client, lidoContractAddr, stakers, oracleAddr, chainID, deployerAddr,
				); err != nil {
					fmt.Println("died here B")
					log.Fatal(err)
				}
				fmt.Println("Added operators & added stakers")
				continue
			}

			if blockNumber == blockDeployLightPrism {
				t, err := deployLightPrism(deployerAddr, client, chainID)
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
