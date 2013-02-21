define([], function() {
	return {
		ApiEndpoint: '/api/v1',
		AuthEndpoint: '/auth',
		Services: {
			"github": {
				id: "github",
				name: "GitHub",
				icon: "glyphicons_401_github.png",
			},
			"google": {
				id: "google",
				name: "Google",
				icon: "glyphicons_382_google_plus.png",
			},
			"facebook": {
				id: "facebook",
				name: "Facebook",
				icon: "glyphicons_410_facebook.png",
			},
			"twitter": {
				id: "twitter",
				name: "Twitter",
				icon: "glyphicons_411_twitter.png",
			},
		},
	}
});
