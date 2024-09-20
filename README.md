# Relay server for chat application

- This is relay server for [Gochat](https://github.com/m1kkY8/gochat)

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
