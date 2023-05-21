package config

// type ConfigFieldLocation struct { // TODO
// 	Line   int
// 	Column int
// }

// type ConfigField[T interface{}] struct { // TODO
// 	ConfigFieldLocation
// 	Value T
// }

// func (*ConfigField[string]) UnmarshallYAML() {

// }

type MatchProfileConfig struct {
	Condition ConditionConfig `yaml:"condition"`
}

type OidcProfileConfig struct {
	OidcDiscoveryUrlTmpl string `yaml:"oidcDiscoveryURL"`
	OidcClientIdTmpl     string `yaml:"oidcClientId"`
	OidcSecretTmpl       string `yaml:"oidcSecret"`
}

type AuthorizationProfileConfig struct {
	Condition ConditionConfig `yaml:"condition"`
	// Condition struct { // TODO
	// 	ConfigFieldLocation
	// 	ConditionConfig
	// } `yaml:"condition"`
}

type RuleConfig struct {
	MatchProfileTmpl         string `yaml:"matchProfile"`
	OidcProfileTmpl          string `yaml:"oidcProfile"`
	AuthorizationProfileTmpl string `yaml:"authorizationProfile"`
}

type RootConfig struct {
	Http struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
		VHost   string `yaml:"vhost"`
	} `yaml:"http"`
	MatchProfiles         map[string]MatchProfileConfig         `yaml:"matchProfiles"`
	OidcProfiles          map[string]OidcProfileConfig          `yaml:"oidcProfiles"`
	AuthorizationProfiles map[string]AuthorizationProfileConfig `yaml:"authorizationProfiles"`
	Rules                 []RuleConfig                          `yaml:"rules"`
}
