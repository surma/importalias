requirejs.config({
    baseUrl: 'js',
    paths: {
    	'angular': '/js/vendor/angular',
    	'jquery': '/js/vendor/jquery-1.8.3.min'
    }
});

requirejs(['angular'], function() {
	window.UserCtrl = function($scope, $http) {
		$scope.userStatus = "loggedOut";
		$scope.login = function(auth) {
			window.open('/auth/'+auth);
		};

		$scope.refreshAuths = function() {
			$http.get('/auth/')
			.success(function(auths) {
				$scope.auths = auths;
			});
		};

		$scope.refreshUser = function() {
			$http.get("/api/v1/me")
			.then(function() {
				$scope.userStatus = "loggedIn";
			}, function() {
				$scope.userStatus = "loggedOut";
			});
		};

		window.addEventListener("message", function(event) {
			if(event.message == "auth_done") {
				$scope.refreshUser();
			}
		}, false);

		$scope.refreshAuths();

	}
});

