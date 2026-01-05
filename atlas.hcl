data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas",
  ]
}

env "local" {
  src = data.external_schema.gorm.url

  // Syntax: docker://image_alias/tag/db_name
  dev = "docker://postgis/16-3.4/dev?search_path=public"

  migration {
    dir = "file://Database/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "production" {
  // For production, we don't use the GORM schema as source
  // Instead, we only apply migrations from the directory

  // Production database URL should be provided via --url flag
  // Example: atlas migrate apply --env production --url "postgres://user:pass@host:5432/dbname?sslmode=require"

  migration {
    dir = "file://Database/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "test" {
  src = data.external_schema.gorm.url

  // Use the same dev database as local
  dev = "docker://postgis/16-3.4/dev?search_path=public"

  migration {
    dir = "file://Database/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
