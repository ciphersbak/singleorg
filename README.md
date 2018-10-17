# PPHLF
PP learning HLF
This is a sample app on HLF to store property records
Original Source Code - https://github.com/hyperledger/fabric-samples && https://github.com/hyperledger/education/tree/master/LFS171x

Learn how to start building applications with Hyperledger Fabric. This example makes use of configtxgen and is for development ONLY. Make use of Fabric-CA for production environments
1. Define/Change crypto-config.yaml (Make sure you set EnableNodeOUs: true) (Network Topology - OrdererOrgs - Name, Domain, Hostname
PeerOrgs - Name, Domain)
Users - Count in addition to Admin)
2. Define/Change docker-compose.yml (FABRIC_CA_SERVER_CA_KEYFILE)
3. Define/Change ./generate.sh (Based on Profiles defined in configtx.yaml)
4. Define/Change ./start.sh (Change channel name)
5. Run startFabric.sh, and replace the new key by inspecting docker logs ca.example.com (Everytime you generate new crypto material, make sure you update "FABRIC_CA_SERVER_CA_KEYFILE" in the "docker-compose.yml" file)
6. Bring up the n/w by running startFabric.sh (> launch network; create channel and join peer to channel > launch CLI container in order to install, instantiate chaincode > bootstrap/ledger with property)
7. Enroll Admin and Register the User (node registerAdmin.js, node registerUser.js)
8. Start server (node server.js http://localhost:8000)
