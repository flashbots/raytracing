package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	m "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

const (
	addrsDeployedFile = "deployed-addrs.json"
)

var (
	checkHash         = flag.String("check", "", "check the txn hash")
	checkAddrContract = flag.String("is_contract", "", "is it a contract")
	addrsDeployed     = flag.String("addrs", addrsDeployedFile, "already deployed addrs")
	freshDeploy       = flag.Bool("fresh", true, "fresh contract deployments")
	clientDial        = flag.String(
		"client_dial", eth1, "could be websocket or IPC",
	)

	deployerKey, _ = crypto.HexToECDSA(
		"7074988e20b9aa7c58ea6dd5a56aaf5faf4bedc2ea7da7b02adfc97c92b7ceb3",
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
	nonce uint64,
	from common.Address,
	client *ethclient.Client,
	chainID *big.Int,
	args []byte,
) (*types.Transaction, error) {

	payload := append(common.Hex2Bytes(lidoByteCode), args...)
	fmt.Println("using nonce ", nonce, "for lido deploy")
	t := types.NewContractCreation(
		nonce, new(big.Int), 400_0000, big.NewInt(10e9), payload,
	)

	t, err := types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
	if err != nil {
		return nil, err
	}

	fmt.Println("deployed lido contract txn hash", t.Hash().Hex())

	return t, client.SendTransaction(context.Background(), t)
}

func deployLightPrism(
	nonce uint64,
	from common.Address,
	client *ethclient.Client,
	chainID *big.Int,
) (*types.Transaction, error) {
	fmt.Println("using nonce", nonce, "for light prism deploy")

	t := types.NewContractCreation(
		nonce, new(big.Int), 4_000_000, big.NewInt(10e9), common.Hex2Bytes(lightPrismByteCode),
	)

	t, err := types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
	if err != nil {
		return nil, err
	}

	fmt.Println("deployed light prism contract ", t.Hash().Hex())
	time.Sleep(time.Second * 17)

	return t, client.SendTransaction(context.Background(), t)
}

func deployMEVDistributor(
	nonce uint64,
	from common.Address,
	client *ethclient.Client,
	chainID *big.Int,
	args []byte,
) (*types.Transaction, error) {

	payload := append(common.Hex2Bytes(lidoMEVdistribByteCode), args...)
	fmt.Println("using nonce", nonce, "for deply mev distributor")
	t := types.NewContractCreation(
		nonce, new(big.Int), 4_000_000, big.NewInt(10e9), payload,
	)

	t, err := types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("deployed MEV distributor contract ", t.Hash().Hex())
	return t, client.SendTransaction(context.Background(), t)
}

func contractDeploy(
	name string,
	nonce uint64,
	from common.Address,
	client *ethclient.Client,
	chainID *big.Int,
	rawByteCode []byte,
) (*types.Transaction, error) {

	fmt.Println("using nonce", nonce, "for contract deploy", name)
	t := types.NewContractCreation(
		nonce, new(big.Int), 400_000_0, big.NewInt(10e9), rawByteCode,
	)

	t, err := types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("a generic contract deployed", t.Hash().Hex())
	return t, client.SendTransaction(context.Background(), t)
}

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
	nonce_ uint64,
	client *ethclient.Client,
	registry common.Address,
	operators []common.Address,
	deployer common.Address,
	chainID *big.Int,
) (uint64, error) {

	nonce := nonce_

	for i, oper := range operators {

		packed, err := abiRegistry.Pack("addNodeOperator", fmt.Sprintf("%d", i), oper, uint64(20))
		if err != nil {
			return 0, err
		}

		nonce = nonce + uint64(i)

		t := types.NewTransaction(
			nonce, registry, new(big.Int), 500_000, big.NewInt(1e9), packed,
		)
		t, err = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
		if err != nil {
			return 0, err
		}

		if err := client.SendTransaction(context.Background(), t); err != nil {
			return 0, err
		}

		pubKeys := make([]byte, 48*20)
		signatures := make([]byte, 96*20)
		rand.Read(pubKeys)
		rand.Read(signatures)

		packed, err = abiRegistry.Pack(
			"addSigningKeys", fmt.Sprintf("%d", i), 20, pubKeys, signatures,
		)

		nonce++

		t = types.NewTransaction(
			nonce, registry, new(big.Int), 200_000, big.NewInt(1e9), packed,
		)

		t, err = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
		if err != nil {
			return 0, err
		}

		if err := client.SendTransaction(context.Background(), t); err != nil {
			return 0, err
		}

	}

	return nonce, nil
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

const (
	eth1 = "ws://138.68.75.41:8546/"
	eth2 = "ws://138.68.75.41:5051/"
)

func program() error {
	client, err := ethclient.Dial(*clientDial)
	if err != nil {
		return err
	}

	type NeededAddrs struct {
		LightPrismAddr            common.Address
		LidoContractAddr          common.Address
		LidoMEVContractAddr       common.Address
		DepositContractAddr       common.Address
		NodeOperatorsRegistryAddr common.Address
		OracleAddr                common.Address
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
	networkID, _ := client.NetworkID(context.Background())
	fmt.Println("using chain id", chainID, "networkID", networkID)

	if err != nil {
		return err
	}

	deployerBal, _ := client.BalanceAt(context.Background(), deployerAddr, nil)
	deployerNonce, _ := client.NonceAt(context.Background(), deployerAddr, nil)
	fmt.Println(
		"deployer addr", deployerAddr.Hex(), " has this much eth on hand\n",
		deployerBal, "\n",
		"and this nonce", deployerNonce,
	)

	if a := *checkAddrContract; a != "" {
		bytes, err := client.CodeAt(context.Background(), common.HexToAddress(a), nil)
		if err != nil {
			return err
		}

		if bytes != nil {
			fmt.Println(common.HexToAddress(a).Hex(), "is a contract")
		} else {
			fmt.Println("not deployed as a contract")
		}
		return nil
	}

	if *checkHash != "" {

		recipt, err := client.TransactionReceipt(
			context.Background(), common.HexToHash(*checkHash),
		)

		if err != nil {
			return err
		}

		fmt.Println("confirmed txn!", recipt)
		return nil
	}

	if *freshDeploy {

		currentBlock, err := client.BlockByNumber(context.Background(), nil)
		if err != nil {
			return err
		}

		fmt.Println("beginning deployment of contracts at live block number",
			currentBlock.Header().Number,
		)

		nonce, err := client.NonceAt(
			context.Background(), deployerAddr, nil,
		)

		t, err := contractDeploy(
			"deposit contract",
			nonce,
			deployerAddr,
			client, chainID, common.Hex2Bytes(depositContractByteCode),
		)
		if err != nil {
			return errors.Wrapf(err, "deposit contract")
		}

		depositContractAddr = crypto.CreateAddress(deployerAddr, t.Nonce())

		t, err = contractDeploy(
			"node operators registry ",
			nonce+1,
			deployerAddr,
			client, chainID, common.Hex2Bytes(nodeOperatorsRegistryByteCode),
		)
		if err != nil {
			return errors.Wrapf(err, "node operators registry")
		}
		nodeOperatorsRegistryAddr = crypto.CreateAddress(
			deployerAddr,
			t.Nonce(),
		)

		fmt.Println("deployed node operator registry at ", nodeOperatorsRegistryAddr.Hex())
		packed, err := abiLido.Pack("", depositContractAddr,
			oracleAddr,
			nodeOperatorsRegistryAddr,
			treasuryAddr,
			treasuryAddr,
		)

		if err != nil {
			return errors.Wrapf(err, "deposit contract addr")
		}

		t2, err := deployLidoContract(
			nonce+2,
			deployerAddr, client, chainID, packed,
		)

		if err != nil {
			return errors.Wrapf(err, "deposit contract addr")
		}

		lidoContractAddr = crypto.CreateAddress(deployerAddr, t2.Nonce())
		fmt.Println("deployed lido contract addr ", lidoContractAddr.Hex())
		packed, err = abiRegistry.Pack("setLido", lidoContractAddr)
		if err != nil {
			return errors.Wrapf(err, "deposit contract addr")
		}
		t = types.NewTransaction(
			nonce+3, nodeOperatorsRegistryAddr, new(big.Int),
			500_000, big.NewInt(1e9), packed,
		)

		t, _ = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)

		if err := client.SendTransaction(context.Background(), t); err != nil {
			return err
		}

		cred := [32]byte{}
		rand.Read(cred[:])

		packed, err = abiLido.Pack("setWithdrawalCredentials", cred)
		t = types.NewTransaction(
			nonce+4, lidoContractAddr, new(big.Int),
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
			nonce+5, lidoContractAddr, new(big.Int),
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
			nonce+6,
			deployerAddr, client, chainID, packed,
		)

		if err != nil {
			return err
		}

		lidoMEVContractAddr = crypto.CreateAddress(deployerAddr, t.Nonce())
		fmt.Println("lidoMEV contract addr deployed at", lidoMEVContractAddr.Hex())

		nonce, err = addOperators(
			nonce+7, client, nodeOperatorsRegistryAddr,
			operators, deployerAddr, chainID,
		)

		if err != nil {
			return err
		}

		if err := stake(
			client, lidoContractAddr, stakers, oracleAddr, chainID, deployerAddr,
		); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Added operators & added stakers, deploying light prism")
		t, err = deployLightPrism(nonce+1, deployerAddr, client, chainID)
		if err != nil {
			return err
		}

		lightPrismAddr = crypto.CreateAddress(
			crypto.PubkeyToAddress(deployerKey.PublicKey),
			t.Nonce(),
		)

		{
			packed, err := abiLightPrism.Pack("queueEther")
			if err != nil {
				return err
			}

			t = types.NewTransaction(
				nonce+2, lightPrismAddr, big.NewInt(3e18),
				200_000, big.NewInt(1e9), packed,
			)
			t, _ = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
			time.Sleep(time.Second * 17)
			if err := client.SendTransaction(context.Background(), t); err != nil {
				log.Fatal(err)
			}
		}

		packed, err = abiLightPrism.Pack(
			"setRecipients", common.Address{}, lidoContractAddr,
		)
		if err != nil {
			return err
		}

		t = types.NewTransaction(
			nonce+3, lightPrismAddr, new(big.Int),
			200_000, big.NewInt(1e9), packed,
		)
		t, _ = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)
		time.Sleep(time.Second * 17)

		if err := client.SendTransaction(context.Background(), t); err != nil {
			log.Fatal(err)
		}
		fmt.Println("set recipients worked")

		allAddrs := NeededAddrs{
			LightPrismAddr:            lightPrismAddr,
			LidoContractAddr:          lidoContractAddr,
			LidoMEVContractAddr:       lidoMEVContractAddr,
			DepositContractAddr:       depositContractAddr,
			NodeOperatorsRegistryAddr: nodeOperatorsRegistryAddr,
			OracleAddr:                deployerAddr,
		}
		s, _ := json.MarshalIndent(allAddrs, "  ", "  ")
		if err := ioutil.WriteFile(addrsDeployedFile, s, 0644); err != nil {
			fmt.Println("some problem on writing the file of nodes")
		}

		return nil
	} else {
		// not a fresh deployment - so lets
		var p NeededAddrs
		common.LoadJSON(addrsDeployedFile, &p)
		pret, _ := json.MarshalIndent(p, " ", " ")
		fmt.Println("Loaded up previously deployed addrs ", string(pret))
		query := ethereum.FilterQuery{
			Addresses: []common.Address{
				p.LightPrismAddr, p.LidoContractAddr, p.LidoMEVContractAddr,
				p.DepositContractAddr, p.NodeOperatorsRegistryAddr,
			},
		}

		logs := make(chan types.Log)
		sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
		if err != nil {
			return err
		}

		go func() {
			tick := time.NewTicker(time.Second * 10)
			for range tick.C {
				fmt.Println("Kicking off arbitrage txn")

				packedQueueEther, _ := abiLightPrism.Pack("queueEther")
				nonce, _ := client.NonceAt(
					context.Background(), deployerAddr, nil,
				)
				packedQueueEther = nil

				t := types.NewTransaction(
					nonce, p.LightPrismAddr, big.NewInt(3e17),
					200_000, big.NewInt(10e9), packedQueueEther,
				)
				t, _ = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)

				fmt.Println("packed queue ether", t.Hash().Hex())

				if err := client.SendTransaction(context.Background(), t); err != nil {
					log.Fatal(err)
				}

				packedPayMiner, err := abiLightPrism.Pack("payMiner")
				if err != nil {
					log.Fatal(err)
				}

				{
					nonce, err = client.NonceAt(
						context.Background(), deployerAddr, nil,
					)

					packedPayMiner = nil

					t = types.NewTransaction(
						nonce, p.LightPrismAddr, big.NewInt(1e18),
						200_000, big.NewInt(10e9), packedPayMiner,
					)

					t, _ = types.SignTx(t, types.NewEIP155Signer(chainID), deployerKey)

					pret, _ := json.MarshalIndent(t, " ", " ")

					fmt.Println(
						"packed payminer called txn hash:",
						t.Hash().Hex(), string(pret),
					)

					if err := client.SendTransaction(context.Background(), t); err != nil {
						log.Fatal(err)
					}
				}

				fmt.Println("arbitrage round ended")
			}
		}()

		type PrismFlashbotsPayment struct {
			Coinbase         common.Address
			ReceivingAddress common.Address
			MsgSender        common.Address
			Amount           *big.Int
			Raw              types.Log // Blockchain specific contextual infos
		}

		type Recipients struct {
			Executor    common.Address
			StakingPool common.Address
		}

		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case vLog := <-logs:
				event := new(PrismFlashbotsPayment)
				if err := abiLightPrism.UnpackIntoInterface(
					event, "FlashbotsPayment", vLog.Data,
				); err != nil {
					fmt.Println("problem on unpacking inerface", err)
				}
				fmt.Println("did event", event)
			}
		}
	}

	return nil
}

func main() {
	flag.Parse()

	if err := program(); err != nil {
		log.Fatal(err)
	}
}
