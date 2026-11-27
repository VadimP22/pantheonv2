package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	sqlitemigrate "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	ctrlloc "backend/controller/item_localizations"
	ctrlresources "backend/controller/stockpile_resources"
	ctrlstockpiles "backend/controller/stockpiles"
	ctrltasks "backend/controller/tasks"
	"backend/lib/langseed"
	repoloc "backend/repository/item_localizations"
	reporesources "backend/repository/stockpile_resources"
	repostockpiles "backend/repository/stockpiles"
	repostasks "backend/repository/tasks"
	svcloc "backend/services/item_localizations"
	svcresources "backend/services/stockpile_resources"
	svcstockpiles "backend/services/stockpiles"
	svctasks "backend/services/tasks"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
}

func run() error {
	// Путь к файлу базы данных SQLite.
	dbPath := "backend.db"
	if p := os.Getenv("DATABASE_PATH"); p != "" {
		dbPath = p
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping sqlite: %w", err)
	}

	driver, err := sqlitemigrate.WithInstance(db, &sqlitemigrate.Config{})
	if err != nil {
		return fmt.Errorf("migrate database driver: %w", err)
	}

	// Абсолютный путь к каталогу миграций (file:// схема требует путь или file-url).
	migrationsDir := filepath.Join(".", "migrations")
	if p := os.Getenv("MIGRATIONS_DIR"); p != "" {
		migrationsDir = p
	}
	abs, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("abs migrations dir: %w", err)
	}

	sourceURL := "file://" + filepath.ToSlash(abs)
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("init migrations: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	ver, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("get migrations version: %w", err)
	}
	log.Printf("migrations up to date (version=%d dirty=%v)", ver, dirty)

	// Заполнение справочных таблиц из данных, встроенных в бинарник.
	if err := langseed.Run(db); err != nil {
		return fmt.Errorf("seed data: %w", err)
	}

	// HTTP-сервер: репозиторий -> сервис -> контроллер.
	repo := repostockpiles.New(db)
	svc := svcstockpiles.New(repo)
	ctrl := ctrlstockpiles.New(svc)

	repoResources := reporesources.New(db)
	svcResources := svcresources.New(repoResources, repo)
	ctrlResources := ctrlresources.New(svcResources)

	repoLoc := repoloc.New(db)
	svcLoc := svcloc.New(repoLoc)
	ctrlLoc := ctrlloc.New(svcLoc)

	repoTasks := repostasks.New(db)
	svcTasks := svctasks.New(repoTasks, svcResources)
	ctrlTasks := ctrltasks.New(svcTasks)

	mux := http.NewServeMux()
	ctrl.Register(mux)
	ctrlResources.Register(mux)
	ctrlLoc.Register(mux)
	ctrlTasks.Register(mux)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("http server listening on %s", addr)
	return server.ListenAndServe()
}
