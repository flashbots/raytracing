# Ray Tracing
Eth2-MEV project with liquid staking (Flashbots-Lido-Nethermind)

What you need to setup:

* Eth2 validator with Rayonism enabled
* Nethermind Eth1 node with Rayonism and MEV plugins enabled
* A bundle producer that can connect to an Eth1 node and send MEV-bundles
* Set of contracts deployed with Lido MEV distributor contract and Flashbots proxy LightPrism (payments splitter for MEV bundles)

![image](https://user-images.githubusercontent.com/498913/117579537-39e45300-b0eb-11eb-9f66-7fb98e7a923d.png)

Team:
 * Lukasz (Nethermind) - MEV Plugin, The Merge Plugin
 * Edgar (flashbots) - arb spotter demo, contract deployments, MevMasterPayer contract
 * Artem (Lido) - infra, Lido flow
 * Jackie (Nethermind) - full flow design and discussion
 * Tomasz (Nethermind) - LightPrism contract, recording / demo

Aknowledgement:
Thanks to everyone at Flashbots, Lido and Nethermind team, and the Rayonism project. In particular, we are grateful for the support from Tina from Flashbots in fluid cross-team coordination bridging Rayonism and ETHGlobal hackathon, Victor and Sam from Lido team providing insights on Lido contract value distribution, Marek and Mateusz from Nethermind for their participation in the Merge launch, and the work of Flashbots Research on ETH2.
