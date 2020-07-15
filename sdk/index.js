const axios = require("axios");
const base64 = require("@stablelib/base64");
const proto = require("./proto/messages_pb");
const { Ed25519PrivateKey } = require("@hashgraph/sdk");
const nacl = require("tweetnacl");

class Client {
    constructor(options) {
        this._api = axios.create({
            baseURL: options.api || "http://localhost:3128",
        });

        this._adminKey = options.adminKey;
    }

    async _send(f) {
        const primitive = new proto.Primitive();
        const message = f(primitive);

        primitive.setHeader(makePrimitiveHeader(message.serializeBinary(), this._adminKey));

        const resp = await this._api.post("/v1/token/transaction", JSON.stringify({
            primitive: base64.encode(primitive.serializeBinary()),
        }), {
            headers: {
                "Content-Type": "application/json",
            }
        });

        return resp.data.message;
    }

    async freeze(username) {
        return await this._send((primitive) => {
            const message = new proto.Freeze();
            message.setAccount(username);

            primitive.setFreeze(message);

            return message;
        });
    }

    async unfreeze(username) {
        return await this._send((primitive) => {
            const message = new proto.UnFreeze();
            message.setAccount(username);

            primitive.setUnfreeze(message);

            return message;
        });
    }

    async clawback(username) {
        return await this._send((primitive) => {
            const message = new proto.Clawback();
            message.setAccount(username);

            primitive.setClawback(message);

            return message;
        });
    }

    async updateAdminKey(newAdminKey) {
        return await this._send((primitive) => {
            const message = new proto.UpdateAdminKey();
            message.setNewAdminKey(newAdminKey);

            primitive.setUpdateAdminKey(message);

            return message;
        });
    }
};

exports.create = function (options) {
    return new Client(options);
};

function makePrimitiveHeader(toSign, privateKey) {
    const rand = Math.floor(Math.random() * 10000000);
    const randomString = rand.toString();

    const randomBytes = new TextEncoder("utf-8").encode(randomString);
    const signMe = Int8Array.from([...toSign, ...randomBytes]);

    const edPrivateKey = Ed25519PrivateKey.fromString(privateKey);
    const publicKey = edPrivateKey.publicKey.toString();

    const signature = nacl.sign.detached(
        Uint8Array.from(signMe),
        edPrivateKey._keyData
    );

    const primitiveHeader = new proto.PrimitiveHeader();
    primitiveHeader.setRandom(rand);
    primitiveHeader.setSignature(signature);
    primitiveHeader.setPublickey(publicKey);

    return primitiveHeader;
}
