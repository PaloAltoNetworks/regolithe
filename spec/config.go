package spec

import ini "gopkg.in/ini.v1"

// Config holds the Specification Config.
type Config struct {
	Author      string
	Copyright   string
	Description string
	Email       string
	Name        string
	Output      string
	ProductName string
	URL         string
	Version     string
}

// NewConfig returns a new APIInfo.
func NewConfig() *Config {
	return &Config{}
}

// LoadConfig loads the config from an ini file.
func LoadConfig(path string) (*Config, error) {

	c := NewConfig()

	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	// Load the sections
	monolitheSection, err := cfg.GetSection("monolithe")
	if err != nil {
		return nil, err
	}

	transformerSection, err := cfg.GetSection("transformer")
	if err != nil {
		return nil, err
	}

	// Set the values
	productNameKey, err := monolitheSection.GetKey("product_name")
	if err != nil {
		return nil, err
	}
	c.ProductName = productNameKey.String()

	copyrightKey, err := monolitheSection.GetKey("copyright")
	if err != nil {
		return nil, err
	}
	c.Copyright = copyrightKey.String()

	outputKey, err := transformerSection.GetKey("output")
	if err != nil {
		return nil, err
	}
	c.Output = outputKey.String()

	nameKey, err := transformerSection.GetKey("name")
	if err != nil {
		return nil, err
	}
	c.Name = nameKey.String()

	urlKey, err := transformerSection.GetKey("url")
	if err != nil {
		return nil, err
	}
	c.URL = urlKey.String()

	authorKey, err := transformerSection.GetKey("author")
	if err != nil {
		return nil, err
	}
	c.Author = authorKey.String()

	emailKey, err := transformerSection.GetKey("email")
	if err != nil {
		return nil, err
	}
	c.Email = emailKey.String()

	descriptionKey, err := transformerSection.GetKey("description")
	if err != nil {
		return nil, err
	}
	c.Description = descriptionKey.String()

	versionKey, err := transformerSection.GetKey("version")
	if err != nil {
		return nil, err
	}
	c.Version = versionKey.String()

	return c, nil
}
