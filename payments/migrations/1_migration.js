const LightPrism = artifacts.require("LightPrism");
const Lido = artifacts.require("LidoMevDistributor");

module.exports = function (deployer) {
  deployer.deploy(LightPrism);
  deployer.deploy(
    Lido,
    `0x5d6d0199912b58220ac15661427af6bf53926385`,
    `${1e18}`
  );
};
