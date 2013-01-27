define(['angular', 'domainlistctrl', 'domaindetailsctrl', 'notificationctrl', 'loginctrl'],
	function(angular, domainlistctrl, domaindetailsctrl, notificationctrl, loginctrl) {
		var importalias = angular.module('importalias', []);
		importalias.controller('notificationctrl', notificationctrl);
	    importalias.controller('loginctrl', loginctrl);
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
