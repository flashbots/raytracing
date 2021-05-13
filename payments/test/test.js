const LightPrism = artifacts.require("LightPrism");
const LidoMevDistributor = artifacts.require("LidoMevDistributor");
const ERC20 = artifacts.require("IERC20");
const fetch = require("node-fetch");
const { Headers } = fetch;

const headers = new Headers({
  "content-type": "application/json",
  accept: "application/json",
});

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

const send_bundle = async (url, txn) => {
  const resp = await fetch(url, {
    headers,
    body: JSON.stringify({
      id: 1,
      jsonrpc: "2.0",
      method: "eth_sendBundle",
      params: [txn],
    }),
  });

  const json = await resp.json();
  return JSON.parse(json.result);
};

// Traditional Truffle test
contract("LightPrism", (accounts) => {
  it("Should pay", async function () {
    
    const lightPrism = await LightPrism.deployed();
    const lidoDistr = await LidoMevDistributor.at("0x775A136c9bB5669677185dEEE09051A7382B1574");
    const steth = await ERC20.at("0x218bC3c5AC79Ba91c8309D01f8D1466Db560dd23");
    const executor = "0x3210000000000000000000000000000000000123";
    const stakingPool = "0x775A136c9bB5669677185dEEE09051A7382B1574";

    let drawBalances = async function() {
      await showBalance(lightPrism.address,  "lightPrism ");
      await showBalance(executor,    "coinbase   ");
      await showBalance(stakingPool, "stakingPool");
      draw____________________________________________();
    };

    let drawStaking = async function() {
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
      draw____________________________________________();
    };

    drawTitle("Presenting a Flashbots MEV-bundle (a list of transactions) sent to a Nethermind node. The bundle includes a transaction that sends an MEV tip to the LightPrism contract which then splits it between the coinbase (Eth1 node operator) and Lido staking pool (further distributing between stakers and validator nodes operators). The Eth1 block is generated in response to the assemble block call from an Eth2 client (Teku).");
    console.log("MEV tip");
    console.log("|");
    console.log("|");
    console.log("|");
    console.log("_______1/3 to coinbase");
    console.log("_______2/3 to staking pool");
    console.log("       |");
    console.log("       |");
    console.log("       |");
    console.log("       |");
    console.log("       _______x%        to validators");
    console.log("       _______100% - x% to stakers");
    draw____________________________________________();
    console.log("");

    drawTitle("Setting recipients of the MEV tip (this needs to be set only once by each coinbase)");
    await lightPrism.setRecipients(executor, stakingPool);
    await waitForTx();
    console.log("   coinbase    <- " + executor);
    console.log("   stakingPool <- " + stakingPool);
    draw____________________________________________();
    await drawBalances();
    console.log("");
    drawTitle("MEV bundle enqueues 0.001 ETH");
    console.log("");
    await lightPrism.queueEther({ value: 1000000000000000 });
    await waitForTx();
    // const recipients = {
    //   executor: "0x0000000000000000000000000000000000000000", // zero will be coinbase
    //   stakingPool: "0xabcd000000000000000000000000000000001234", // lidoDistr contract here
    // };

    drawTitle("stETH balances before");
    await drawBalances();
    await drawStaking();

    console.log("");
    drawTitle("MEV bundle enqueues 0.001 ETH");
    await lightPrism.queueEther({ value: 1000000000000000 });
    await waitForTx();
    await drawBalances();
    console.log("");
    drawTitle("MEV bundle makes a payment to a contract which splits it between the coinbase and the staking pool");
    await lightPrism.payMiner();
    await waitForTx();
    await drawBalances();
    console.log("");
    drawTitle("Lido distributes MEV-tip between stakers and validator nodes operators");
    await lidoDistr.distribureMev();
    await waitForTx();
    await drawBalances();
    console.log("");
    drawTitle("stETH balances after");
    await drawStaking();
    await drawBalances();
  });
});

function drawTitle(titleText) {
  draw____________________________________________();
  console.log(titleText);
  draw____________________________________________();
}

function draw____________________________________________() {
  console.log("================================================================");
}

async function waitForTx() {
  await new Promise(r => setTimeout(r, 500));
}

async function showBalance(address, name) {
  let balance = await web3.eth.getBalance(address);
  console.log(`${name}:   `, balance);
}

