'use strict';

const {
    Gateway,
    Wallets
} = require('fabric-network');
const path = require('path');
const fs = require('fs');
const args = process.argv;
const identityString = args[2];

const networkIntage = async () => {
    const ccpPath = path.resolve(__dirname, '..', '..', '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'intage.example.com', 'connection-intage.json');
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
    const walletPath = path.resolve(__dirname, '..', 'scripts', 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    // Create a new gateway for connecting to our peer node.
    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: identityString,
        discovery: {
            enabled: true,
            asLocalhost: true
        }
    });
    // Get the network (channel) our contract is deployed to.
    const network = await gateway.getNetwork('mychannel');
    return network;
}

const networkWS1 = async () => {
    const ccpPath = path.resolve(__dirname, '..', '..', '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'wholesale1.example.com', 'connection-wholesale1.json');
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
    const walletPath = path.resolve(__dirname, '..', 'scripts', 'walletWS1');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    // Create a new gateway for connecting to our peer node.
    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: 'admin',
        discovery: {
            enabled: true,
            asLocalhost: true
        }
    });
    // Get the network (channel) our contract is deployed to.
    const network = await gateway.getNetwork('mychannel');
    return network;
}
module.exports = {
    networkIntage: networkIntage,
    networkWS1: networkWS1
}