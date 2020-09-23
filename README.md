# Data collector from cryptocurrency exchange written on Go

![go-logo](https://miro.medium.com/max/607/0*VetT99h2ijw0KyZ8.png)

## Installation Go

Download the  [archive](https://golang.org/dl/go1.15.1.linux-amd64.tar.gz "Go1.15.1.linux-amd64.tar.gz") and extract it into /usr/local, creating a Go tree in /usr/local/go

Run the following as root or through sudo:
```bash
$ tar -C /usr/local -xzf go1.14.3.linux-amd64.tar.gz
```
---
Add /usr/local/go/bin to the PATH environment variable

You can do this by adding the following line to your $HOME/.profile or /etc/profile (for a system-wide installation):
```bash
export PATH=$PATH:/usr/local/go/bin
```

**Note**: Changes made to a profile file may not apply until the next time you log into your computer. To apply the changes immediately, just run the shell commands directly or execute them from the profile using a command such as source $HOME/.profile.

---

Verify that you've installed Go by opening a command prompt and typing the following command
```bash
$ go version
```

Confirm that the command prints the installed version of Go.

## installation Tarantool

To install Tarantool, run the following command as a user who has access to sudo:

```bash
# install these utilities if they are missing
apt-get -y install sudo
sudo apt-get -y install gnupg2
sudo apt-get -y install curl

curl https://download.tarantool.org/tarantool/release/2.4/gpgkey | sudo apt-key add -

# install the lsb-release utility and use it to identify the current OS code name;
# alternatively, you can set the OS code name manually, e.g. stretch
sudo apt-get -y install lsb-release
release=`lsb_release -c -s`

# install https download transport for APT
sudo apt-get -y install apt-transport-https

# append two lines to a list of source repositories
sudo rm -f /etc/apt/sources.list.d/*tarantool*.list
echo "deb https://download.tarantool.org/tarantool/release/2.4/ubuntu/ ${release} main" | sudo tee /etc/apt/sources.list.d/tarantool_2_5.list
echo "deb-src https://download.tarantool.org/tarantool/release/2.4/ubuntu/ ${release} main" | sudo tee -a /etc/apt/sources.list.d/tarantool_2_5.list

## install tarantool

sudo apt-get -y update
sudo apt-get -y install tarantool
```

**Note**: Just copy and paste

## How to run project

To start project run the following command:

```bash
$ go run main.go
```