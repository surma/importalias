define(['jquery', 'config', 'angular'], function($, config) {
	return function($scope, $routeParams, $http) {
		$scope.domain = $routeParams.domain;
		$scope.aliases = [];
		$scope.alias = {};
		$scope.deleteAlias = function(id) {
			$http.delete(config.ApiEndpoint + '/domains/' + $scope.domain + '/' + id)
			.success(refresh);
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
				dialog.modal('hide');
				refresh();
			})
			.error(function(data) {
				console.log('Error: ' + data);
			})
		}

		var dialog = $('#aliasdialog');

		var refresh = function() {
			$http.get(config.ApiEndpoint + '/domains/' + $scope.domain)
			.success(function(data) {
				$scope.aliases = data;
			})
			.error(function() {
				console.log('error. Todo: Redirect + error');
			});
		}
		refresh();
	};
})
