#!/bin/bash

# Create the project directory structure
mkdir -p webapp
cd webapp || exit

# Create the index.php file
cat <<EOL > index.php
<?php
require_once 'config.php';
echo "<h1>Welcome to the PHP WebApp</h1>";
?>
EOL

# Create the config.php file
cat <<EOL > config.php
<?php
// Configuration settings
define('APP_NAME', 'PHP WebApp');
define('APP_ENV', 'production');

// Error reporting settings
if (APP_ENV === 'production') {
    error_reporting(0);
    ini_set('display_errors', '0');
} else {
    error_reporting(E_ALL);
    ini_set('display_errors', '1');
}
?>
EOL

# Create the .htaccess file
cat <<EOL > .htaccess
# Disable directory listing
Options -Indexes

# Protect sensitive files
<FilesMatch "(^\.|config\.php|composer\.json|composer\.lock)">
    Require all denied
</FilesMatch>

# Prevent PHP execution in uploads directory
<Directory "uploads">
    php_flag engine off
</Directory>

# Disable server signature
ServerSignature Off
EOL

# Create a README file
cat <<EOL > README.md
# PHP WebApp

This is a basic PHP web application with security configurations.

## How to Test

1. Start the PHP built-in server:
   \`\`\`
   php -S localhost:8000
   \`\`\`

2. Open your browser and navigate to http://localhost:8000.

## Security Features

- Directory listing is disabled.
- Sensitive files like \`config.php\` are protected.
- PHP execution is disabled in the \`uploads\` directory.
- Server signature is turned off.
EOL

# Print completion message
echo "PHP WebApp setup complete. Navigate to the 'webapp' directory and start the PHP server with 'php -S localhost:8000'."