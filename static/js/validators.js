define([], function() {
	var DOMAIN_REGEXP = /^([a-z0-9-]+\.)+[a-z]{2,4}$/i
	var PATH_REGEXP = /^(\/[^\/]+)+$/i
	return {
		'domain': function() {
			return {
				require: 'ngModel',
				link: function(scope, elm, attrs, ctrl) {
					ctrl.$parsers.unshift(function(viewValue) {
						if (DOMAIN_REGEXP.test(viewValue)) {
							console.log('valid');
							ctrl.$setValidity('domain', true);
							return viewValue;
						} else {
							console.log('invalid');
							ctrl.$setValidity('domain', false);
							return undefined;
						}
					});
				}
			};
		},
		'path': function() {
			return {
				require: 'ngModel',
				link: function(scope, elm, attrs, ctrl) {
					ctrl.$parsers.unshift(function(viewValue) {
						if (PATH_REGEXP.test(viewValue)) {
							console.log('valid');
							ctrl.$setValidity('path', true);
							return viewValue;
						} else {
							console.log('invalid');
							ctrl.$setValidity('path', false);
							return undefined;
						}
					});
				}
			};
		},
	};
});
