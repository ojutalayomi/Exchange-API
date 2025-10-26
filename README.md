# Step 2 API - Countries Data Management

A RESTful API built with Go and Gin that manages country data, including population, currencies, exchange rates, and estimated GDP calculations.

## Features

- 🌍 **Country Data Management**: Fetch, store, and manage country information
- 💱 **Exchange Rate Integration**: Real-time exchange rate data from external APIs
- 📊 **GDP Estimation**: Calculate estimated GDP using population and exchange rates
- 🖼️ **Summary Image Generation**: Automatic generation of summary statistics images
- 🚀 **High Performance**: Optimized database operations with batch processing
- 🔄 **Data Refresh**: Bulk refresh capabilities with update/insert logic

## API Endpoints

### Countries Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/countries/refresh` | Refresh all countries data from external APIs |
| `GET` | `/countries` | Get all countries from database |
| `GET` | `/countries/:name` | Get specific country by name |
| `DELETE` | `/countries/:name` | Delete specific country |

### Statistics & Images

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/stats` | Get database statistics |
| `GET` | `/countries/image` | Get generated summary image |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Basic health check |

## Data Structure

### Country Response
```json
{
  "name": "United States",
  "capital": "Washington, D.C.",
  "region": "Americas",
  "population": 331002651,
  "currency_code": "USD",
  "exchange_rate": 1.0,
  "estimated_gdp": 2.1e13,
  "flag_url": "https://flagcdn.com/us.svg",
  "last_refreshed_at": "2025-01-26T10:30:00Z"
}
```

### Nullable Fields
- `currency_code`: `null` if country has no currencies
- `exchange_rate`: `null` if currency not found in exchange API
- `estimated_gdp`: `null` if exchange rate unavailable

## Installation & Setup

### Prerequisites
- Go 1.24+
- MySQL database
- External API access (Countries API & Exchange Rates API)

### Environment Variables
Create a `.env` file in the project root:

```env
MYSQL_STRING=username:password@tcp(host:port)/database?tls=skip-verify
COUNTRIES_API_URL=https://api.example.com/countries
EXCHANGE_RATES_API_URL=https://api.example.com/exchange-rates
```

### Database Setup
The application automatically creates the required table:

```sql
CREATE TABLE countries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    capital VARCHAR(255),
    region VARCHAR(255),
    population BIGINT NOT NULL,
    currency_code VARCHAR(10),
    exchange_rate DOUBLE,
    estimated_gdp DECIMAL(20,1),
    flag_url VARCHAR(512),
    last_refreshed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Running the Application

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Run the application:**
   ```bash
   go run main.go
   ```

3. **The server will start on port 8080:**
   ```
   🚀 Step 2 API server starting on port 8080
   📝 API Documentation available at: /
   🏥 Health check available at: /health
   🔗 Me endpoint: GET /me
   ```

## Usage Examples

### Refresh Countries Data
```bash
curl -X POST http://localhost:8080/countries/refresh
```
**Response:** `204 No Content` (successful refresh)

### Get All Countries
```bash
curl http://localhost:8080/countries
```

### Get Specific Country
```bash
curl http://localhost:8080/countries/United%20States
```

### Get Statistics
```bash
curl http://localhost:8080/stats
```
**Response:**
```json
{
  "total_countries": 195,
  "last_refreshed_at": "2025-01-26T10:30:00Z"
}
```

### Get Summary Image
```bash
curl http://localhost:8080/countries/image
```
**Response:** PNG image file

## Performance Optimizations

### Database Operations
- **Batch Processing**: Single query to fetch all existing countries
- **Map Lookup**: O(1) country existence checking instead of individual queries
- **Bulk Operations**: Separate batch insert and update operations

### API Efficiency
- **Single Exchange Rate Call**: Fetch exchange rates once per refresh
- **Asynchronous Image Generation**: Non-blocking image creation
- **Bounds Checking**: Prevent GDP overflow errors

### Expected Performance
- **Before Optimization**: 3+ minutes response time
- **After Optimization**: <30 seconds response time

## Error Handling

### HTTP Status Codes
- `200`: Success
- `204`: No Content (successful refresh)
- `400`: Bad Request (missing parameters)
- `404`: Not Found (country not found)
- `503`: Service Unavailable (external API or database errors)

### Error Response Format
```json
{
  "error": "Error description",
  "details": "Additional error details (optional)"
}
```

## Project Structure

```
step2/
├── main.go              # Main application and routes
├── db/
│   └── mysql.go         # Database operations
├── utils/
│   ├── api.go           # API utilities and data structures
│   └── image.go         # Image generation utilities
├── cache/
│   └── summary.png      # Generated summary image
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
└── README.md            # This file
```

## Dependencies

- **Gin**: Web framework
- **MySQL Driver**: Database connectivity
- **GG**: Image generation library
- **Godotenv**: Environment variable management

## Development

### Adding New Features
1. Add new routes in `main.go`
2. Implement database methods in `db/mysql.go`
3. Add utility functions in `utils/` package
4. Update this README with new endpoints

### Testing
```bash
# Test all endpoints
curl -X POST http://localhost:8080/countries/refresh
curl http://localhost:8080/countries
curl http://localhost:8080/stats
curl http://localhost:8080/countries/image
```

## License

This project is part of the HNG internship program.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

---

**Built with ❤️ using Go and Gin**
