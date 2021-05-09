const LightPrism = artifacts.require("LightPrism");

// Traditional Truffle test
contract("LightPrism", accounts => {
  it("Should pay", async function() {
    const lightPrism = await LightPrism.new();
    console.log("queue")
    await lightPrism.queueEther({value : 3});
    const recipients = {
      executor :    '0x0100000000000000000000000000000000000000',
      stakingPool : '0x0200000000000000000000000000000000000000'
    }
    console.log("recipients")
    await lightPrism.setRecipients(recipients);
    console.log("queue")
    await lightPrism.queueEther({value : 3});
    console.log("pay")
    await lightPrism.payMiner();
  });
});