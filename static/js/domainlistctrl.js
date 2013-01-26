define(['jquery', 'config', 'bootstrap', 'angular'], function($, config) {
	return function($scope, $http, $location) {
		$scope.newDomainName = "";
		$scope.domains = [];

		$scope.openNewDomainDialog = function() {
			dialog.modal();
		};
		$scope.saveNewDomain = function() {
			$http.post(config.ApiEndpoint + '/domains/' + $scope.newDomainName)
			.success(function() {
				dialog.modal('hide');
				refresh();
			})
			.error(function(data) {
				console.log('Error: ' + data);
			})
		}
		$scope.deleteDomain = function(domain) {
			$http.delete(config.ApiEndpoint + '/domains/' + domain)
			.success(function() {
				refresh();
			})
			.error(function(data) {
				console.log('Error: ' + data);
			});
		}

		var dialog = $('#newdomaindialog');

		var refresh = function() {
			$http.get(config.ApiEndpoint + '/domains')
			.success(function(data) {
				$scope.domains = data;
			}).error(function() {
				$location
				.path('/')
				.search('error','Apparently you are not logged in');
			});
		};
		refresh();
	};
})
