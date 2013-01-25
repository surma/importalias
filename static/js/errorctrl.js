define(['angular'], function() {
	return function($scope, $routeParams) {
		console.log($routeParams);
		$scope.error = $routeParams.error;
	}
});
