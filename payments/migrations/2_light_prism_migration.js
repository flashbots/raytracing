const LightPrism = artifacts.require("LightPrism");
const Lido = artifacts.require("LidoMevDistributor");

module.exports = function (deployer) {
  deployer.deploy(LightPrism);
  deployer.deploy(Lido);
};
