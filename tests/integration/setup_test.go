package integration

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/internal/app"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	_ "modernc.org/sqlite"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

type TestServer struct {
	DB     *gorm.DB
	Route  *chi.Mux
	Server *httptest.Server
}

var TestServerInstance *TestServer
var DBName = "todo-golang-test.db"

func TestMain(m *testing.M) {
	env.LoadEnv(".env.testing")
	TestServerInstance = setupGlobalServer()
	code := m.Run()
	tearDownGlobalServer(TestServerInstance)
	os.Exit(code)
}

func setupGlobalServer() *TestServer {
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite", // <-- must match the imported driver
		DSN:        DBName,
	}, &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate models
	err = db.AutoMigrate(&database.User{}, &database.TokenBlacklist{})
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	appContainer := app.NewAppContainer(db)
	appContainer.RegisterRoutes(r)

	return &TestServer{
		DB:     db,
		Route:  r,
		Server: httptest.NewServer(r),
	}
}

func tearDownGlobalServer(ts *TestServer) {

	// Close the test server
	if ts.Server != nil {
		ts.Server.Close()
	}

	if ts.DB != nil {
		sqlDB, err := ts.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	// Clean up test database file if db is sqlite
	if ts.DB.Dialector.Name() == "sqlite" {
		_ = os.Remove(DBName)
	}

}

func ClearAllTables(t *testing.T, db *gorm.DB) {
	// Get all table names
	tables, err := getAllTableNames(db)
	if err != nil {
		t.Errorf("Failed to get table names: %v", err)
		return
	}

	if len(tables) == 0 {
		t.Log("No tables found to clean")
		return
	}

	//t.Logf("Found %d tables to clean: %s", len(tables), strings.Join(tables, ", "))

	// Disable foreign key constraints
	if err := disableForeignKeys(db); err != nil {
		t.Logf("Warning: Could not disable foreign keys: %v", err)
	}

	// Clear all tables
	for _, table := range tables {
		result := db.Exec(fmt.Sprintf("DELETE FROM %s", quoteTableName(db, table)))
		if result.Error != nil {
			t.Errorf("Failed to clear table %s: %v", table, result.Error)
		}
	}

	// Re-enable foreign key constraints
	if err := enableForeignKeys(db); err != nil {
		t.Logf("Warning: Could not re-enable foreign keys: %v", err)
	}
}

func getAllTableNames(db *gorm.DB) ([]string, error) {
	var tables []string

	switch db.Dialector.Name() {
	case "sqlite":
		query := "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
		err := db.Raw(query).Scan(&tables).Error
		return tables, err

	case "mysql":
		query := "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_TYPE = 'BASE TABLE'"
		err := db.Raw(query).Scan(&tables).Error
		return tables, err

	default:
		return nil, fmt.Errorf("unsupported database type: %s", db.Dialector.Name())
	}
}

func disableForeignKeys(db *gorm.DB) error {
	switch db.Dialector.Name() {
	case "sqlite":
		return db.Exec("PRAGMA foreign_keys = OFF").Error
	case "mysql":
		return db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	default:
		return nil
	}
}

func enableForeignKeys(db *gorm.DB) error {
	switch db.Dialector.Name() {
	case "sqlite":
		return db.Exec("PRAGMA foreign_keys = ON").Error
	case "mysql":
		return db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	default:
		return nil
	}
}

func quoteTableName(db *gorm.DB, tableName string) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "`" + tableName + "`"
	default:
		return tableName
	}
}

func seedUser[T any](t *testing.T, input T) *database.User {
	t.Helper()

	// Defaults
	user := database.User{
		FirstName:           "John",
		LastName:            "Doe",
		Email:               "testing@gmail.com",
		Password:            "password",
		ResetToken:          nil,
		ResetTokenExpiresAt: nil,
	}

	inputValue := reflect.ValueOf(&input).Elem()
	inputType := inputValue.Type()

	if inputType.Kind() != reflect.Struct {
		t.Fatalf("seedUser expects a struct, got %s", inputType.Kind())
	}

	defaultValue := reflect.ValueOf(&user).Elem()
	defaultType := defaultValue.Type()

	for i := 0; i < defaultType.NumField(); i++ {
		field := defaultValue.Field(i)
		fieldType := defaultType.Field(i)

		inField := inputValue.FieldByName(fieldType.Name)
		if inField.IsValid() && !inField.IsZero() && field.CanSet() {
			field.Set(inField)
		}
	}

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	result := TestServerInstance.DB.Create(&user)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	return &user
}

func GenerateTestJwtToken(t *testing.T, userID uint) string {
	ttl := time.Now().Add(time.Hour * time.Duration(env.FetchInt("JWT_TTL")))
	token, err := utils.GenerateJwtToken(env.FetchString("JWT_SECRET"), ttl, userID)
	if err != nil {
		t.Fatal("uable to generate JWT token", err)
	}
	return token
}
