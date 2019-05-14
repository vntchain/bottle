var Migrations = artifacts.require("./contracts/Migrations.c");

module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
