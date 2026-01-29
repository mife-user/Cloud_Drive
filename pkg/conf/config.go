package conf

var globalConfig Config

// 配置结构体
type Config struct {
	Env   string      `mapstructure:"env"`
	Mysql MysqlConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	Gin   GinConfig   `mapstructure:"gin"`
}

// mysql配置结构体
type MysqlConfig struct {
	Dsn     string `mapstructure:"dsn"`
	MaxIdle int    `mapstructure:"max_idle"`
	MaxOpen int    `mapstructure:"max_open"`
}

// redis配置结构体
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// gin配置结构体
type GinConfig struct {
	Mode string     `mapstructure:"mode"`
	Port int        `mapstructure:"port"`
	Cors CorsConfig `mapstructure:"cors"`
}

// CORS配置结构体
type CorsConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
}
