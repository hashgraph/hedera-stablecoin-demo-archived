# Hedera Stable Coin Demo

> This is a sample implementation of an ERC20-like token in a 
> Hedera Consensus Service (HCS) Decentralized Application.

## Dependencies

 * PostgreSQL 12+
 
 * Go 1.14+ 

 * Node.js v13.14.0 for the UI
    
 * [migrate](https://github.com/golang-migrate/migrate) <sub>†</sub>
 
 * Protobuf Compiler (protoc) <sub>†</sub>
 
 * [`protoc-gen-go`](https://github.com/golang/protobuf) <sub>†</sub> 
 
 * (optional) docker version 19.03.8 (or above)
 * (optional) docker compose version 1.25.4 (or above)

<sup><sub>† Required only for development.</sub></sup>

## Architecture

 * API – `api/`
 * Mirror – `mirror/`
 * Mirror API – `mirror/api/`
 * Mirror State (working memory) – `mirror/state/`
 * Database (persistence) – `data/`

## License

Licensed under Apache License,
Version 2.0 – see [LICENSE](LICENSE) in this repo
or [apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

## Deployment and installation

In this example, all components have been co-located into a single server, it is however possible to separate each component if desired. The server is a n2d-standard-16 Google Compute virtual machine comprising 16 vCPUs and 64Gb of memory with a 2TB SSD hard drive running Debian GNU/Linux 10 (buster).

*Note: This specification is likely oversize, however lower GCP specifications result in lower network throughput and lower disk IOPS* 

## Docker compose

A docker compose file enables you to quickly get the solution up and running. 

*Note: This is not recommended for performance testing where an individual component installation is preferable*

### Setup the application parameters

```shell script
cd ~/hedera-stable-coin-demo

cp .env.sample .env

nano .env
```

complete the following values in the .env file with your own, you may leave the others as they are.

```text
# hedera network operator
OPERATOR_ID=0.0.__
OPERATOR_KEY=302...

# topic used for the token
TOPIC_ID=0.0.__
```
### start docker images

```shell script
cd ~/hedera-stable-coin-demo

docker-compose up
```

All the necessary services should start up, you can navigate to the UI at `http://localhost:8080`

*Note: adding a `-d` flag to the `docker-compose up` command will run the docker images in the background.

## Installing individual components

Please follow the instructions below to install each of the solution's components individually.

### PostgreSQL installation

#### Install wget

```shell script
sudo apt update
sudo apt install wget
```

#### Install the PostgreSQL version 12 software.

```shell script
# Create the file repository configuration:
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'

# Import the repository signing key:
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

# Update the package lists:
sudo apt-get update

# Install the latest version of PostgreSQL.
sudo apt-get install postgresql-12
```

after the installation is complete, connect to the database and issue the following command to update the postgres user’s password

```shell script
sudo -u postgres psql

ALTER USER postgres PASSWORD 'secretpasswordhere';

CTRL+D to exit psql
```

It may be necessary to update the networking permissions by editing the `/etc/postgresql/12/main/pg_hba.conf` file

```shell script
sudo nano /etc/postgresql/12/main/pg_hba.conf
```

And adding the following lines at the end of the file.

```text
hostnossl    all          all            0.0.0.0/0  trust
host         all          all            0.0.0.0/0  trust
host         all          all            all        md5
```

*Note: this fully opens the database to any incoming connection on any IP address, only a firewall if installed protects the database from outside interference on port `5432`. You may wish to be more specific in how this is configured.*

Also edit `/etc/postgresql/12/main/postgresql.conf`

```shell script
# backup the original file
sudo nano cp sudo nano /etc/postgresql/12/main/postgresql.conf sudo nano /etc/postgresql/12/main/postgresql.conf.installed

sudo nano /etc/postgresql/12/main/postgresql.conf
```

with the contents of the `postgresql.conf` file in this repository.

If you prefer to use your own configuration, ensure that the port the PostgreSQL server is listening to is `5432` (port = 5432). 

Finally, restart the database server to ensure all changes are taken into account

```shell script
sudo pg_ctlcluster 12 main restart
```

### Application components

#### Obtain a Hedera Testnet account

Navigate to [https://portal.hedera.com](https://portal.hedera.com) and create a testnet account if you don’t already have one.
Make a note of your account id (0.0.xxxxx) and private key (302…..)

#### Install GoLang

Install GoLang version 1.14.4 following these instructions [how-to-install-go-on-debian-9](https://www.digitalocean.com/community/tutorials/how-to-install-go-on-debian-9)

*Note: Ensure you install version 1.14.4 and not 1.12.5 with this command*

```shell script
curl -O https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
```

#### Install Git

```shell script
sudo apt update

sudo apt install git

git version
```

> Output

```
git version 2.11.0
```

Download the example source code from git

```shell script
git clone https://github.com/hashgraph/hedera-stable-coin-demo.git
```

#### Setup the application parameters

```shell script
cd ~/hedera-stable-coin-demo

cp .env.sample .env

nano .env
```

complete the following values in the .env file


```text
DATABASE_URL=postgres://postgres:password@localhost/stabl

# hedera network 
HEDERA_NETWORK=network/testnet.json
HEDERA_MIRROR_NODE=hcs.testnet.mirrornode.hedera.com:5600
# if true, the start time for the subscription will be epoch time if there is no data in the database
# if false, the start time will be now - 10s if there is no data in the database
MIRROR_CATCHUP=false

# hedera network operator
OPERATOR_ID=0.0.__
OPERATOR_KEY=302...

# topic used for the token
# leave commented out to have the service create one on boot and exit immediately
TOPIC_ID=0.0.__

# stabl coin network issuer (join/mint)
ISSUER_PRIVATE_KEY=302e020100300506032b6570042204207ce25f7ac7a4fa7284efa8453f153922e16ede6004c36778d3870c93d5dfbee5
# stabl coin network admin (freeze, unfreeze, clawback)
ADMIN_PUBLIC_KEY=302a300506032b6570032100ce99793e964002a8d268ee34a14cb1a33effd9054ba1a423e8ac82108ae6a157
#private 302e020100300506032b6570042204205f89127785c690602c383b03202267fea8e9f63ef6ed6a55f2975fd2fc0d2f36

# token metadata
# returned from API GET /v1/token
TOKEN_NAME=STABL
SYMBOL=STB
TOTAL_SUPPLY=0
DECIMALS=0

# API Port
PORT=3128

# Mirror Port
MIRROR_PORT=3129

# Mirror address (for API)
MIRROR_ADDRESS=http://localhost:3129

# Logging
API_LOG=WARN # api/
MIRROR_LOG=WARN # mirror/
MIRROR_API_LOG=WARN # mirror/api/
```

*Note: Ports `3128` and `3129` must be reachable for load testing and the UI. You may change them for `8081` and `8082` respectively (don't forget to reflect the change for the `MIRROR_ADDRESS` parameter).*

#### Run the REST API

```shell script
export GIN_MODE=release
go run ./api
```

> If all is well, the following output should show on screen after go has downloaded the necessary dependencies

```shell script
⇨ http server started on [::]:3128
```

You can now `CTRL+C` to stop the REST api for now.

#### Run the mirror subscriber

```shell script
export GIN_MODE=release
go run ./mirror
```

> The following output should show on screen after go has downloaded the necessary dependencies

```text
Secs, Count, AvgTPS, BlockTPS, BlockSecs (BlockSize 2000)
```

You can now `CTRL+C` to stop the mirror subscriber for now.

#### Install NodeJS for the UI

```shell script
cd ~

curl -sL https://deb.nodesource.com/setup_13.x | sudo -E bash 

sudo apt-get install -y nodejs

node
```

> output

```text
Welcome to Node.js v13.14.0.
Type ".help" for more information.
>
```

Type `.exit` to exit the node command line

#### Configure the UI environemnt

Ensure that `config/dev.env.js` and `config/prod.env.js` (if applicable) are configured properly:

`POST_PORT: '3128'` must match the `PORT=3128` specified in the `.env` file earlier in this document
`GET_PORT: '3129'` must match the `MIRROR_PORT=3129` specified in the `.env` file earlier in this document  
  
#### Install compiler tools

```shell script
sudo apt-get update

sudo apt-get install build-essential
```

#### Install the UI server

```shell script
cd ~/hedera-stable-coin-demo/ui

npm install
```

*Note: Ignore any vulnerability warnings for now.*

#### Run the UI server

```shell script
npm run dev
```

*Note: Ignore UnhandledPromiseRejectionWarning errors*

#### Check UI server is running

Navigate to `http://your-server-ip-address:8081/`, the UI below should appear

 ![Register](assets/images/register.png)

### Testing the UI

#### Start all the components

To test the UI end to end, all components must be up and running, open an ssh connection to your server for each of the commands below

**Terminal 1 – REST API**

```shell script
cd ~/hedera-stable-coin-demo; export GIN_MODE=release; go run ./api
```

**Terminal 2 – Mirror subscriber**

```shell script
cd ~/hedera-stable-coin-demo; export GIN_MODE=release; go run ./mirror
```

**Terminal 3 – UI Server**

```shell script
cd ~/hedera-stable-coin-demo/ui; npm run dev
```

**Web browser (chrome preferred)**

Navigate to `http://your-server-ip-address:8081/`

In the top left menu, click on Explore 0.0.xxxx (This is your topic id). A new tab will open. 

This tab is a view from a third party mirror node ([Kabuto](https://kabuto.sh)) which will show you messages sent to the Hedera Consensus Service.

Note: [Dragonglass](https://testnet.dragonglass.me/hedera/home) is an alternative third party mirror you may also use.

#### Registering as a new user

In the UI below, input a Username and press the REGISTER button.

![Register](assets/images/registerUser.png)

This will generate a private/public key which will be stored in a Cookie and then will send a request including the public key and username to JOIN the decentralized application network to the REST API server which will package and sign a request and send it to the Hedera Consensus service.

Within a few seconds, you should see a new message appear in the [Kabuto](https://kabuto.sh) Explorer. The same message is picked up by the Mirror Subscriber which verifies the message and adds the new user to the database.

At the same time, the UI should refresh and present the accounts page.

*Note: You may reset your registration (this does not remove data from the database, only from the cookies) by clicking on the menu and choosing the Reset menu option.*

#### Accounts

After successful registration and refresh (*Note: if you close the browser and navigate to the same URL, the cookies stored in the browser will indicate registration is complete and the UI will switch to the Accounts view automatically*)

![Register](assets/images/accounts.png)

Cash and Visa are for display purposes only, there is no functionality behind the buttons.
STABL shows the user’s current STABL balance. Clicking on the row enables the user to operate the STABL account.

#### Operating the STABL account

After clicking on the STABL account, the UI presents a screen where it’s possible to operate the STABL account.

![Register](assets/images/operate.png)

This screen enables the end user to BUY STABL tokens in exchange for FIAT currency, REDEEM existing tokens for FIAT currency and SEND STABL tokens to other users of the decentralized application network.

Recent activity shows the last 5 operations on the account.

*Note: `See More` and `REQUEST` have not been implemented at the time of writing this document.*

#### Buy STABL

Buying STABL exchanges FIAT currency for STABL tokens.

Clicking on the BUY button brings up a dialog

![Register](assets/images/dialogBuy.png)

Input an amount to purchase and click on PROCEED

A request is sent to the REST server which creates a Hedera Consensus Transaction which you can witness in the [Kabuto](https://kabuto.sh) or [Dragonglass](https://testnet.dragonglass.me/hedera/home) explorers.

When the request is received by the Mirror Subscriber, it checks the signatures on the operation and credits the user’s account with the amount specified.

The result of the operation appears as a toaster at the bottom of the screen after a few seconds, and the page refreshes to show the new balance.

The balance now shows the amount you specified.

![Register](assets/images/operateAfterBuy.png)

#### REDEEM STABL 

Redeeming is the opposite of buying. STABL tokens are exchanged for FIAT currency.

Click on the REDEEM button, input an amount to redeem and click on PROCEED.

The result of the operation appears as a toaster at the bottom of the screen after a few seconds, and the page refreshes to show the new balance.

![Register](assets/images/operateAfterRedeem.png)

If an amount greater than the balance is specified (e.g. 120), a request is generated and sent to the application network. The Mirror Subscriber will check the balance of the account when it receives the request to redeem 120 and will fail the operation since there aren’t sufficient funds in the account.

Operations that negatively affect the balance of an account such as this redeem operation are signed by the end user using their private key. The Mirror Subscriber checks the signature on the operation corresponds to the public key it has on record for the user and rejects the operation if the signature is invalid.

Also, if the same operation was replayed on the Hedera Network by a malicious actor, the Mirror Subscriber would detect a duplicate signature (signatures using ED25519 keys can be considered unique) and reject the operation.

#### SEND STABL

Sending STABL transfers STABL tokens from one account to another. In order to perform this operation, another user must exist in the decentralized application database. Using another browser (to avoid cookie overlap), register a new user.

After successfully registering a new user (here we created a user called Hudson), click on the SEND button.

A dialog appears, where you can input an amount to send and type the name of a user to send to (the page will search through the list of names matching the search and return the first 10).

![Register](assets/images/dialogSend.png)

Click the PROCEED button, witness the message in the explorer
The result of the operation appears as a toaster at the bottom of the screen after a few seconds, and the page refreshes to show the new balance.

![Register](assets/images/accountsHudson.png)
