define(['underscore', 'config'], function(_, config) {
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
				$scope.auths = _.map(auths, function(authname) {
					return config.Services[authname]
				});
			});
		};
		var refreshUser = function(stateChange) {
			$http.get(config.ApiEndpoint + '/me')
			.success(function(data) {
				$scope.user = data;
				if(stateChange) {
					window.notify('success', 'Logged in');
				}
			})
			.error(function() {
				$scope.user = null;
				if(stateChange) {
					$location.path('/');
				}
			})
			.then(function() {
				$scope.dropDownOpen = false;
			});
		};
		var refresh = function() {
			refreshAuths();
			refreshUser();
		};

		window.addEventListener('message', function(event) {
			if(event.data == 'auth_done') {
				refreshUser(true);
			}
		}, false);
		refresh();
	};
});
