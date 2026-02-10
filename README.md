<h1 align="center" style="border-bottom: none;">Hanzo VM</h1>
<h3 align="center">An open-source cloud operating system management platform developed by Go and React.</h3>
<p align="center">
  <a href="#badge">
    <img alt="semantic-release" src="https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg">
  </a>
  <a href="https://github.com/hanzoai/vm/releases/latest">
    <img alt="GitHub Release" src="https://img.shields.io/github/v/release/hanzoai/vm.svg">
  </a>
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/hanzoai/vm">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/hanzoai/vm?style=flat-square">
  </a>
  <a href="https://github.com/hanzoai/vm/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/hanzoai/vm?style=flat-square" alt="license">
  </a>
  <a href="https://github.com/hanzoai/vm/issues">
    <img alt="GitHub issues" src="https://img.shields.io/github/issues/hanzoai/vm?style=flat-square">
  </a>
  <a href="https://github.com/hanzoai/vm/stargazers">
    <img alt="GitHub stars" src="https://img.shields.io/github/stars/hanzoai/vm?style=flat-square">
  </a>
  <a href="https://github.com/hanzoai/vm/network">
    <img alt="GitHub forks" src="https://img.shields.io/github/forks/hanzoai/vm?style=flat-square">
  </a>
</p>

## Architecture

Hanzo VM contains 2 parts:

| Name     | Description                       | Language               | Source code                                       |
|----------|-----------------------------------|------------------------|---------------------------------------------------|
| Frontend | Web frontend UI for Hanzo VM      | Javascript + React     | https://github.com/hanzoai/vm/tree/master/web     |
| Backend  | RESTful API backend for Hanzo VM  | Golang + Beego + MySQL | https://github.com/hanzoai/vm                     |

## Installation

Hanzo VM uses Casdoor as the authentication system. So you need to create an organization and an application for Hanzo VM in a Casdoor instance.

### Necessary configuration

#### Get the code

```shell
go get github.com/casdoor/casdoor
go get github.com/hanzoai/vm
```

or

```shell
git clone https://github.com/casdoor/casdoor
git clone https://github.com/hanzoai/vm
```

#### Setup database

Hanzo VM will store its users, nodes and topics information in a MySQL database named: `hanzo_vm`, will create it if not existed. The DB connection string can be specified at: https://github.com/hanzoai/vm/blob/master/conf/app.conf

```ini
dataSourceName = root:123@tcp(localhost:3306)/
```

Hanzo VM uses XORM to connect to DB, so all DBs supported by XORM can also be used.

#### Configure Casdoor

After creating an organization and an application for Hanzo VM in a Casdoor, you need to update `clientID`, `clientSecret`, `casdoorOrganization` and `casdoorApplication` in app.conf.

#### Run Hanzo VM

- Configure and run Hanzo VM by yourself.
- Open browser: http://localhost:19000/

### Optional configuration

#### Setup your Hanzo VM to enable some third-party login platform

  Hanzo VM uses Casdoor to manage members. If you want to log in with oauth, you should see [casdoor oauth configuration](https://casdoor.org/docs/provider/oauth/overview).

#### OSS, Email, and SMS

  Hanzo VM uses Casdoor to upload files to cloud storage, send Emails and send SMSs. See Casdoor for more details.

#### RDP

Run vmd (Hanzo VM Daemon) for RDP connection.

```shell
docker run --name vmd -d -p 4822:4822 ghcr.io/hanzovm/vmd
```

## Contribute

For Hanzo VM, if you have any questions, you can give Issues, or you can also directly start Pull Requests(but we recommend giving issues first to communicate with the community).

## License

[Apache-2.0](LICENSE)
