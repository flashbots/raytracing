# Ray Tracing
Eth2-MEV project with liquid staking (Flashbots-Lido-Nethermind)

What you need to setup:

* Eth2 validator with Rayonism enabled
* Nethermind Eth1 node with Rayonism and MEV plugins enabled
* A bundle producer that can connect to an Eth1 node and send MEV-bundles
* Set of contracts deployed with Lido MEV distributor contract and Flashbots proxy LightPrism (payments splitter for MEV bundles)

![image](https://user-images.githubusercontent.com/498913/117579537-39e45300-b0eb-11eb-9f66-7fb98e7a923d.png)

Team:
  Lukasz (Nethermind) - MEV Plugin, The Merge Plugin
  Edgar (flashbots) - arb spotter demo, contract deployments, MevMasterPayer contract
  Artem (Lido) - infra, Lido flow
  Jackie (Nethermind) - full flow design and discussion
  Tomasz (Nethermind) - LightPrism contract, recording / demo

Thanks to everyone at flashbots, Lido and Nethermind who helped the team better understand the flow. Many slides from the presentation were reused from previous demo on the topic by Alex Obadia (prepared on top of the flashbots research). Great thank you to Victor and Sam from Lido for explaining how the Lido contracts distribute value. Special thanks to Tina from flashbots for encouraging and pushing for timely delivery and help with presentation design. Thanks to Marek and Mateusz from Nethermind for their participation in The Merge launch and infra after Luaksz delivered the tonnes of code.
