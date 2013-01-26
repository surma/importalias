define(['config', 'angular'], function(config) {

	return function($scope, $http, $location) {
		$scope.user = null;
		$scope.login = function(auth) {
			window.open(config.AuthEndpoint + '/' + auth + '/');
		};
		$scope.logout = function() {
			window.open(config.AuthEndpoint + '/logout');
		};
		$scope.isLoggedIn = function() {
			return $scope.user != null;
		};
		$scope.newApiKey = function() {
			alert("API key regeneration is not yet implemented");
		};

		var refreshAuths = function() {
			$http.get(config.AuthEndpoint + '/')
			.success(function(auths) {
				$scope.auths = auths;
			});
		};
		var refreshUser = function() {
			$http.get(config.ApiEndpoint + '/me')
			.success(function(data) {
				$scope.user = data;
			})
			.error(function() {
				$scope.user = null;
				$location.path('/');
			});
		};
		var refresh = function() {
			refreshAuths();
			refreshUser();
		};

		window.addEventListener('message', function(event) {
			if(event.data == 'auth_done') {
				refreshUser();
			}
		}, false);
		refresh();
	};
});
