/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const services = require('../services');
const contractServices = services;
const dataServices = contractServices.dataServices;

const uploadPRivData = async (request, response) => {
    const {
        body
    } = request;
    const {
        dataOwner,
        dataDesc
    } = body;
    try {
        const result = await dataServices.uploadPrivateData(dataOwner, dataDesc);
        return response.json({
            errorcode: 200,
            message: `Success upload data: ${result.toString()}`,
        });
    } catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to upload data: ${error}`
        });
    }
}

const compareAndPut = async (request, response) => {
    const {
        body
    } = request;
    const {
        mspOfDataOwner,
        dataID,
        dataString,
        from,
        to
    } = body;
    try {
        const result = await dataServices.compareAndPutPrivateData(mspOfDataOwner, dataID, dataString, from, to);
        return response.json({
            errorcode: 200,
            message: `Success to request view data: ${result.toString()}`,
        });
    } catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to request view data: ${error}`
        });
    }
}

const getPrivData = async (request, response) => {
    const {
        body
    } = request;
    const {
        dataID,
        mspid
    } = body;
    try {
        const result = await dataServices.getPrivateData(dataID, mspid);
        return response.json({
            errorcode: 200,
            message: `Success to request view data: ${result.toString()}`,
        });
    } catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to request view data: ${error}`
        });
    }
}


module.exports = {
    uploadPRivData: uploadPRivData,
    compareAndPut: compareAndPut,
    getPrivData: getPrivData
}