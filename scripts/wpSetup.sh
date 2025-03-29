#!/usr/bin/env bash

# Make wp-cli executable
sudo apt install php
sudo apt install php-mysql
chmod +x wp-cli.phar
sudo mv wp-cli.phar /usr/local/bin/wp

# Generate secure, random passwords
DB_PASS="$(openssl rand -hex 16)"
ADMIN_PASS="$(openssl rand -hex 16)"

# Download WP core
wp core download --allow-root
# Create MySQL database and user

# Create WP config with random DB password
wp config create \
    --dbname="secure_db" \
    --dbuser="secure_user" \
    --dbpass="$DB_PASS" \
    --skip-check \
    --allow-root

# Secure configurations
wp config set DISALLOW_FILE_EDIT true --raw --allow-root
wp config set table_prefix "wp_secure_" --raw --allow-root
wp config set FORCE_SSL_ADMIN true --raw --allow-root
wp config set XMLRPC_ENABLED false --raw --allow-root

# Set secure authentication keys and salts
wp config set AUTH_KEY "$(openssl rand -hex 32)" --raw --allow-root
wp config set SECURE_AUTH_KEY "$(openssl rand -hex 32)" --raw --allow-root
wp config set LOGGED_IN_KEY "$(openssl rand -hex 32)" --raw --allow-root
wp config set NONCE_KEY "$(openssl rand -hex 32)" --raw --allow-root
wp config set AUTH_SALT "$(openssl rand -hex 32)" --raw --allow-root
wp config set SECURE_AUTH_SALT "$(openssl rand -hex 32)" --raw --allow-root
wp config set LOGGED_IN_SALT "$(openssl rand -hex 32)" --raw --allow-root
wp config set NONCE_SALT "$(openssl rand -hex 32)" --raw --allow-root

# Install WP with randomised, secure admin password
wp core install \
    --url="http://localhost" \
    --title="SecureWP" \
    --admin_user="secureAdmin" \
    --admin_password="$ADMIN_PASS" \
    --admin_email="admin@example.com" \
    --allow-root

# Set proper file permissions
chmod 640 wp-config.php
chmod -R 755 wp-content