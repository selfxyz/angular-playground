import 'dotenv/config';
import express from 'express';
import cors from 'cors';
import { verifyHandler } from './api/verify';
import { saveOptionsHandler } from './api/saveOptions';

import dotenv from 'dotenv';
dotenv.config();

const app = express();
const PORT = process.env['PORT'] || 3001;

// Middleware
app.use(cors());
app.use(express.json());

// API Routes
app.post('/api/verify', verifyHandler);
app.post('/api/saveOptions', saveOptionsHandler);

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'OK', message: 'Server is running' });
});

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
