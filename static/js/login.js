define(['loginctrl', 'angular'], function(loginctrl) {
	var importalias = angular.module('login', []);
	importalias.controller('loginctrl', loginctrl);
	return importalias;
});
