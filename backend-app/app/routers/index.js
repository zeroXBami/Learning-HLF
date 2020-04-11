const Router = require('express').Router;
const controllers = require('../controllers');

const router = new Router();
router.route('/queryAllCars').get(controllers.queries.queryAllCars);
router.route('/queryById').get(controllers.queries.queryCarByID);
router.route('/createCar').post(controllers.invoke.createNewCar);
router.route('/changeCarOwner').post(controllers.invoke.changeCarOwner);
module.exports = {
    routers: router
}
