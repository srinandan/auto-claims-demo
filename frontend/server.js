/**
 * Copyright 2026 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const express = require('express');
const path = require('path');
const rateLimit = require('express-rate-limit');
const { createProxyMiddleware } = require('http-proxy-middleware');
const app = express();

app.set('trust proxy', 1);

const apiBackendUrl = process.env.API_BACKEND_SERVICE_URL || 'http://localhost:8080';

const { GoogleAuth } = require('google-auth-library');

app.use(createProxyMiddleware('/api', {
    target: apiBackendUrl,
    changeOrigin: true,
}));

// Proxy for OTLP traces
const auth = new GoogleAuth({
    scopes: ['https://www.googleapis.com/auth/cloud-platform'],
});
const otlpEndpoint = process.env.OTLP_EXPORTER_URL || 'https://telemetry.googleapis.com';

app.use('/v1/traces', async (req, res, next) => {
    try {
        const client = await auth.getClient();
        const headers = await client.getRequestHeaders();
        for (const key in headers) {
            req.headers[key.toLowerCase()] = headers[key];
        }
    } catch (err) {
        console.error('Error fetching Google Auth headers:', err);
    }
    next();
}, createProxyMiddleware({
    target: otlpEndpoint,
    changeOrigin: true,
}));

const splatLimiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 minutes
    max: 100, // limit each IP to 100 requests per windowMs
});

app.use(express.static(path.join(__dirname, 'dist')));

app.get('/{*splat}', splatLimiter, (req, res) => {
    res.sendFile(path.join(__dirname, 'dist', 'index.html'));
});

const port = process.env.PORT || 8080;
app.listen(port, () => {
    console.log(`Server is running on port ${port}`);
});