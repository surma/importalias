requirejs.config({
	baseUrl: 'js',
	paths: {
		'angular': '/js/vendor/angular',
		'jquery': '/js/vendor/jquery-1.8.3.min',
		'bootstrap': '/js/vendor/bootstrap',
		'underscore': '/js/vendor/underscore',
	},
	shim: {
		'underscore': {
			exports: '_',
		},
		'bootstrap': {
			deps: ['jquery'],
			exports: '$',
		},
		'angular': {
			deps: ['jquery'],
			exports: 'angular',
		},
	},
});

requirejs(['underscore', 'angular', 'filters', 'validators', 'importalias'], function(_, angular, filters, validators, importalias) {
	_.each(filters, function(filter, name) {
		importalias.filter(name, filter);
	});
	_.each(validators, function(validator, name) {
		importalias.directive(name, validator);
	});
	angular.bootstrap(document, ['importalias']);
});

