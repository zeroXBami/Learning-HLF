/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const services = require('../services');
const contractServices = services.contractServices


const queryAllCars = async (request, response) => {
    try {
        // load the network configuration
       const result = await contractServices.queryAll();

        // Evaluate the specified transaction.
        // queryCar transaction - requires 1 argument, ex: ('queryCar', 'CAR4')
        // queryAllCars transaction - requires no arguments, ex: ('queryAllCars')
        // const result = await contract.evaluateTransaction('queryAllCars');
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        return response.json({
            errorcode: 200,
            message: `Transaction has been evaluated, result is: ${result.toString()}`,
        });

    } catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to evaluate transaction: ${error}`
        });
    }
}

const queryCarByID = async (request, response) => {
    const {body} = request;
    const { id } = body;
    try {
        const result = await contractServices.queryById(id);
        return response.json({
            errorcode: 200,
            message: `Transaction has been evaluated, result is: ${result.toString()}`,
        });
    }  catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to evaluate transaction: ${error}`
        });
    }
}
module.exports = {
    queryAllCars: queryAllCars,
    queryCarByID: queryCarByID
}
