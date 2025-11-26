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
    dir = "file://database/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
