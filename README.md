# Hedera Stable Coin Demo

> This is a sample implementation of an ERC20-like token in a 
> Hedera Consensus Service (HCS) Decentralized Application.

## Dependencies

 * PostgreSQL 12+
 
 * Go 1.14+ <sub>†</sub>
    
 * [migrate](https://github.com/golang-migrate/migrate) <sub>†</sub>
 
 * Protobuf Compiler (protoc) <sub>†</sub>
 
 * [`protoc-gen-go`](https://github.com/golang/protobuf) <sub>†</sub> 

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
