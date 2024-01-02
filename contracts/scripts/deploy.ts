import {ethers} from "hardhat";

async function main() {
    const contract = await ethers.deployContract("gpn-core");

    await contract.waitForDeployment();
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
