requirejs.config({
    baseUrl: 'js',
    paths: {
    	'angular': '/js/vendor/angular',
    	'jquery': '/js/vendor/jquery-1.8.3.min',
    	'bootstrap': '/js/vendor/bootstrap'
    }
});

requirejs(['jquery', 'filters', 'importalias', 'login', 'bootstrap', 'angular'], function($, filters, importalias) {
	$.each(filters, function(filtername, filter) {
		importalias.filter(filtername, filter);
	})
	angular.bootstrap(document, ['login', 'importalias']);
});

