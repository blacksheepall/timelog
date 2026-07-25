package config

type Config struct {
	DevMode bool `yaml:"dev_mode" env-default:"false"`
	App     struct {
		Timezone string `yaml:"timezone" env-default:"Asia/Singapore"`
	} `yaml:"app"`
	Database struct {
		Host string `yaml:"host"`
	} `yaml:"database"`
	Server struct {
		Addr         string   `yaml:"addr" env-default:""`
		Port         int      `yaml:"port" env-required:"true"`
		AllowOrigins []string `yaml:"allow_origins"`
		HTTPSEnabled bool     `yaml:"https_enabled" env-default:"false"`
		CertFile     string   `yaml:"cert_file" env-default:"./certs/cert.pem"`
		KeyFile      string   `yaml:"key_file" env-default:"./certs/key.pem"`
	} `yaml:"server"`
	Log struct {
		Level    string `yaml:"level" env-default:"info"`
		Path     string `yaml:"path" env-default:"app.log"`
		Rotation struct {
			MaxSize    int `yaml:"max_size"`
			MaxBackups int `yaml:"max_backups"`
			MaxAge     int `yaml:"max_age"`
		} `yaml:"rotation"`
		ORMLogLevel int `yaml:"orm_log_level"`
	} `yaml:"log"`
	Passkey struct {
		Enabled      bool     `yaml:"enabled" env-default:"false"`
		RPID         string   `yaml:"rp_id"`
		RPName       string   `yaml:"rp_name"`
		RPOrigins    []string `yaml:"rp_origins"`
		TokenTTL     int      `yaml:"token_ttl" env-default:"86400"`
		TempPassword struct {
			TTL int `yaml:"ttl" env-default:"900"`
		} `yaml:"temp_password"`
	} `yaml:"passkey"`
	MCP struct {
		Enabled    bool   `yaml:"enabled" env:"MCP_ENABLED" env-default:"false"`
		Level      string `yaml:"level" env:"MCP_LEVEL" env-default:"debug"`
		Path       string `yaml:"path" env:"MCP_PATH" env-default:"logs/mcp.log"`
		Transport  string `yaml:"transport" env:"MCP_TRANSPORT" env-default:"stdio"`
		ListenAddr string `yaml:"listen_addr" env:"MCP_LISTEN_ADDR" env-default:":8080"`
		Token      string `yaml:"token" env:"MCP_TOKEN" env-default:""`
	} `yaml:"mcp"`
	Datasources []DatasourceConfig `yaml:"datasources"`
	Test struct {
		Flush bool `yaml:"flush" env-default:"false"`
	} `yaml:"test"`
}

type DatasourceConfig struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Enabled bool                   `yaml:"enabled" env-default:"true"`
	Config  map[string]interface{} `yaml:"config"`
}
