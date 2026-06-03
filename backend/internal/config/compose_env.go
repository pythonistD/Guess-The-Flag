package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// DeployConfig — параметры Docker/Traefik (секция deploy в config.yml).
// Backend их не использует; нужны только для docker compose.
type DeployConfig struct {
	Domain              string `yaml:"domain"`
	TraefikNetwork      string `yaml:"traefik_network"`
	TraefikEntrypoint   string `yaml:"traefik_entrypoint"`
	TraefikCertResolver string `yaml:"traefik_cert_resolver"`
	TraefikRouterName   string `yaml:"traefik_router_name"`
	ImageTag            string `yaml:"image_tag"`
}

type fileConfig struct {
	Database DBConfig     `yaml:"database"`
	Deploy   DeployConfig `yaml:"deploy"`
}

// ComposeEnv — переменные для docker compose (из config.yml, без отдельного .env).
type ComposeEnv struct {
	PostgresUser        string
	PostgresPassword    string
	PostgresDB          string
	Domain              string
	TraefikNetwork      string
	TraefikEntrypoint   string
	TraefikCertResolver string
	TraefikRouterName   string
	ImageTag            string
}

func LoadComposeEnv(path string) (ComposeEnv, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ComposeEnv{}, fmt.Errorf("read config: %w", err)
	}
	var cfg fileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return ComposeEnv{}, fmt.Errorf("parse config: %w", err)
	}

	db := cfg.Database
	if db.Username == "" {
		return ComposeEnv{}, fmt.Errorf("database.username is required in %s", path)
	}
	if db.Password == "" {
		return ComposeEnv{}, fmt.Errorf("database.password is required in %s", path)
	}
	if db.DbName == "" {
		return ComposeEnv{}, fmt.Errorf("database.dbname is required in %s", path)
	}

	d := cfg.Deploy
	if d.Domain == "" {
		return ComposeEnv{}, fmt.Errorf("deploy.domain is required in %s (for Traefik Host rule)", path)
	}

	return ComposeEnv{
		PostgresUser:        db.Username,
		PostgresPassword:    db.Password,
		PostgresDB:          db.DbName,
		Domain:              d.Domain,
		TraefikNetwork:      defaultIfEmpty(d.TraefikNetwork, "traefik"),
		TraefikEntrypoint:   defaultIfEmpty(d.TraefikEntrypoint, "websecure"),
		TraefikCertResolver: defaultIfEmpty(d.TraefikCertResolver, "letsencrypt"),
		TraefikRouterName:   defaultIfEmpty(d.TraefikRouterName, "guess-the-flag"),
		ImageTag:            defaultIfEmpty(d.ImageTag, "latest"),
	}, nil
}

func (e ComposeEnv) FormatDotenv() string {
	lines := []string{
		fmt.Sprintf("POSTGRES_USER=%s", e.PostgresUser),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", e.PostgresPassword),
		fmt.Sprintf("POSTGRES_DB=%s", e.PostgresDB),
		fmt.Sprintf("DOMAIN=%s", e.Domain),
		fmt.Sprintf("TRAEFIK_NETWORK=%s", e.TraefikNetwork),
		fmt.Sprintf("TRAEFIK_ENTRYPOINT=%s", e.TraefikEntrypoint),
		fmt.Sprintf("TRAEFIK_CERT_RESOLVER=%s", e.TraefikCertResolver),
		fmt.Sprintf("TRAEFIK_ROUTER_NAME=%s", e.TraefikRouterName),
		fmt.Sprintf("IMAGE_TAG=%s", e.ImageTag),
	}
	return strings.Join(lines, "\n") + "\n"
}

func defaultIfEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
