# PPHLF
PP learning HLF
This is a sample app on HLF to store property records 

| Repo         | Link         |
| ------------ | ------------ |
| Fabric Repo  | https://github.com/hyperledger/fabric-samples |
| Edx Repo     | https://github.com/hyperledger/education/tree/master/LFS171x |
| yeasy Repo   | https://github.com/yeasy/docker-compose-files/tree/master/hyperledger_fabric/v1.2.0 |
| AWS Template | https://docs.aws.amazon.com/blockchain-templates/latest/developerguide/blockchain-templates-hyperledger.html#blockchain-hyperledger-launch |

Learn how to start building applications with Hyperledger Fabric. 
This example makes use of configtxgen and is for development ONLY. Make use of Fabric-CA for production environments

1. Define/Change crypto-config.yaml (Make sure you set EnableNodeOUs: true) 
2. Network Topology - OrdererOrgs - Name, Domain, Hostname, PeerOrgs - Name, Domain Users - Count in addition to Admin
3. Define/Change docker-compose.yml (FABRIC_CA_SERVER_CA_KEYFILE)
4. Define/Change ./generate.sh (Based on Profiles defined in configtx.yaml)
5. Define/Change ./start.sh (Change channel name)
6. Run startFabric.sh, and replace the new key by inspecting docker logs ca.example.com (Everytime you generate new crypto material, make sure you update "FABRIC_CA_SERVER_CA_KEYFILE" in the "docker-compose.yml" file)
7. Bring up the n/w by running startFabric.sh (> launch network; create channel and join peer to channel > launch CLI container in order to install, instantiate chaincode > bootstrap/ledger with property)
8. Enroll Admin and Register the User (node registerAdmin.js, node registerUser.js)
9. Start server (node server.js http://localhost:8000)

To Do/Findings:

Investigate the certificates generated by cryptogen using `openssl x509 -in ./path/to/cert.pem -text -noout` 

`msp/config.yaml` should NOT be placed into both the orderer and peers' MSP directory. It must ONLY go into the peers.
