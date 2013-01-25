define(['config', 'angular'], function(config) {
	return function($scope, $http, $location) {
		$scope.domains = [];
		$http.get(config.ApiEndpoint + '/domains')
		.success(function(data) {
			$scope.domains = data;
		}).error(function() {
			$location
			.path('/')
			.search('error','Apparently you are not logged in');
		});
	};
})
