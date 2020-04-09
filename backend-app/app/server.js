const config = require('./config');
const APIRoutes = require('./routers');
const express = require("express");
const bodyParser = require("body-parser");
const cors = require("cors");

const app = express();

var corsOptions = {
  origin: "http://localhost:8081"
};

app.use(cors(corsOptions));

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use('/v1', APIRoutes.routers);
// simple route


const PORT = process.env.PORT || config.localConfig.HTTP_SERVER_PORT;
app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}.`);
});