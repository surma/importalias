define(['domainlistctrl', 'domaindetailsctrl', 'newaliasctrl', 'errorctrl', 'angular'],
	function(domainlistctrl, domaindetailsctrl, newaliasctrl, errorctrl) {
		var importalias = angular.module('importalias', []);
		importalias.controller('errorctrl', errorctrl);
		importalias.config(function($locationProvider, $routeProvider) {
			$locationProvider.html5Mode(false);
			$routeProvider
			.when('/', {
				templateUrl: 'partials/home.html'
			})
			.when('/wtf', {
				templateUrl: 'partials/wtf.html'
			})
			.when('/domains', {
				templateUrl: 'partials/domainlist.html',
				controller: domainlistctrl
			})
			.when('/domains/:domain', {
				templateUrl: 'partials/domaindetails.html',
				controller: domaindetailsctrl
			})
			.when('/domains/:domain/new', {
				templateUrl: 'partials/newalias.html',
				controller: newaliasctrl
			})
			.otherwise({redirectTo: '/'});
		});
		return importalias;
	});
