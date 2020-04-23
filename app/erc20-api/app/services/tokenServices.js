'use strict';
const network = require('../config').network;

const getTokenInfor = async () => {
    const networkConnection = await network.networkIntage();
    const contract = networkConnection.getContract('erc20');
    const result = await contract.evaluateTransaction('getTokenInfor','');
    return result; 
}

const getBalanceOf = async (address) => {
    const networkConnection = await network.networkIntage();
    const contract = networkConnection.getContract('erc20');
    const result = await contract.evaluateTransaction('balanceOf', address);
    return result;
}

const mintToken = async (to, amount) => {
    const networkConnection = await network.networkIntage();
    const contract = networkConnection.getContract('erc20');
    const result = await contract.submitTransaction('mint',to, amount);
    return result;
}

module.exports = {
    getTokenInfor: getTokenInfor,
    getBalanceOf: getBalanceOf,
    mintToken: mintToken
}