define(['bootstrap', 'config', 'parsley'], function($, config, parsley) {
	return function($scope, $http, $location) {
		$scope.newDomainName = "";
		$scope.domains = [];

		$scope.openNewDomainDialog = function() {
			dialog.modal({
				backdrop: false,
			})
			.on('hidden', function() {
				console.log('hidden');
				$('form').parsley('destroy');
			})
			.on('shown', function() {
				console.log('shown');
				$('form').parsley();
			});
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
