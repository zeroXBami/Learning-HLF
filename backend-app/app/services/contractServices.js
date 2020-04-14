'use strict';
const network = require('../config').network;


const queryAllCars = async () => {
    const networkConnection = await network.network();
    const contract = networkConnection.getContract('fabcar');
    const result = await contract.evaluateTransaction('queryAllCars');
    return result;
}

const queryCarByID = async (id) => {
    const networkConnection = await network.network();
    const contract = networkConnection.getContract('fabcar');
    const result = await contract.evaluateTransaction('queryCar', id);
    return result; 
}


const createNewCar = async (carNumber, make, model, color, owner) => {
    const networkConnection = await network.network();
    const contract = networkConnection.getContract('fabcar');
    const result = await contract.submitTransaction('createCar', carNumber, make, model, color, owner);
    return result; 
}

const transferCarOwnership = async (id, newOwner) => {
    const networkConnection = await network.network();
    const contract = networkConnection.getContract('fabcar');
    const result = await contract.submitTransaction('changeCarOwner', id, newOwner);
    return result; 
}

const getTokenInfor = async () => {
    const networkConnection = await network.network();
    const contract = networkConnection.getContract('erc20');
    const result = await contract.submitTransaction('getTokenInfor','');
    return result; 
}

const getBalanceOf = async (address) => {
    const networkConnection = await network.network();
    const contract = networkConnection.getContract('erc20');
    const result = await contract.submitTransaction('balanceOf',address);
    return result;
}

const uploadData = async (dataId, dataOwner) => {
    const networkConnection = await network.network();
    const contract = networkConnection.getContract('erc20');
    const result = await contract.submitTransaction('uploadData',dataId, dataOwner);
    return result;
}

const requestViewData = async (requester, dataId) => {
    const networkConnection = await network.network();
    const contract = networkConnection.getContract('erc20');
    const result = await contract.submitTransaction('requestViewData',requester, dataId);
    return result;
}

module.exports = {
    queryAll: queryAllCars,
    queryById: queryCarByID,
    createNewCar: createNewCar,
    transferCarOwnership: transferCarOwnership,
    getTokenInfor: getTokenInfor,
    getBalanceOf: getBalanceOf,
    uploadData: uploadData,
    requestViewData: requestViewData
}