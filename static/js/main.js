requirejs.config({
    baseUrl: 'js',
    paths: {
    	'angular': '/js/vendor/angular',
    	'jquery': '/js/vendor/jquery-1.8.3.min'
    }
});

requirejs(['jquery', 'importalias', 'login', 'angular'], function($) {
	angular.bootstrap($('#login').get(), ['login']);
	angular.bootstrap($('#main').get(), ['importalias']);
});

