#!/bin/bash

npm init -y


mkdir -p src
mkdir -p logs
mkdir -p config

# Install common dependencies
npm install express dotenv helmet cors --save

# Create a secure configuration file
cat <<EOL > config/default.json
{
    "port": 3000,
    "env": "development",
    "debug": false,
    "security": {
        "cors": {
            "enabled": true,
            "origin": "*"
        },
        "helmet": {
            "enabled": true,
            "contentSecurityPolicy": true,
            "referrerPolicy": { "policy": "no-referrer" },
            "xssFilter": true,
            "noSniff": true
        }
    },
    "rateLimit": {
        "enabled": true,
        "windowMs": 60000,
        "maxRequests": 100
    },
    "logging": {
        "enabled": true,
        "level": "info"
    }
}
EOL

# Create a basic server file
cat <<EOL > src/server.js
const express = require('express');
const helmet = require('helmet');
const cors = require('cors');
const dotenv = require('dotenv');
const config = require('config');

dotenv.config();

const app = express();
const port = process.env.PORT || config.get('port');

// Middleware
if (config.get('security.helmet.enabled')) app.use(helmet());
if (config.get('security.cors.enabled')) app.use(cors({ origin: config.get('security.cors.origin') }));

app.use(express.json());

// Basic route
app.get('/', (req, res) => {
    res.send('Hello, secure Node.js app!');
});

// Start server
app.listen(port, () => {
    console.log(\`Server running on port \${port}\`);
});
EOL

# Create a .env file for sensitive environment variables
cat <<EOL > .env
PORT=3000
NODE_ENV=development
EOL

# Add .gitignore file
cat <<EOL > .gitignore
node_modules/
.env
logs/
EOL
