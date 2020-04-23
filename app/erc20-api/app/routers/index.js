const Router = require('express').Router;
const controllers = require('../controllers');
const tokenController = controllers.token;
const dataController = controllers.data;
const router = new Router();

router.route('/getTokenInfor').get(tokenController.getTokenInfor);
router.route('/balanceOf').get(tokenController.balanceOf);
router.route('/mintToken').post(tokenController.mintToken);

router.route('/uploadPRivData').post(dataController.uploadPRivData);
router.route('/compareAndPut').post(dataController.compareAndPut);
router.route('/getPrivData').get(dataController.getPrivData);
module.exports = {
    routers: router
}
