const LightPrism = artifacts.require("LightPrism");
const Lido = artifacts.require("LidoMevDistributor");

// Traditional Truffle test
contract("LightPrism", (accounts) => {
  it("Should pay", async function () {
    const lightPrism = await LightPrism.deployed();
    const lido = await Lido.deployed();
    console.log("queue");
    await lightPrism.queueEther({ value: 3 });
    const executor = "0x0000000000000000000000000000000000000000";
    const stakingPool = "0x5f3519e5cbb1d4ba25efea577f61e5ced54a6cea";
    // const recipients = {
    //   executor: "0x0000000000000000000000000000000000000000", // zero will be coinbase
    //   stakingPool: "0xabcd000000000000000000000000000000001234", // lido contract here
    // };
    console.log("recipients");
    await lightPrism.setRecipients(executor, stakingPool);
    console.log("queue");
    await lightPrism.queueEther({ value: 3 });
    console.log("pay");
    await lightPrism.payMiner();
    console.log("paid");
    await lido.distribureMev();
    console.log("MEV distributed");
  });
});
