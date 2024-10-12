# Relay server for chat application

- This is relay server for [Lockbox](https://github.com/m1kkY8/lockbox)

# Before start

- Server can use both HTTP and HTTPS

# Prerequisites

- Familiarity with Linux command line
- Linux server with **static IP address** and Docker installed
- This guide will cover only on how to setup server using HTTP. for HTTPS you will need TLS certificate and domain

## HTTPS specifics

- TLS certificate from CA, [Let's Encrypt](https://letsencrypt.org/) for example
- Domain from [no-ip](https://www.noip.com/) or any other domain provider

## Setting up Linux server

**_Disclaimer:_** **This is not free**

Linode or any other cloud service provider will do but this guide will be on how to set up linode instance

## Chosing where we will host the server

For this we will chose Linode as their 5$/month plan is best value for this

1. Create account on [Linode](https://www.linode.com/), you can also sign in with your Github account

### Creating the computing instance

After successful account creation we need to create insance

1. Click on `Create` button on the top left corner
2. Choose `Linode` from the dropdown
3. Choose `Nanode 1GB` as it is the cheapest option
4. Choose operating system that you like (I chose Ubuntu 24.04 LTS)
5. Choose `Region` closest to you
6. Choose `Linode Plan` as `Nanode 1GB`
7. Set `Linode Label` in Details as you like
8. Set `Root Password` in Security
9. Click on `Create Linode`

After this you will be redirected to the dashboard where you can see your newly created instance

### Setting up the instance

1. Click on the newly created instance
2. Click on `Boot` button
3. SSH into the instance using the provided IP address and root password

```bash
ssh root@<ip_address>
```

4. Update the system

```bash
apt update && apt upgrade -y
```

## Creating the SSH key/pair

We will create the ssh key pair to avoid using password for ssh

On your local machine make sure you have ssh installed and .ssh folder in your home directory

```bash
ssh-keygen -t ed25519 -f <key_name> -C "meaningful comment"
```

This will create two files `<key_name>` and `<key_name>.pub` in your .ssh folder

Copy the content of the `<key_name>.pub` file

```bash
cat ~/.ssh/<key_name>.pub
```

Paste the content in the `~/.ssh/authorized_keys` file on the server

```bash
echo "<content>" >> ~/.ssh/authorized_keys
```

Or use `ssh-copy-id` command

```bash
ssh-copy-id -i ~/.ssh/<key_name>.pub root@<ip_address>
```

## Hardening the SSH

```bash
# Use your prefered text editor
vim|nano /etc/ssh/sshd_config
```

1. Find `#Port 22` and change it to `Port <port_number>`, changing default port will prevent a lot of automated bots from cloging up our log file
2. Disable root login by changing `PermitRootLogin yes` to `PermitRootLogin no`
3. Disable password login by changing `PasswordAuthentication yes` to `PasswordAuthentication no`
4. Disable empty password by changing `PermitEmptyPasswords yes` to `PermitEmptyPasswords no`
5. Disable X11 forwarding by changing `X11Forwarding yes` to `X11Forwarding no`
6. Set `MaxAuthTries` to 1 since we are using ssh keys
7. Change `KbdInteractiveAuthentication yes` to `KbdInteractiveAuthentication no`
8. Change `UsePAM yes` to `UsePAM no`
9. Set `ClientAliveInterval 600`
10. Set `ClientAliveCountMax 0`
11. Change `Protocol 2` to use newer version of protocol
12. Restart the ssh service

```bash
sudo systemctl restart ssh
```

## Setting up the firewall

1. Check the status of the firewall

```bash
ufw status
```

2. If the firewall is not active enable it

```bash
ufw enable
```

3. Allow port that we set for SSH

```bash
ufw allow <port_number>
```

4. Check the status of the firewall

```bash
ufw status
```

5. Reboot the server

```bash
sudo reboot now
```

After the reboot we can ssh into the server using the new port and key that we generated

```bash
ssh -p <port_number> -i ~/.ssh/<key_name> <your_user>@<ip_address>
```

## Installing docker

For installation guide you can visit the [official docker documentation](https://docs.docker.com/engine/install/ubuntu/)

1. Add Docker's official GPG key

```bash
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
```

2. Install Docker

```bash
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

3. Add your user to the docker group in order to run docker commands without sudo

```bash
sudo usermod -aG docker $USER
```

4. Reboot the server

```bash
sudo reboot now
```

5. Check if docker is installed correctly

```bash
docker --version
```

```bash
docker run hello-world
```

## Installing the server

1. Clone the repository

```bash
git clone https://github.com/m1kkY8/lockbox-relay.git
```

```bash
docker compose up -d
```

2. Check if the server is running

```bash
docker ps
```

3. Check if the server is accessible

```bash
curl https://<your_ip>/health
```
