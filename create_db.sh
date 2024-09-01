#!/bin/bash

sqlite3 database.db <<EOF
-- Enable Write-Ahead Logging for better concurrency and performance
PRAGMA journal_mode = WAL;

-- Enable Foreign Keys for integrity and performance optimization
PRAGMA foreign_keys = ON;

-- Set synchronous mode to NORMAL for a good balance between safety and performance
PRAGMA synchronous = NORMAL;

-- Increase cache size to improve read performance (measured in pages, default page size is 4096 bytes)
PRAGMA cache_size = -2000;

-- Use memory for temporary storage instead of disk
PRAGMA temp_store = MEMORY;

-- Optimize queries by setting automatic indexing
PRAGMA automatic_index = ON;

-- Enable query optimizer options (disabled by default in some versions)
PRAGMA optimize;

-- Add a vacuum to compact the database
VACUUM;

EOF

echo "db created with performance optimizations."
