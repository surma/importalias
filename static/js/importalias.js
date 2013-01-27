define(['domainlistctrl', 'domaindetailsctrl', 'errorctrl', 'angular'],
	function(domainlistctrl, domaindetailsctrl, errorctrl) {
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
			.when('/legal', {
				templateUrl: 'partials/legal.html'
			})
			.when('/domains', {
				templateUrl: 'partials/domainlist.html',
				controller: domainlistctrl
			})
			.when('/domains/:domain', {
				templateUrl: 'partials/domaindetails.html',
				controller: domaindetailsctrl
			})
			.otherwise({redirectTo: '/'});
		});
		return importalias;
	});
