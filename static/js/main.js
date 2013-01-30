requirejs.config({
    baseUrl: 'js',
    paths: {
    	'angular': '/js/vendor/angular',
    	'jquery': '/js/vendor/jquery-1.8.3.min',
    	'bootstrap': '/js/vendor/bootstrap',
        'underscore': '/js/vendor/underscore',
        'parsley': '/js/vendor/parsley.min',
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
        'parsley': {
            deps: ['jquery'],
        }
    },
});

requirejs(['underscore', 'angular', 'filters', 'importalias'], function(_, angular, filters, importalias) {
	_.each(filters, function(filter, name) {
		importalias.filter(name, filter);
	});
	angular.bootstrap(document, ['importalias']);
});

