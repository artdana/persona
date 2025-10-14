package persona

type Profile struct {
	Name        string `mapstructure:"name"`
	User 	    string `mapstructure:"user"`
	Email       string `mapstructure:"email"`
	SigningKey  string `mapstructure:"signingKey"`
	Description string `mapstructure:"description"`
}