define(['bootstrap', 'config'], function($, config) {
	return function($scope, $routeParams, $http) {
		$scope.domain = $routeParams.domain;
		$scope.aliases = [];
		$scope.alias = {};
		$scope.deleteAlias = function(id) {
			$http.delete(config.ApiEndpoint + '/domains/' + $scope.domain + '/' + id)
			.success(function() {
				window.notify('success', 'Alias deleted');
				refresh();
			})
			.error(function(data) {
				window.notify('error', data);
			});
		}
		$scope.openAliasDialog = function(alias) {
			if(alias) {
				$scope.alias = alias;
			} else {
				$scope.alias = {
					repo_type: "git",
				}
			}
			dialog.modal();
		}
		$scope.saveNewAlias = function() {
			$http.put(config.ApiEndpoint + '/domains/' + $scope.domain, $scope.alias)
			.success(function() {
				window.notify('success', 'Alias added');
				dialog.modal('hide');
				refresh();
			})
			.error(function(data) {
				window.notify('error', data);
			})
		}

		var dialog = $('#aliasdialog');

		var refresh = function() {
			$http.get(config.ApiEndpoint + '/domains/' + $scope.domain)
			.success(function(data) {
				$scope.aliases = data;
			})
			.error(function() {
				window.notify('error', 'Are you not logged in?');
				$location.path('/');
			});
		}
		refresh();
	};
})
