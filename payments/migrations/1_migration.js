const LightPrism = artifacts.require("LightPrism");
const LidoMevDistributor = artifacts.require("LidoMevDistributor");
const NodeOperatorsRegistry = artifacts.require("NodeOperatorsRegistry");
const DepositContractMock = artifacts.require("DepositContractMock");
const Lido = artifacts.require("Lido");

module.exports = function (deployer) {
  var operatorsAddress = deployer.deploy(NodeOperatorsRegistry);
  var depositAddress = deployer.deploy(DepositContractMock);
  console.log(address);
  var lidoAddress = deployer.deploy(Lido,
    depositAddress,
    `0x0000000000000000000000000000000000000001`,
    operatorsAddress,
    `0x0000000000000000000000000000000000000002`,
    `0x0000000000000000000000000000000000000003`,
    );
  deployer.deploy(
    LidoMevDistributor,
    lidoAddress,
    `${4e18}`
  );
  deployer.deploy(LightPrism);
};
