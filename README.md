# Generic proving network

This is an implementation of the generic proving network, designed to improve current usage of ZK-provers by:
- Decreasing the cost of the proper decentralization of the provers
- Increasing the throughput of the provers network for any specific consumer(rollup, bridge, etc.)
- Increasing censorship resistance
- Making running a prover more economically viable

## How it works
The general idea is to have a network of provers, that can be used by any consumer. The consumer can be a rollup, a bridge, or any other system that requires a ZK-prover.
Each node inside the network is able to compute the ZK-proofs for the subset of the consumers that it committed to.
This allows to maintain a high level of decentralization, while also keeping low idle time for the provers due to the ability to switch between the consumers.

This is achieved by keeping a list of Docker images of the ZK-provers installed on each node. It is assumed that the prover is able to compute the proof for any consumer it committed to, given the input data.
The selection of the prover is done randomly, using on the latest proof generated for the specific consumer as a random seed.

## Proving workflow
1. The consumer sends a request to the network to generate a proof for the specific data.
2. The network selects a random prover from the list of the provers that committed to the consumer.
3. The selected prover computes the proof and broadcasts it into the network.
4. The members of the network, who are also committed to the consumer, verify the proof and broadcast the verification result.
5. If the verification result is negative, the next prover is selected.
6. If the verification result is positive, the list of signatures from the verifiers is submitted to the EVM blockchain by the proving node.
7. The smart contract verifies the signatures and starts a 24h window, where any other prover can try and submit more signatures of the network members to claim the reward.
8. After the 24h window, the reward is distributed between the prover and the verifiers(to incentivize the validation of others' proofs).

## Development
- Each node of the network has to communicate with the EVM blockchain. To support this, the project uses [abigen](https://geth.ethereum.org/docs/tools/abigen) to generate the bindings for the smart contract. Use `/internal/abi/gen.sh` script to generate the bindings.
- The communication of the consumers with the network is done via gRPC. To update the proto contract, first update the `.proto` file and then run `gen.sh` with proto directory as a first argument and the output directory as a second argument.
