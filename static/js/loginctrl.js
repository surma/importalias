define(['config', 'angular'], function(config) {

	return function($scope, $http) {
		$scope.user = null;
		$scope.login = function(auth) {
			window.open(config.AuthEndpoint + '/' + auth + '/');
		};
		$scope.isLoggedIn = function() {
			return $scope.user != null;
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
