@echo off

REM Initialize package.json with default values
echo Initializing package.json...
call npm init -y
IF %ERRORLEVEL% NEQ 0 (
    echo Failed to initialize package.json
    exit /b %ERRORLEVEL%
)

REM Install common dependencies
echo Installing dependencies...
call npm install express dotenv helmet cors config --save
IF %ERRORLEVEL% NEQ 0 (
    echo Failed to install dependencies
    exit /b %ERRORLEVEL%
)

REM Install development dependencies
echo Installing development dependencies...
call npm install nodemon --save-dev
IF %ERRORLEVEL% NEQ 0 (
    echo Failed to install development dependencies
    exit /b %ERRORLEVEL%
)

REM Create necessary directories
echo Creating directories...
mkdir src
mkdir logs
mkdir config

REM Create a secure configuration file
echo Creating configuration file...
echo { > config\default.json
echo     "port": 3000, >> config\default.json
echo     "env": "development", >> config\default.json
echo     "debug": false, >> config\default.json
echo     "security": { >> config\default.json
echo         "cors": { >> config\default.json
echo             "enabled": true, >> config\default.json
echo             "origin": "*" >> config\default.json
echo         }, >> config\default.json
echo         "helmet": { >> config\default.json
echo             "enabled": true, >> config\default.json
echo             "contentSecurityPolicy": true, >> config\default.json
echo             "referrerPolicy": { "policy": "no-referrer" }, >> config\default.json
echo             "xssFilter": true, >> config\default.json
echo             "noSniff": true >> config\default.json
echo         } >> config\default.json
echo     }, >> config\default.json
echo     "rateLimit": { >> config\default.json
echo         "enabled": true, >> config\default.json
echo         "windowMs": 60000, >> config\default.json
echo         "maxRequests": 100 >> config\default.json
echo     }, >> config\default.json
echo     "logging": { >> config\default.json
echo         "enabled": true, >> config\default.json
echo         "level": "info" >> config\default.json
echo     } >> config\default.json
echo } >> config\default.json

REM Create the server.js file
echo Creating server.js file...
echo const express = require('express'); > src\server.js
echo const helmet = require('helmet'); >> src\server.js
echo const cors = require('cors'); >> src\server.js
echo const dotenv = require('dotenv'); >> src\server.js
echo const config = require('config'); >> src\server.js
echo. >> src\server.js
echo dotenv.config(); >> src\server.js
echo. >> src\server.js
echo const app = express(); >> src\server.js
echo. >> src\server.js
echo // Middleware >> src\server.js
echo if (config.get('security.helmet.enabled')) app.use(helmet()); >> src\server.js
echo if (config.get('security.cors.enabled')) app.use(cors({ origin: config.get('security.cors.origin') })); >> src\server.js
echo. >> src\server.js
echo app.use(express.json()); >> src\server.js
echo. >> src\server.js
echo // Basic route >> src\server.js
echo app.get('/', (req, res) =^> { >> src\server.js
echo     res.send('Hello World!'); >> src\server.js
echo }); >> src\server.js
echo. >> src\server.js
echo // Start server >> src\server.js
echo const PORT = process.env.PORT; >> src\server.js
echo app.listen(PORT, () =^> { >> src\server.js
echo     console.log(`Server is running on port ${PORT}`); >> src\server.js
echo }); >> src\server.js

REM Create a .env file for sensitive environment variables
echo Creating .env file...
echo PORT=3000 > .env
echo NODE_ENV=development >> .env

REM Add .gitignore file
echo Creating .gitignore file...
echo node_modules/ > .gitignore
echo .env >> .gitignore
echo logs/ >> .gitignore

echo Setup complete!