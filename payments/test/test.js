const LightPrism = artifacts.require("LightPrism");
const Lido = artifacts.require("LidoMevDistributor");

// Traditional Truffle test
contract("LightPrism", (accounts) => {
  it("Should pay", async function () {
    const lightPrism = await LightPrism.new();
    console.log("queue");
    await lightPrism.queueEther({ value: 3 });
    const recipients = {
      executor: "0x0000000000000000000000000000000000000000", // zero will be coinbase
      stakingPool: "0xabcd000000000000000000000000000000001234", // lido contract here
    };
    console.log("recipients");
    await lightPrism.setRecipients(recipients);
    console.log("queue");
    await lightPrism.queueEther({ value: 3 });
    console.log("pay");
    await lightPrism.payMiner();

    const lido = await Lido.deployed();
  });
});
