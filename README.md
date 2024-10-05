# Relay server for chat application

- This is relay server for [Lockbox](https://github.com/m1kkY8/lockbox)

## Prerequisites

- Familiarity with linux command line
- Linux server with docker and nginx installed
- TLS certificate from CA, [Let's Encrypt](https://letsencrypt.org/) for example
- Domain from [no-ip](https://www.noip.com/) or any other domain provider

## Installation guide

### Setting up Linux server

_Disclaimer:_ This is not free

For server we need linux machine with static ip for ease of access
Linode or any other cloud service provider will do but this guid will be on how to set up linode instance

#### Chosing where we will host the server

For this we will chose Linode as their 5$/month plan is best value for this

1. Create account on [Linode](https://www.linode.com/), you can also sign in with your Github account
2. Link your credit card or other prefered payment method to your account
3. Buy some credits OR find promo code from streamers and youtubers to get free 100$ for 60 days worth of credit

This is the [referal code](https://linode.com/theprimeagen) i used to get my free 100$

#### Creating the computing instance

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

#### Setting up the instance

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

#### Setting up the domain

For this we will use [no-ip](https://www.noip.com/)

1. Create account on [no-ip](https://www.noip.com/)
2. Add a new host
3. Choose a domain name

#### Creating the SSH key/pair

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

##### Hardening the SSH

1. We will change the default port to avoid bots, use your prefered text editor (I use vim)

2. Find the line `#Port 22` and change it to `Port <port_number>`

3. Disable root login by changing `PermitRootLogin yes` to `PermitRootLogin no`

4. Disable password login by changing `PasswordAuthentication yes` to `PasswordAuthentication no`

5. Disable empty password by changing `PermitEmptyPasswords yes` to `PermitEmptyPasswords no`

6. Disable X11 forwarding by changing `X11Forwarding yes` to `X11Forwarding no`

7. Set `MaxAuthTries` to 1 since we are using ssh keys

8. Change `KbdInteractiveAuthentication yes` to `KbdInteractiveAuthentication no`

9. Change `UsePAM yes` to `UsePAM no`

10. Restart the ssh service

```bash
sudo systemctl restart ssh
```

##### Setting up the firewall

1. Check the status of the firewall

```bash
ufw status
```

2. If the firewall is not active enable it

```bash
ufw enable
```

3. Allow the ssh port

```bash
ufw allow <port_number>
```

4. Allow the http and https ports

```bash
ufw allow 80
ufw allow 443
```

5. Enable the firewall

```bash
ufw enable
```

6. Check the status of the firewall

```bash
ufw status
```

3. Reboot the server

```bash
sudo reboot now
```

After the reboot you can ssh into the server using the new port

```bash
ssh -p <port_number> -i ~/.ssh/<key_name> <your_user>@<ip_address>
```

##### Installing docker

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

##### Installing Nginx

We will need nginx for our server to get the TLS certificate from Let's Encrypt

1. Install Nginx

```bash
sudo apt update && sudo apt install nginx
```

2. Start the Nginx service

```bash
sudo systemctl start nginx
```

##### Getting the TLS certificates

For this we will use certbot

1. Install certbot

```bash
sudo apt install certbot
```

2. Get the certificate

```bash
sudo certbot certonly --standalone -d <your_domain>
```

3. Check if the certificate is installed correctly

```bash
sudo ls /etc/letsencrypt/live/<your_domain>
```

4. We can now disable nginx

```bash
sudo systemctl stop nginx
```

##### Setting up the relay server

1. Compose file for the server is provided in this repository alongside the nginx configuration

2. Clone the repository

```bash
git clone https://github.com/m1kkY8/gochat-relay.git
```

3. Create directory for the server where you want to run docker compose file from and copy contents of the configs folder to that directory

```bash
mkdir -p ~/<server_directory>
cd ~/<server_directory>
```

4. Copy the contents of the configs folder to the server directory

```bash
cp -r ~/gochat-relay/configs/* ~/<server_directory>
```

Note: Server is currently running on docker image i built and pushed to docker hub, you can build the image yourself by running the following command in the server directory

```bash
# Build docker image on your local machine and push it to docker hub
# You need to have docker installed on your local machine
# And will need docker account to push the image
docker build -t <image_name> . --push
```

5. Run the docker compose file

```bash
docker compose pull
docker compose up -d
```

6. Check if the server is running

```bash
docker ps
```

7. Check if the server is accessible

```bash
curl https://<your_domain>/health

```
