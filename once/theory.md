## Once

- ``sync.Once`` is a synchronization primitive provided by Go's ``sync`` package.

- It ensures that a piece of code is executed only once regardless of how many goroutines are trying to execute it.

## Anatomy

- ``sync.Once`` is a struct with single method ``Do(f func())``

- It guarantees that the function ``f`` is called at ``most once,`` even if Do is invoked multiple times concurrently.

- Common use cases for ``sync.Once`` include:
    - Initializing shared resources
    - Setting up singletons
    - Performing expensive computations only once
    - Loading configuration files

```go
var instance *singleton
var once sync.Once

func getInstance() *singleton {
    once.Do(func() {
        instance = &singleton{}
    })
    return instance
}
```

## Singleton example

- The Singleton pattern is a classic software design pattern that restricts the instantiation of a class to a single instance

- It's particularly useful when exactly one object is needed to coordinate actions across the system

## Lazy initialization

- Lazy initialization is a design pattern where we delay the creation of an object, the calculation of a value, or some other expensive process until the first time it's needed

- This strategy can significantly improve performance and resource usage, especially for applications with heavy initialization cost

- In Go, ``sync.Once ``provides an excellent mechanism for implementing thread-safe lazy initialization.

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "sync"

    _ "github.com/lib/pq"
)

type DatabaseConnection struct {
    db *sql.DB
}

var (
    dbConn *DatabaseConnection
    once   sync.Once
)

func GetDatabaseConnection() (*DatabaseConnection, error) {
    var initError error
    once.Do(func() {
        fmt.Println("Initializing database connection...")
        db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
        if err != nil {
            initError = fmt.Errorf("failed to open database: %v", err)
            return
        }
        if err = db.Ping(); err != nil {
            initError = fmt.Errorf("failed to ping database: %v", err)
            return
        }
        dbConn = &DatabaseConnection{db: db}
    })
    if initError != nil {
        return nil, initError
    }
    return dbConn, nil
}

func main() {
    // Simulate multiple goroutines trying to get the database connection
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            conn, err := GetDatabaseConnection()
            if err != nil {
                log.Printf("Goroutine %d: Error getting connection: %v\n", id, err)
                return
            }
            log.Printf("Goroutine %d: Got connection %p\n", id, conn)
        }(i)
    }
    wg.Wait()
}
```

## Error handling

If an error occurs during initialization, it's captured and returned to all callers.

## Realistic use-case
This pattern is commonly used in real-world applications for managing shared resources like database connections, configuration loading, or cache initialization.

## Benefits

- **Efficiency**: Resources are allocated only when they're actually needed, which can significantly reduce startup time and memory usage.

- **Thread-safety**: ``sync.Once`` ensures that even if multiple goroutines try to initialize the resource simultaneously, initialization happens exactly once.

- **Simplicity**: Compared to manual locking mechanisms, ``sync.Once`` provides a cleaner and less error-prone approach.

- **Separation of concerns**: The initialization logic is encapsulated within the ``once.Do()`` function, making the code more modular and easier to maintain.

## Gotchas

It's important to note that while ``sync.Once`` is powerful, it's not always the best solution. For instance, if you need the ability to re-initialize a resource (e.g., reconnecting to a database after a connection loss), you might need to use other synchronization primitives like ``mutexes``.


## Real world applications

##  Database Connection Pooling

```go
import (
    "database/sql"
    "sync"

    _ "github.com/lib/pq"
)

var (
    dbPool *sql.DB
    poolOnce sync.Once
)

func GetDBPool() (*sql.DB, error) {
    var err error
    poolOnce.Do(func() {
        dbPool, err = sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
        if err != nil {
            return
        }
        dbPool.SetMaxOpenConns(25)
        dbPool.SetMaxIdleConns(25)
        dbPool.SetConnMaxLifetime(5 * time.Minute)
    })
    if err != nil {
        return nil, err
    }
    return dbPool, nil
}
```

This approach ensures that the database connection pool is initialized only once, regardless of how many goroutines call ``GetDBPool()``. It's both efficient and thread-safe.

## Configuration loading scenario

```go
    import (
        "encoding/json"
        "os"
        "sync"
    )

    type Config struct {
        APIKey string `json:"api_key"`
        Debug bool `json:"debug"`
    }

    var (
        config *Config
        configOnce sync.Once
    )

    func GetConfig() (*Config, error) {
        var err error
        configOnce.Do(func() {
            file, err := os.Open("config.json")
            if err != nil {
                return
            }
            defer file.Close()
            
            config = &Config{}
            err = json.NewDecoder(file).Decode(config)
        })
        if err != nil {
            return nil, err
        }
        return config, nil
}
```

This pattern ensures that the potentially expensive operation of reading and parsing a configuration file happens only once, even if multiple parts of your application request the configuration concurrently.


##  Plugin Initialization in a Modular Go Application

```go
type Plugin struct {
    Name string
    initialized bool
    initOnce sync.Once
}

func (p *Plugin) Initialize() error {
    var err error
    p.initOnce.Do(func() {
        // Simulate complex initialization
        time.Sleep(2 * time.Second)
        if p.Name == "BadPlugin" {
            err = fmt.Errorf("failed to initialize plugin: %s", p.Name)
        }
        p.initialized = true
        fmt.Printf("Plugin %s initialized\n", p.Name)
    })
    return err
}

func UsePlugin(name string) error {
    plugin := &Plugin{Name: name}
    if err := plugin.Initialize(); err != nil {
        return err
    }
    // Use the plugin...
    return nil
}
```

| Usecase   |      Benefits      |  Drawbacks |
|----------|:-------------:|------:|
| DB Connection Pooling	 | Ensures single pool creation, thread-safe	 | May delay error detection until first use|
| Config Loading | Lazy loading, consistent config across app | Might complicate dynamic config updates
|
| Plugin Initialization	 | Efficient for rarely used plugins	 | Could increase complexity in plugin management
 

 - If you need to ``reinitialize`` a resource (e.g., reconnecting to a database after a connection loss), ``sync.Once`` isn't suitable as it only runs ``once``.

## Further reading

[sync.Once](https://cristiancurteanu.com/understanding-go-sync-once/)