const LightPrism = artifacts.require("LightPrism");
const LidoMevDistributor = artifacts.require("LidoMevDistributor");
const NodeOperatorsRegistry = artifacts.require("NodeOperatorsRegistry");
const DepositContractMock = artifacts.require("DepositContractMock");
const Lido = artifacts.require("Lido");

module.exports = function (deployer) {
  var operatorsAddress = deployer.deploy(NodeOperatorsRegistry);
  var depositAddress = deployer.deploy(DepositContractMock);
  console.log(address);
  var lightPrismAddress = deployer.deploy(LightPrism);
  var lidoAddress = deployer.deploy(Lido,
    depositAddress,
    `0x0000000000000000000000000000000000000001`,
    operatorsAddress,
    `0x0000000000000000000000000000000000000002`,
    `0x0000000000000000000000000000000000000003`,
    );
  deployer.deploy(
    lidoAddress,
    `0x5d6d0199912b58220ac15661427af6bf53926385`,
    `${4e18}`
  );
};
