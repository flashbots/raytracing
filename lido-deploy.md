# Deploying Lido

This document descripbes how to deploym the Lido protocol mock and MEV rewards distributor contract.

#### 1. Install the `eth-brownie` Python package

See: https://eth-brownie.readthedocs.io/en/stable/install.html

#### 2. Add Brownie network and local accounts

In a posix-compatible shell, e.g. Bash:

```text
$ brownie networks add Ethereum merge-hackathon host='http://138.68.75.41:8545' chainid=700

$ brownie accounts new merge-hackathon-deployer
<enter the PK>

$ brownie accounts new merge-hackathon-lido-treasury
<enter the PK>

$ brownie accounts new merge-hackathon-op-1
<enter the PK>

$ brownie accounts new merge-hackathon-op-2
<enter the PK>

$ brownie accounts new merge-hackathon-op-3
<enter the PK>

$ brownie accounts new merge-hackathon-st-1
<enter the PK>

$ brownie accounts new merge-hackathon-st-2
<enter the PK>

$ brownie accounts new merge-hackathon-st-3
<enter the PK>

$ brownie console --network merge-hackathon
```

The last command will start a python Brownie console, all of the next steps will be in this console.

#### 3. Load the accounts

```py
>>> deployer = accounts.load('merge-hackathon-deployer')
>>> treasury = accounts.load('merge-hackathon-lido-treasury')
>>> ops = [ accounts.load(f'merge-hackathon-op-{i+1}') for i in range(3) ]
>>> stakers = [ accounts.load(f'merge-hackathon-st-{i+1}') for i in range(3) ]
```

Enter passwords for unlocking the accounts (chosen in the step 2).

#### 4. Populate stakers' accounts with ETH

```py
>>> deployer.transfer(stakers[0], amount='100 ether')
>>> deployer.transfer(stakers[1], amount='100 ether')
>>> deployer.transfer(stakers[2], amount='100 ether')
```

#### 5. Deploy the Lido protocol and MEV distributor

```py
>>> from scripts.deploy import deploy, add_operators, stake
>>> (lido, registry, distributor) = deploy(Lido, NodeOperatorsRegistry, DepositContractMock, LidoMevDistributor, deployer, treasury)
```

#### 6. Add Lido node operators

```py
>>> add_operators(registry, operators, deployer)
```

#### 7. Stake some ETH

```py
>>> stake(lido, [(stakers[0], 1 * 10**18), (stakers[1], 32 * 10**18), (stakers[2], 96 * 10**18)], deployer)
```

This will also pretend beacon balance has increased due to staking rewards.

#### 8. Print stakers' and node operators' stETH balances

```py
>>> [ lido.balanceOf(staker) / 10**18 for staker in stakers ]
>>> [ lido.balanceOf(op) / 10**18 for op in ops ]
```

#### 9. Test MEV distribution

To test the distribution, send some ETH to the MEV distributor contract:

```py
>>> tx = distributor.distribureMev({'from': deployer, 'value': '10 ether'})
>>> tx.info()
```

Print the stETH balances once again to see them increased:

```py
>>> [ lido.balanceOf(staker) / 10**18 for staker in stakers ]
>>> [ lido.balanceOf(op) / 10**18 for op in ops ]
```
