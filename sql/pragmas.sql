-- Use Write-Ahead Logging for better concurrency and crash recovery
PRAGMA journal_mode = WAL;

-- Balance data safety and performance (fsync at critical moments only)
PRAGMA synchronous = NORMAL;

-- Enforce referential integrity constraints
PRAGMA foreign_keys = ON;

-- Store temporary tables and indices in memory for better performance
PRAGMA temp_store = MEMORY;

-- Limit WAL file size to 64 megabytes to prevent excessive disk usage
PRAGMA journal_size_limit = 67108864;

-- Allocate 128 megabytes for memory-mapped I/O to improve read performance
PRAGMA mmap_size = 134217728;

-- Set database cache to 2 megabytes to reduce disk I/O
PRAGMA cache_size = -2000;

-- Wait up to 5 seconds when database is locked before returning an error
PRAGMA busy_timeout = 5000;
