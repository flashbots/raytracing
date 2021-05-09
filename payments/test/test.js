const LightPrism = artifacts.require("LightPrism");

// Traditional Truffle test
contract("LightPrism", accounts => {
  it("Should pay", async function() {
    const lightPrism = await LightPrism.new();
    const recipients = {
      executor : '0x0100000000000000000000000000000000000000',
      stakingPool : '0x0200000000000000000000000000000000000000"'
    }
    await lightPrism.setRecipients(recipients);
    await lightPrism.queueEther({value : 3});
  });
});