requirejs.config({
    baseUrl: 'js',
    paths: {
    	'angular': '/js/vendor/angular',
    	'jquery': '/js/vendor/jquery-1.8.3.min'
    }
});

requirejs(['userctrl'], function(userctrl) {
	window.UserCtrl = userctrl;
});

