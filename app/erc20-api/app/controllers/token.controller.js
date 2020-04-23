'use strict';

const services = require('../services');
const contractServices = services;
const tokenServices = contractServices.tokenServices;

const mintToken = async (request, response) => {
    const {
        body
    } = request;
    const {
        to,
        amount
    } = body;
    try {
        const result = await tokenServices.mintToken(to, amount);
        return response.json({
            errorcode: 200,
            message: `Success to mint token: ${result.toString()}`,
        });
    } catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to mint token: ${error}`
        });
    }
}

const getTokenInfor = async (request, response) => {
    try {
        // load the network configuration
        const result = await tokenServices.getTokenInfor();
        return response.json({
            errorcode: 200,
            message: `Token name, token symbol, token publisher and total supply is: ${result.toString()}`,
        });

    } catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to get token information: ${error}`
        });
    }
}

const balanceOf = async (request, response) => {
    const {
        body
    } = request;
    const {
        address
    } = body;
    try {
        const result = await tokenServices.getBalanceOf(address);
        return response.json({
            errorcode: 200,
            message: `Balance of ${address} is: ${result.toString()}`,
        });
    } catch (error) {
        return response.json({
            errorcode: 401,
            message: `Failed to get balance of ${address} with error: ${error}`
        });
    }
}

module.exports = {
    getTokenInfor: getTokenInfor,
    balanceOf: balanceOf,
    mintToken: mintToken
}