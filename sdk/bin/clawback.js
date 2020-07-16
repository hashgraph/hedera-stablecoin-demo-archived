// node bin/unfreeze.js --username greg

require('dotenv').config();
const { argv } = require("yargs");

async function main() {
    const stabl = require("../").create({
        adminKey: process.env.ADMIN_PRIVATE_KEY
    });

    const resp = await stabl.clawback(argv.username);

    console.log(resp);
}

main();
