// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

contract Network {
    uint256 public constant MIN_ETH_AMOUNT_CONSUMER = 1 ether;
    uint256 public constant MIN_ETH_AMOUNT_PROVER = 0.5 ether;
    uint256 private constant SECONDS_IN_DAY = 86400;

    event ProverUpdate(address addr, bool isAdded);

    struct Consumer {
        uint256 balance;
        string containerName;
    }

    struct ConsumerView {
        address addr;
        uint256 balance;
        string containerName;
    }

    struct Prover {
        uint256 balance;
    }

    struct ProvingRewardClaim {
        // TODO: store validators here for validation payouts
        uint16 validations;
    }

    struct ProvingPayout {
        address consumer;
        mapping(address => ProvingRewardClaim) claimers;
        address[] claimersAddresses;
        uint256 claimableAfterTimestamp;
    }

    mapping(address => Consumer) public consumers;
    address[] public consumerAddresses;

    mapping(address => Prover) public provers;
    address[] public proverAddresses;

    mapping(string => ProvingPayout) public payouts;
    string[] public payoutRequestIds;

    modifier willHaveEnoughEth() {
        require(
            msg.value + consumers[msg.sender].balance >=
                MIN_ETH_AMOUNT_CONSUMER,
            "Insufficient amount of Ether sent"
        );
        _;
    }

    function registerConsumer(string calldata _containerName) external payable {
        require(msg.value >= MIN_ETH_AMOUNT_CONSUMER);
        require(bytes(_containerName).length != 0);
        require(bytes(consumers[msg.sender].containerName).length == 0);

        consumers[msg.sender] = Consumer(msg.value, _containerName);
        consumerAddresses.push(msg.sender);
    }

    function depositEth() external payable willHaveEnoughEth {
        require(bytes(consumers[msg.sender].containerName).length != 0);
        consumers[msg.sender].balance += msg.value;
    }

    function withdrawConsumer() external {
        require(consumers[msg.sender].balance != 0);

        uint256 balance = consumers[msg.sender].balance;
        consumers[msg.sender] = Consumer(0, "");

        payable(msg.sender).transfer(balance);
    }

    // registerProver requires a deposit from a prover to economically secure the network
    function registerProver() external payable {
        require(msg.value >= MIN_ETH_AMOUNT_PROVER);
        require(provers[msg.sender].balance == 0);

        provers[msg.sender] = Prover(msg.value);
        proverAddresses.push(msg.sender);

        emit ProverUpdate(msg.sender, true);
    }

    function withdrawRewards() external {
        require(provers[msg.sender].balance != 0);

        uint256 withdrawalAmount = provers[msg.sender].balance - MIN_ETH_AMOUNT_PROVER;
        provers[msg.sender].balance = MIN_ETH_AMOUNT_PROVER; 

        payable(msg.sender).transfer(withdrawalAmount);
    }

    function withdrawProver() external {
        require(provers[msg.sender].balance != 0);

        uint256 balance = provers[msg.sender].balance;
        provers[msg.sender] = Prover(0);

        payable(msg.sender).transfer(balance);

        emit ProverUpdate(msg.sender, false);
    }

    function getConsumers() external view returns (ConsumerView[] memory) {
        uint256 length = consumerAddresses.length;
        ConsumerView[] memory result = new ConsumerView[](length);

        for (uint256 i = 0; i < length; i++) {
            result[i] = ConsumerView(
                consumerAddresses[i],
                consumers[consumerAddresses[i]].balance,
                consumers[consumerAddresses[i]].containerName
            );
        }

        return result;
    }

    function getProvers() external view returns (address[] memory) {
        uint256 length = proverAddresses.length;
        address[] memory result = new address[](length);

        for (uint256 i = 0; i < length; i++) {
            result[i] = proverAddresses[i];
        }

        return result;
    }

    function validationOutputToJson(
        string memory requestId,
        address proverAddress,
        bool isValid
    ) internal pure returns (bytes memory) {
        return
            abi.encodePacked(
                '{"request_id":"',
                requestId,
                '","prover_address":"',
                Strings.toHexString(uint160(proverAddress), 20),
                '","is_valid":',
                isValid ? "true" : "false",
                "}"
            );
    }

    // rs[0], ss[0], vs[0] are consumer parameters of a signature of the request
    // TODO: make it callable more than once to increase the number of validations in case of any malicious actions from
    // other participants of the network, check if no one signed more than one message
    function submitSignedProof(
        string calldata requestId,
        uint256 reward,
        bytes32[] calldata rs,
        bytes32[] calldata ss,
        uint8[] calldata vs
    ) external {
        require(rs.length == ss.length);
        require(vs.length == ss.length);

        address consumer = ecrecover(
            keccak256(abi.encodePacked(requestId, reward)),
            vs[0],
            rs[0],
            ss[0]
        );
        require(consumers[consumer].balance != 0);

        uint16 validationsCnt = 0;
        {
            bytes memory json = validationOutputToJson(
                requestId,
                msg.sender,
                true
            );
            bytes32 hash = keccak256(json);

            for (uint256 i = 1; i < rs.length; ++i) {
                address validator = ecrecover(hash, vs[i], rs[i], ss[i]);
                if (provers[validator].balance != 0) {
                    // TODO: add incentive for validators
                    validationsCnt++;
                }
            }
        }

        payouts[requestId].consumer = consumer;
        payouts[requestId].claimersAddresses.push(msg.sender);
        payouts[requestId].claimableAfterTimestamp =
            block.timestamp +
            SECONDS_IN_DAY;
        payouts[requestId].claimers[msg.sender].validations = validationsCnt;
    }
}
