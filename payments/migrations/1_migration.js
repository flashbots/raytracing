const LightPrism = artifacts.require("LightPrism");
const LidoMevDistributor = artifacts.require("LidoMevDistributor");
// const NodeOperatorsRegistry = artifacts.require("NodeOperatorsRegistry");
// const DepositContractMock = artifacts.require("DepositContractMock");
// const Lido = artifacts.require("Lido");

module.exports = function (deployer) {
  // var operatorsAddress = deployer.deploy(NodeOperatorsRegistry);
  // var depositAddress = deployer.deploy(DepositContractMock);
  // var lidoAddress = deployer.deploy(Lido,
  //   depositAddress,
  //   `0x0000000000000000000000000000000000000001`,
  //   operatorsAddress,
  //   `0x0000000000000000000000000000000000000002`,
  //   `0x0000000000000000000000000000000000000003`,
  //   );
  // var lidoAddress = '0x5304e3c2b42BEaA8f3e37585bFed1274D1055E47'
  // var address = deployer.deploy(
  //   LidoMevDistributor,
  //   lidoAddress,
  //   `${1e18}`
  // );
  // console.log(address);
  deployer.deploy(LightPrism);
};
