// node bin/unfreeze.js --username joe

require('dotenv').config();
const { argv } = require("yargs");

async function main() {
    const stabl = require("../").create({
        adminKey: process.env.ADMIN_PRIVATE_KEY
    });

    const resp = await stabl.unfreeze(argv.username);

    console.log(resp);
}

main();
