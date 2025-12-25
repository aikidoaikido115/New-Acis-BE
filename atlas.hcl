data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./pkg/database/atlas/loader.go",
  ]
}

env "dev" {
  src = data.external_schema.gorm.url
  url = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:${getenv("DB_PORT")}/${getenv("DB_NAME")}?sslmode=${getenv("SSL_Mode")}"
  dev = "docker://postgres/16/dev?search_path=public"
  
  migration {
    dir = "file://migrations"
  }
  
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}