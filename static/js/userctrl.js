define(['angular'], function() {
	return function($scope, $http) {
		$scope.userStatus = 'loggedOut';
		$scope.login = function(auth) {
			window.open('/auth/'+auth+'/');
		};

		$scope.refreshAuths = function() {
			$http.get('/auth/')
			.success(function(auths) {
				$scope.auths = auths;
			});
		};

		$scope.refreshUser = function() {
			$http.get('/api/v1/me')
			.then(function() {
				$scope.userStatus = 'loggedIn';
			}, function() {
				$scope.userStatus = 'loggedOut';
			});
		};

		$scope.refresh = function() {
			$scope.refreshAuths();
			$scope.refreshUser();
		}

		window.addEventListener('message', function(event) {
			if(event.data == 'auth_done') {
				$scope.refreshUser();
			}
		}, false);

		$scope.refresh();
	};
});
