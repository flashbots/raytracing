const LightPrism = artifacts.require("LightPrism");
const LidoMevDistributor = artifacts.require("LidoMevDistributor");
const ERC20 = artifacts.require("IERC20");

const nodeOperators = [
  "0x15a2209CE319914B7E6e39FEF2d04A764e522f7E",
  "0xFA5B77d2b4A5071B6Da9D85787270260026d1042",
  "0x238F1b5Db4fC0Bff012278d59A1B60e915f55d40",
];

const stakers = [
  "0x40acf4b2c6894aBBF880dbcb46791cE2bB273882",
  "0x14fe850e76Af0617eb5c01D525aB451ecc9f55eE",
  "0x5E90A2d51653CDc1Be80aAeF35Cd41F7cb34c3e4",
];

// Traditional Truffle test
contract("LightPrism", (accounts) => {
  it("Should pay", async function () {
    const lightPrism = await LightPrism.deployed();
    const lidoDistr = await LidoMevDistributor.at('0x5304e3c2b42BEaA8f3e37585bFed1274D1055E47');
    const steth = await ERC20.at("0x8953454B243E11012DD63b1849D6a4cdb64aB3EB");

    const executor = "0x3210000000000000000000000000000000000123";
    const stakingPool = "0x5304e3c2b42BEaA8f3e37585bFed1274D1055E47";
    console.log("");
    console.log("=============================");
    console.log("Setting recipients of the MEV tip (this needs to be set only once by each coinbase)");
    await lightPrism.setRecipients(executor, stakingPool);
    console.log("   coinbase    <- " + executor);
    console.log("   stakingPool <- " + stakingPool);
    console.log("=============================");
    console.log("");
    let executorBalance = await web3.eth.getBalance(executor);
    console.log("coinbase:", executorBalance);
    let stakingBalance = await web3.eth.getBalance(stakingPool);
    console.log("stakingPool:", stakingBalance);
    console.log("");
    console.log("=============================");
    console.log("MEV bundle enqueues 0.001 ETH");
    console.log("=============================");
    console.log("");
    await lightPrism.queueEther({ value: 1000000000000000 });
    // const recipients = {
    //   executor: "0x0000000000000000000000000000000000000000", // zero will be coinbase
    //   stakingPool: "0xabcd000000000000000000000000000000001234", // lidoDistr contract here
    // };

    console.log("=============================");
    console.log("stETH balances before");
    console.log("=============================");

    executorBalance = await web3.eth.getBalance(executor);
    console.log("coinbase:", executorBalance);
    stakingBalance = await web3.eth.getBalance(stakingPool);
    console.log("stakingPool:", stakingBalance);

    let opsBalances = await Promise.all(
      nodeOperators.map((op) => steth.balanceOf(op))
    );
    console.log(
      "node operators:",
      opsBalances.map((x) => x.toString())
    );
    let stakersBalances = await Promise.all(
      stakers.map((op) => steth.balanceOf(op))
    );
    console.log(
      "stakers:       ",
      stakersBalances.map((x) => x.toString())
    );

    console.log("");
    console.log("=============================");
    console.log("MEV bundle enqueues 0.001 ETH");
    console.log("=============================");
    console.log("");
    await lightPrism.queueEther({ value: 1000000000000000 });

    console.log("=============================");
    console.log("MEV bundle makes a payment to a contract which splits it between the coinbase and the staking pool");
    console.log("=============================");

    await lightPrism.payMiner();
    executorBalance = await web3.eth.getBalance(executor);
    console.log("coinbase:", executorBalance);
    stakingBalance = await web3.eth.getBalance(stakingPool);
    console.log("stakingPool:", stakingBalance);

    console.log("=============================");
    console.log("Lido distributes MEV-tip between stakers and validator nodes operators");
    await lidoDistr.distribureMev();
    console.log("=============================");
    console.log("stETH balances after");
    console.log("=============================");

    executorBalance = await web3.eth.getBalance(executor);
    console.log("coinbase:", executorBalance);
    stakingBalance = await web3.eth.getBalance(stakingPool);
    console.log("stakingPool:", stakingBalance);

    opsBalances = await Promise.all(
      nodeOperators.map((op) => steth.balanceOf(op))
    );
    console.log(
      "node operators:",
      opsBalances.map((x) => x.toString())
    );
    stakersBalances = await Promise.all(
      stakers.map((op) => steth.balanceOf(op))
    );
    console.log(
      "stakers:       ",
      stakersBalances.map((x) => x.toString())
    );
    console.log("=============================");
  });
});
