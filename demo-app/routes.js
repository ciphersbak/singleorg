//SPDX-License-Identifier: Apache-2.0

var property = require('./controller.js');

module.exports = function(app){

  app.get('/get_property/:id', function(req, res){
    property.get_property(req, res);
  });
  app.get('/add_property/:property', function(req, res){
    property.add_property(req, res);
  });
  app.get('/get_all_property', function(req, res){
    property.get_all_property(req, res);
  });
  app.get('/change_holder/:holder', function(req, res){
    property.change_holder(req, res);
  });
}
