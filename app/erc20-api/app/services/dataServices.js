'use strict';
const fs = require('fs');
const path = require('path')
const network = require('../config').network;

const uploadPrivateData = async (dataOwner, dataDesc) => {
    var networkConnection;
    if (dataOwner == "Intage") {
        networkConnection = await network.networkIntage();
    } else if (dataOwner == "WS1") {
        networkConnection = await network.networkWS1();
    }
    const contract = networkConnection.getContract('privateData');
    const result = await contract.submitTransaction('uploadPrivateData', dataOwner, dataDesc);
    const stringResult = result.toString();
    const fileName = `${dataOwner}.${JSON.parse(stringResult).Id}`
    const filePath = path.join(__dirname, '../data', fileName)
    fs.writeFileSync(filePath, stringResult, (err) => {
        if (err) return console.log("Err", err);
    });
    return result;
}

const compareAndPutPrivateData = async (mspOfDataOwner, dataID, dataString, from, to ) => {
    var networkConnection;
    if (from == "Intage") {
        networkConnection = await network.networkIntage();
    } else if (from == "WS1") {
        networkConnection = await network.networkWS1();
    }
    const contract = networkConnection.getContract('privateData');
    const result = await contract.submitTransaction('compareAndPutPrivateData', mspOfDataOwner, dataID, dataString, from, to);
    return result;
}

const getPrivateData = async (dataID, mspid) => {
    const networkConnection = await network.networkWS1();
    const contract = networkConnection.getContract('privateData');
    const result = await contract.evaluateTransaction('getPrivateData', dataID, mspid);
    return result;
}
module.exports = {
    uploadPrivateData: uploadPrivateData,
    compareAndPutPrivateData: compareAndPutPrivateData,
    getPrivateData: getPrivateData
}