define(['config', 'angular'], function(config) {
	return function($scope, $routeParams, $http) {
		$scope.domain = $routeParams.domain;
		$scope.aliases = [];
		$http.get(config.ApiEndpoint + '/domains/' + $scope.domain)
		.success(function(data) {
			$scope.aliases = data;
		})
		.error(function() {
			console.log('error. Todo: Redirect + error');
		});
	};
})
