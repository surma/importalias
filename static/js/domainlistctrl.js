define(['config', 'angular'], function(config) {
	return function($scope, $http) {
		$scope.domains = [];
		$http.get(config.ApiEndpoint + '/domains')
		.success(function(data) {
			$scope.domains = data;
		}).error(function() {
			console.log('error. Todo: Redirect + error');
		});
	};
})
