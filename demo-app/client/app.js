// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();

	$scope.queryAllProperty = function(){

		appFactory.queryAllProperty(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_property = array;
		});
	}

	$scope.queryProperty = function(){

		var id = $scope.property_id;

		appFactory.queryProperty(id, function(data){
			$scope.query_property = data;

			if ($scope.query_property == "Could not locate property"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	$scope.recordProperty = function(){

		appFactory.recordProperty($scope.property, function(data){
			$scope.create_property = data;
			$("#success_create").show();
		});
	}

	$scope.changeHolder = function(){

		appFactory.changeHolder($scope.holder, function(data){
			$scope.change_holder = data;
			if ($scope.change_holder == "Error: no property found"){
				$("#error_holder").show();
				$("#success_holder").hide();
			} else{
				$("#success_holder").show();
				$("#error_holder").hide();
			}
		});
	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllProperty = function(callback){

    	$http.get('/get_all_property/').success(function(output){
			callback(output)
		});
	}

	factory.queryProperty = function(id, callback){
    	$http.get('/get_property/'+id).success(function(output){
			callback(output)
		});
	}

	factory.recordProperty = function(data, callback){

		data.location = data.longitude + ", "+ data.latitude;

		var property = data.id + "-" + data.location + "-" + data.timestamp + "-" + data.holder + "-" + data.propertyname;

    	$http.get('/add_property/'+property).success(function(output){
			callback(output)
		});
	}

	factory.changeHolder = function(data, callback){

		var holder = data.id + "-" + data.name;

    	$http.get('/change_holder/'+holder).success(function(output){
			callback(output)
		});
	}

	return factory;
});

