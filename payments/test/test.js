const LightPrism = artifacts.require("LightPrism");
const LidoMevDistributor = artifacts.require("LidoMevDistributor");
const ERC20 = artifacts.require("IERC20");

const nodeOperators = [
    '0x15a2209CE319914B7E6e39FEF2d04A764e522f7E',
    '0xFA5B77d2b4A5071B6Da9D85787270260026d1042',
    '0x238F1b5Db4fC0Bff012278d59A1B60e915f55d40',
]

const stakers = [
    '0x40acf4b2c6894aBBF880dbcb46791cE2bB273882',
    '0x14fe850e76Af0617eb5c01D525aB451ecc9f55eE',
    '0x5E90A2d51653CDc1Be80aAeF35Cd41F7cb34c3e4'
]

// Traditional Truffle test
contract("LightPrism", (accounts) => {
  it("Should pay", async function () {
    const lightPrism = await LightPrism.deployed();
    const lidoDistr = await LidoMevDistributor.deployed();
    const steth = await ERC20.at("0x8953454B243E11012DD63b1849D6a4cdb64aB3EB");
    console.log("queue");
    await lightPrism.queueEther({ value: 3 });
    const executor = "0x0000000000000000000000000000000000000000";
    const stakingPool = "0x5f3519e5cbb1d4ba25efea577f61e5ced54a6cea";
    // const recipients = {
    //   executor: "0x0000000000000000000000000000000000000000", // zero will be coinbase
    //   stakingPool: "0xabcd000000000000000000000000000000001234", // lidoDistr contract here
    // };

    console.log("stETH balances before");
    const opsBalances = await Promise.all(nodeOperators.map(op => steth.balanceOf(op)))
    console.log("node operators:", opsBalances.map(x => x.toString()));
    const stakersBalances = await Promise.all(stakers.map(op => steth.balanceOf(op)))
    console.log("stakers:", stakersBalances.map(x => x.toString()));

    console.log("recipients");
    await lightPrism.setRecipients(executor, stakingPool);
    console.log("queue");
    await lightPrism.queueEther({ value: 3 });
    console.log("pay");
    await lightPrism.payMiner();
    console.log("paid");
    await lidoDistr.distribureMev();
    console.log("MEV distributed");

    console.log("stETH balances after");
    const opsBalances = await Promise.all(nodeOperators.map(op => steth.balanceOf(op)))
    console.log("node operators:", opsBalances.map(x => x.toString()));
    const stakersBalances = await Promise.all(stakers.map(op => steth.balanceOf(op)))
    console.log("stakers:", stakersBalances.map(x => x.toString()));
  });
});
