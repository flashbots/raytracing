const LightPrism = artifacts.require("LightPrism");

// Traditional Truffle test
contract("LightPrism", accounts => {
  it("Should pay", async function() {
    const lightPrism = await LightPrism.new();
    await lightPrism.queueEther.call({value : 3});
  });
});