const Router = require('express').Router;
const controllers = require('../controllers');

const router = new Router();
router.route('/queryAllCars').get(controllers.queries.queryAllCars);
router.route('/queryById').get(controllers.queries.queryCarByID);
router.route('/createCar').post(controllers.invoke.createNewCar);
router.route('/changeCarOwner').post(controllers.invoke.changeCarOwner);

router.route('/token/getTokenInfor').get(controllers.invoke.getTokenInfor);
router.route('/token/balanceOf').get(controllers.invoke.getBalanceOf);
router.route('/token/uploadData').post(controllers.invoke.uploadData);
router.route('/token/requestViewData').post(controllers.invoke.requestViewData);

module.exports = {
    routers: router
}
