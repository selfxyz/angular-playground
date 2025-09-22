# Angular Playground - Self Protocol API

This Angular project provides API endpoints identical to the Next.js playground, implementing the Self protocol verification system.

## Project Structure

```
angular-playground/
├── src/                    # Angular application
├── server/                 # Express API server
│   ├── api/
│   │   ├── verify.ts      # Production verify endpoint
│   │   ├── saveOptions.ts # Production saveOptions endpoint
│   │   ├── mock-verify.ts # Mock verify endpoint
│   │   └── mock-saveOptions.ts # Mock saveOptions endpoint
│   ├── server.ts          # Production server
│   └── mock-server.ts     # Mock server for testing
└── package.json
```

## Setup

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Environment Variables:**
   Copy `env.example` to `.env` and configure:
   ```bash
   KV_REST_API_URL=your_upstash_redis_url
   KV_REST_API_TOKEN=your_upstash_redis_token
   PORT=3001
   ```

## Available Scripts

### Angular Application
- `npm start` - Start Angular development server (port 4200)
- `npm run build` - Build Angular application
- `npm test` - Run Angular tests

### API Servers
- `npm run server` - Start production API server (port 3001)
- `npm run server:dev` - Start production API server with auto-reload
- `npm run server:mock` - Start mock API server (port 3002)
- `npm run server:mock:dev` - Start mock API server with auto-reload

## API Endpoints

### Production Endpoints (Port 3001)
These endpoints use the actual Self protocol and require valid Redis configuration:

- `POST /api/verify` - Verify Self protocol proofs
- `POST /api/saveOptions` - Save verification options to Redis
- `GET /health` - Health check

### Mock Endpoints (Port 3002)
These endpoints provide mock responses for testing without external dependencies:

- `POST /api/verify` - Mock verification endpoint
- `POST /api/saveOptions` - Mock save options (in-memory storage)
- `GET /api/getOptions/:userId` - Retrieve saved options
- `GET /health` - Health check
- `GET /` - API documentation

## Testing with cURL

### 1. Start the Mock Server
```bash
npm run server:mock
```

### 2. Test Health Check
```bash
curl -X GET http://localhost:3002/health
```

### 3. Test Save Options
```bash
curl -X POST http://localhost:3002/api/saveOptions \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "test-user-123",
    "options": {
      "minimumAge": 18,
      "ofac": true,
      "excludedCountries": ["USA", "CAN"]
    }
  }'
```

### 4. Test Verify Endpoint
```bash
curl -X POST http://localhost:3002/api/verify \
  -H "Content-Type: application/json" \
  -d '{
    "attestationId": "test-id",
    "proof": "test-proof",
    "publicSignals": ["signal1"],
    "userContextData": {"test": "data"}
  }'
```

### 5. Test Get Options
```bash
curl -X GET http://localhost:3002/api/getOptions/test-user-123
```

## Production Usage

### 1. Configure Environment
Set up your Upstash Redis credentials in `.env`:
```bash
KV_REST_API_URL=https://your-redis-instance.upstash.io
KV_REST_API_TOKEN=your-actual-token
```

### 2. Start Production Server
```bash
npm run server
```

### 3. Test Production Endpoints
```bash
# Health check
curl -X GET http://localhost:3001/health

# Save options (requires valid Redis)
curl -X POST http://localhost:3001/api/saveOptions \
  -H "Content-Type: application/json" \
  -d '{"userId": "user123", "options": {...}}'

# Verify (requires valid Self protocol data)
curl -X POST http://localhost:3001/api/verify \
  -H "Content-Type: application/json" \
  -d '{"attestationId": "...", "proof": "...", "publicSignals": [...], "userContextData": {...}}'
```

## Dependencies

### Core Dependencies
- `@selfxyz/core` - Self protocol core SDK
- `@selfxyz/qrcode` - Self protocol QR code SDK
- `@upstash/redis` - Redis client for Upstash
- `express` - Web server framework
- `cors` - CORS middleware

### Development Dependencies
- `ts-node` - TypeScript execution
- `nodemon` - Auto-reload for development
- `@types/*` - TypeScript definitions

## API Response Examples

### Successful Verify Response
```json
{
  "status": "success",
  "result": true,
  "credentialSubject": {
    "name": "John Doe",
    "dateOfBirth": "1990-01-01",
    "nationality": "USA",
    "issuingState": "California",
    "idNumber": "P123456789",
    "gender": "M",
    "expiryDate": "2030-01-01"
  },
  "verificationOptions": {
    "minimumAge": 18,
    "ofac": true,
    "excludedCountries": ["IRN", "PRK", "SYR"]
  }
}
```

### Successful Save Options Response
```json
{
  "message": "Options saved successfully",
  "userId": "test-user-123",
  "savedAt": "2025-09-21T04:37:19.483Z"
}
```

## Notes

- The production endpoints require valid Self protocol data and Redis configuration
- The mock endpoints are perfect for development and testing
- Both servers implement identical logic to the Next.js playground
- Environment variables can be set inline or via `.env` file
- The mock server includes in-memory storage with 30-minute expiration