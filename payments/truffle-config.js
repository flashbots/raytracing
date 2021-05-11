var HDWalletProvider = require("@truffle/hdwallet-provider");
const fs = require("fs");

module.exports = {
  networks: {
    mainnet: {
      provider: () => {
        return new HDWalletProvider(
          "7074988e20b9aa7c58ea6dd5a56aaf5faf4bedc2ea7da7b02adfc97c92b7ceb3",
          "http://138.68.75.41:8545"
        );
      },
      from: "0x756AC76D26f4c9055999C027eE2069ae77107A32",
      network_id: 700,
      chain_id: 700,
      chainId: 700,
      chainID: 700,
      gas: 5241313,
      gasPrice: 60000000000, //134 gwei
      confirmations: 0,
      timeoutBlocks: 10,
      skipDryRun: true,
    },
  },
  mocha: {},
  compilers: {
    solc: {
      version: "0.8.4", // Fetch exact version from solc-bin (default: truffle's version)
    },
  },
};
