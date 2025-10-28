# Backend Challenge - E-Commerce API

A high-performance Go-based e-commerce API server with advanced coupon code validation using optimized file search algorithms.

## 🚀 Features

- **Product Management**: List products with pagination and retrieve individual product details
- **Order Processing**: Create orders with item validation and coupon code support
- ** Coupon Validation**:  coupon code validation using  HDD file reader

## 📋 API Operations

### Health Check
- **GET** `/api/health` - Service health status with timestamp

### Product Operations
- **GET** `/api/product` - List products with pagination
  - Query Parameters: `page` (default: 1), `limit` (default: 5)
- **GET** `/api/product/{productId}` - Get product by ID

### Order Operations
- **POST** `/api/order` - Create new order with optional coupon code validation

## 🧠 HDD File Reader Logic

The application implements a coupon code validation system optimized for large files (~1GB) using the **HDDFileReader**.

### Problem Statement
- Large coupon files (~1GB) cannot be loaded entirely into memory
- Sequential file reading would result in O(N) time complexity in worst-case scenarios
- Need to support concurrent searches across multiple files

### Optimization Strategy

The HDD File Reader implements following optimization approach:

#### 1. **Pre-sorted File Assumption**
- Files are assumed to be sorted in ascending order before application startup
- Enables binary search algorithms for efficient lookups

#### 2. **Partial Indexing System**
- Application reads all file contents during startup
- Creates lightweight indexes for fast file filtering

#### 3. **Chunk-based Indexing**
- For every `chunkSize` lines (configurable, default: 100,000), stores:
  - Line content (coupon code)
  - File offset position

#### 4. **File Range Tracking**
- Captures first and last coupon code for each file
- Enables quick file elimination during search

#### 6. **Binary Search on Partial Index**
- Uses `sort.Search` on chunk keys to locate potential data region
- Reduces search space from entire file to specific chunks

#### 7. **Concurrent Worker Pool Processing**
- Implements worker pool pattern to limit goroutine creation
- Configurable pool size (default: 5 workers)
- Context-based cancellation for early termination when matches found


## 🛠️ Build & Run

### Prerequisites
- Go 1.19 or higher
- Sorted coupon code files in the specified directory

### Environment Variables
```bash
export PORT=8080
export LOG_LEVEL=debug
export ENVIRONMENT=development
export COUPON_CODE_FOLDER_PATH=/path/to/coupon/files
export COUPON_CODE_FILE_PARTIAL_INDEX_CHUNK_SIZE=100000
export COUPON_CODE_FILE_CONCURRENT_POOL_SIZE=5
export GIN_MODE=release
```

### Using Shell Scripts

#### Build the Application
```bash
chmod +x build.sh
./build.sh
```

#### Run the Application
```bash
chmod +x run.sh
./run.sh
```

