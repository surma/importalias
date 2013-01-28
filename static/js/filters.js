define(['config'], function(config) {
    return {
        service: function() {
            return function(input) {
                var name = config.Services[input];
                if(!name) {
                    return input;
                }
                return name;
            };
        },
    };
});
