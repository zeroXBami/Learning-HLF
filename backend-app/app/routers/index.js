const Router = require('express').Router;
const controllers = require('../controllers');

const router = new Router();
router.route('/queryAllCars').get(controllers.queries.queryAllCars);

module.exports = {
    routers: router
}
