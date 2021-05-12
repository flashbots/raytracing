const LightPrism = artifacts.require("LightPrism");
const Lido = artifacts.require("LidoMevDistributor");

// Traditional Truffle test
contract("LightPrism", (accounts) => {
  it("Should pay", async function () {
    const lightPrism = await LightPrism.deployed();
    // const lido = await Lido.deployed();
    const lido = await Lido.at("0x5304e3c2b42BEaA8f3e37585bFed1274D1055E47");
    // console.log(lido)
    console.log("queue");
    await lightPrism.queueEther({ value: 100 });
    const executor = "0x1230000000000000000000000000000000000321";
    const stakingPool = "0x5304e3c2b42BEaA8f3e37585bFed1274D1055E47";
    // const recipients = {
    //   executor: "0x0000000000000000000000000000000000000000", // zero will be coinbase
    //   stakingPool: "0xabcd000000000000000000000000000000001234", // lido contract here
    // };
    console.log("recipients");
    await lightPrism.setRecipients(executor, stakingPool);
    console.log("queue");
    await lightPrism.queueEther({ value: 100 });
    console.log("pay");
    await lightPrism.payMiner();
    console.log("paid");
    await lido.distribureMev();
    console.log("MEV distributed");
    await lido.explodeToSeeEvents();
  });
});
