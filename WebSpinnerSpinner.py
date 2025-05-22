import os
import random
import string
import subprocess
import sys

def generate_password(length=16):
    chars = string.ascii_letters + string.digits + string.punctuation
    return ''.join(random.SystemRandom().choice(chars) for _ in range(length))

def install_wsl_distro(distro='Ubuntu'):
    # Install WSL if not present
    subprocess.run(['wsl', '--install', '-d', distro], check=True)

def set_root_password(distro='Ubuntu', password=''):
    # Set root password in the WSL distro
    cmd = f'echo "root:{password}" | wsl -d {distro} sudo chpasswd'
    subprocess.run(cmd, shell=True, check=True)

def main():
    distro = 'Ubuntu'
    print(f"Installing WSL distro: {distro}")
    install_wsl_distro(distro)
    password = generate_password()
    print(f"Generated root password: {password}")
    print("Setting root password in WSL...")
    set_root_password(distro, password)
    print("Done.")


if os.name != 'nt':
    print("This script is intended to run on Windows.")
    sys.exit(1)

main()