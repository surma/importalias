define(['config', 'angular'], function(config) {
	return function($scope, $routeParams, $http) {
		$scope.domain = $routeParams.domain;
	};
})
