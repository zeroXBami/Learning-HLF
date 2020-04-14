/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const services = require('../services');
const contractServices = services.contractServices


const createNewCar = async (request, response) => {
    const { body } = request;
    const { carNumber, make, model, color, owner } = body;
    try {
        const result = await contractServices.createNewCar(carNumber, make, model, color, owner);
        return response.json({
            errorcode: 200,
            message: `Transaction has been submitted, result is: ${result.toString()}`,
        });
    }  catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to submit transaction: ${error}`
        });
    }
}

const changeCarOwner =  async (request, response) => {
    const { body } = request;
    const { id, newOwner } = body;
    try {
        const result = await contractServices.transferCarOwnership(id, newOwner);
        return response.json({
            errorcode: 200,
            message: `Transaction has been submitted, result is: ${result.toString()}`,
        });
    }  catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to submit transaction: ${error}`
        });
    }
} 

const getTokenInfor = async (request, response) => {
    try {
        const result = await contractServices.getTokenInfor();
        return response.json({
            errorcode: 200,
            message: `Token information: ${result.toString()}`,
        });
    }  catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to get token information: ${error}`
        });
    }
}

const getBalanceOf = async (request, response) => {
    const { body } = request;
    const { address } = body;
    try {
        const result = await contractServices.getBalanceOf(address);
        return response.json({
            errorcode: 200,
            message: `Balance of ${address} is: ${result.toString()}`,
        });
    }  catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to get balance: ${error}`
        });
    }
}

const uploadData = async (request, response) => {
    const { body } = request;
    const { dataId, dataOwner } = body;
    try {
        const result = await contractServices.uploadData(dataId, dataOwner);
        return response.json({
            errorcode: 200,
            message: `Success upload data: ${result.toString()}`,
        });
    }  catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to upload data: ${error}`
        });
    }
}

const requestViewData = async (request, response) => {
    const { body } = request;
    const { requester, dataId } = body;
    try {
        const result = await contractServices.requestViewData(requester, dataId);
        return response.json({
            errorcode: 200,
            message: `Success to request view data: ${result.toString()}`,
        });
    }  catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to request view data: ${error}`
        });
    }
}

module.exports = {
    createNewCar: createNewCar,
    changeCarOwner: changeCarOwner,
    getTokenInfor: getTokenInfor,
    getBalanceOf: getBalanceOf,
    uploadData: uploadData,
    requestViewData: requestViewData
}
 