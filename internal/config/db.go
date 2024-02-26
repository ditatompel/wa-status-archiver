package config

import (
    "os"
    "strconv"
    "time"
)

// DB holds the DB configuration
type DB struct {
    Host            string
    Port            int
    Name            string
    User            string
    Password        string
    MaxOpenConn     int
    MaxIdleConn     int
    MaxConnLifetime time.Duration
}

var db = &DB{}

// DBCfg returns the default DB configuration
func DBCfg() *DB {
    return db
}

// LoadDBCfg loads DB configuration
func LoadDBCfg() {
    db.Host = os.Getenv("DB_HOST")
    db.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
    db.User = os.Getenv("DB_USER")
    db.Password = os.Getenv("DB_PASSWORD")
    db.Name = os.Getenv("DB_NAME")
    db.MaxOpenConn, _ = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))
    db.MaxIdleConn, _ = strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
    lifeTime, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))
    db.MaxConnLifetime = time.Duration(lifeTime) * time.Second
}
