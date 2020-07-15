// node bin/update-admin-key.js 302e020...

require('dotenv').config();
const { argv } = require("yargs");

async function main() {
    const stabl = require("../").create({
        adminKey: process.env.ADMIN_PRIVATE_KEY
    });

    const resp = await stabl.updateAdminKey(argv._[0]);

    console.log(resp);
}

main();
