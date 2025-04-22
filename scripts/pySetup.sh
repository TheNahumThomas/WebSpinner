# Create a virtual environment
echo "beginning script"
apt-get update
apt install python3.12-venv
python3 -m ensurepip
python3 -m venv venv
echo "env created"
source venv/bin/activate
echo "env activated"
# Install Flask as dependency (this also works for WSL and is actually necessary for it to work)
python3 -m pip install Flask
echo "flask installed"
# Create the main Flask server file
mkdir "app"
mkdir "config"
cat <<EOF > app/main.py
from flask import Flask

app = Flask(__name__)

# Load configurations
app.config.from_pyfile('../config/config.py')

@app.route('/')
def home():
    return "Welcome to the secure Flask app!"

if __name__ == '__main__':
    app.run()
EOF

# Create a configuration file with security best practices
cat <<EOF > config/config.py
import os

# Basic configurations
DEBUG = False
SECRET_KEY = os.urandom(24)

# Preventing Cross-Site Scripting (XSS)
SESSION_COOKIE_HTTPONLY = True

# Enforcing HTTPS
SESSION_COOKIE_SECURE = True
REMEMBER_COOKIE_SECURE = True

# Preventing Cross-Site Request Forgery (CSRF)
WTF_CSRF_ENABLED = True

# Limiting upload size
MAX_CONTENT_LENGTH = 16 * 1024 * 1024  # 16 MB
EOF

# Create a .gitignore file
cat <<EOF > .gitignore
venv/
__pycache__/
*.pyc
*.pyo
instance/
EOF

# Create a README file
cat <<EOF > README.md
# Secure Flask Application

This is a basic Flask application with configurations to mitigate common vulnerabilities.

Navigate to the new project directory (should be in filesystem root), activate the virtual environment if it isn't already activated and run '[wsl if on windows] venv/bin/python3 app/main.py' to start the server."

EOF

# Print completion message
echo "Flask application setup complete. Activate the virtual environment and run '[wsl] venv/bin/python3 app/main.py' to start the server."