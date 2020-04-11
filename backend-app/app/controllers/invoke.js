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
module.exports = {
    createNewCar: createNewCar,
    changeCarOwner: changeCarOwner
}
