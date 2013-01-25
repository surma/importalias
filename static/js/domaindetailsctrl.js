define(['config', 'angular'], function(config) {
	return function($scope, $routeParams, $http) {
		$scope.domain = $routeParams.domain;
		$scope.aliases = [];
		$scope.deleteAlias = function(id) {
			$http.delete(config.ApiEndpoint + '/domains/' + $scope.domain + '/' + id)
			.success(refresh);
		}

		var refresh = function() {
			$http.get(config.ApiEndpoint + '/domains/' + $scope.domain)
			.success(function(data) {
				$scope.aliases = data;
			})
			.error(function() {
				console.log('error. Todo: Redirect + error');
			});
		}
		refresh();
	};
})
