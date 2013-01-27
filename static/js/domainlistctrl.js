define(['bootstrap', 'config'], function($, config) {
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
				window.notify('success', 'Domain added');
				refresh();
			})
			.error(function(data) {
				window.notify('error', data);
			})
		}
		$scope.deleteDomain = function(domain) {
			$http.delete(config.ApiEndpoint + '/domains/' + domain)
			.success(function() {
				window.notify('success', 'Domain deleted');
				refresh();
			})
			.error(function(data) {
				window.notify('error', data);
			});
		}

		var dialog = $('#newdomaindialog');

		var refresh = function() {
			$http.get(config.ApiEndpoint + '/domains')
			.success(function(data) {
				$scope.domains = data;
			}).error(function() {
				window.notify('error', 'Are you not logged in?');
				$location.path('/');
			});
		};
		refresh();
	};
})
