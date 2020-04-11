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

module.exports = {
    queryAll: queryAllCars,
    queryById: queryCarByID,
    createNewCar: createNewCar,
    transferCarOwnership: transferCarOwnership
}