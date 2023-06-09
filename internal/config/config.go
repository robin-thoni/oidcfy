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
	Oidc struct {
		DiscoveryUrlTmpl string `yaml:"discoveryURL"`
		ClientIdTmpl     string `yaml:"clientId"`
		ClientSecretTmpl string `yaml:"clientSecret"`
		ScopesTmpl       string `yaml:"scopes"`
		StandardFlow     struct {
			CallbackTimeoutTmpl string `yaml:"callbackTimeout"`
		} `yaml:"standardFlow"`
	} `yaml:"oidc"`
	LoginTimeoutTmpl string `yaml:"loginTimeout"`
	Cookie           struct {
		DomainTmpl string `yaml:"domain"`
		PathTmpl   string `yaml:"path"`
		SecureTmpl string `yaml:"secure"`
	} `yaml:"cookie"`
}

type AuthorizationProfileConfig struct {
	Condition ConditionConfig `yaml:"condition"`
	// Condition struct { // TODO
	// 	ConfigFieldLocation
	// 	ConditionConfig
	// } `yaml:"condition"`
}

type MutatorProfileConfig struct {
	Mutators []MutatorConfig `yaml:"mutators"`
}

type RuleConfig struct {
	MatchProfileTmpl         string `yaml:"matchProfile"`
	OidcProfileTmpl          string `yaml:"oidcProfile"`
	AuthorizationProfileTmpl string `yaml:"authorizationProfile"`
	MutatorProfileTmpl       string `yaml:"mutatorProfile"`
}

type RootConfig struct {
	Http struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
		BaseUrl string `yaml:"baseUrl"`
	} `yaml:"http"`
	MatchProfiles         map[string]MatchProfileConfig         `yaml:"matchProfiles"`
	OidcProfiles          map[string]OidcProfileConfig          `yaml:"oidcProfiles"`
	AuthorizationProfiles map[string]AuthorizationProfileConfig `yaml:"authorizationProfiles"`
	MutatorProfiles       map[string]MutatorProfileConfig       `yaml:"mutatorProfiles"`
	Rules                 []RuleConfig                          `yaml:"rules"`
}
