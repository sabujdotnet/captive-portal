package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	
	Database struct {
		Type     string `yaml:"type"` // mysql or mikrotik
		DSN      string `yaml:"dsn"`
	} `yaml:"database"`
	
	MikroTik struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mikrotik"`
	
	PHPNuxBill struct {
		URL      string `yaml:"url"`
		APIKey   string `yaml:"api_key"`
	} `yaml:"phpnuxbill"`
	
	Radius struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		Secret   string `yaml:"secret"`
	} `yaml:"radius"`
	
	SocialAuth struct {
		Facebook struct {
			ClientID     string `yaml:"client_id"`
			ClientSecret string `yaml:"client_secret"`
			RedirectURL  string `yaml:"redirect_url"`
		} `yaml:"facebook"`
		Google struct {
			ClientID     string `yaml:"client_id"`
			ClientSecret string `yaml:"client_secret"`
			RedirectURL  string `yaml:"redirect_url"`
		} `yaml:"google"`
	} `yaml:"social_auth"`
	
	Branding struct {
		CompanyName string `yaml:"company_name"`
		LogoURL     string `yaml:"logo_url"`
		PrimaryColor string `yaml:"primary_color"`
		SecondaryColor string `yaml:"secondary_color"`
	} `yaml:"branding"`
	
	Features struct {
		VoucherEnabled bool `yaml:"voucher_enabled"`
		SocialLoginEnabled bool `yaml:"social_login_enabled"`
		TermsConditionsEnabled bool `yaml:"terms_conditions_enabled"`
	} `yaml:"features"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	cfg := &Config{}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}
	
	return cfg, nil
}