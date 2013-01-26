define(['config', 'angular'], function(config) {
    return {
        servicename: function() {
            return function(input) {
                var name = config.ServiceNames[input];
                if(!name) {
                    return input;
                }
                return name;
            };
        },
    };
});
