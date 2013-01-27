define(['underscore'], function(_) {
	return function($scope, $timeout) {
		var notification_buffer = [];
		$scope.notifications = [];

		var refresh = function() {
			notification_buffer = _.filter(notification_buffer, function(elem) { return elem.visible});
			$scope.notifications = notification_buffer;
		};

		window.notify = function(type, message, permanent) {
			var notification = {
				type: type,
				message: message,
				visible: true,
			};
			notification_buffer.push(notification);
			refresh();
			$timeout(function() {
				notification.visible = false;
				refresh();
			}, 5000);
		};
	};
});
