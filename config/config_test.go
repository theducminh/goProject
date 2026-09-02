package config

import "testing"

func TestLoadConfigUsesEnvValues(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("PORT", "9090")
	t.Setenv("DB_HOST", "postgres-db")
	t.Setenv("DB_USER", "appuser")
	t.Setenv("DB_PASSWORD", "strongpass")
	t.Setenv("DB_NAME", "godb")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("ES_ADDRESS", "https://elasticsearch:9200")
	t.Setenv("ES_PASSWORD", "es-password")
	t.Setenv("REDIS_ADDR", "redis-cache:6379")
	t.Setenv("REDIS_PASSWORD", "redis-password")
	t.Setenv("KAFKA_BROKERS", "kafka:9092")
	t.Setenv("KAFKA_TLS", "true")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}

	if cfg.AppEnv != "production" {
		t.Fatalf("AppEnv = %q, want %q", cfg.AppEnv, "production")
	}
	if cfg.Port != "9090" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.DBHost != "postgres-db" {
		t.Fatalf("DBHost = %q, want %q", cfg.DBHost, "postgres-db")
	}
	if cfg.DBUser != "appuser" {
		t.Fatalf("DBUser = %q, want %q", cfg.DBUser, "appuser")
	}
	if cfg.DBPassword != "strongpass" {
		t.Fatalf("DBPassword = %q, want %q", cfg.DBPassword, "strongpass")
	}
	if cfg.ESAddress != "https://elasticsearch:9200" {
		t.Fatalf("ESAddress = %q, want %q", cfg.ESAddress, "https://elasticsearch:9200")
	}
	if cfg.RedisAddr != "redis-cache:6379" {
		t.Fatalf("RedisAddr = %q, want %q", cfg.RedisAddr, "redis-cache:6379")
	}
}
