const { CheckpointFactories, FileSystemWallet, Gateway } = require('fabric-network');

const connectionProfile = require('/home/quanth/Desktop/HLF/Test/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json');

const walletIdentity = 'admin';
const walletPath = './wallet';
const channelName = 'mychannel';

const listenerName = 'ChairsContract-listener';
const chaincodeName = 'erc20';
const contractName = 'ERC20Chaincode';
const eventName1 = 'transferEvent';

const wallet = new FileSystemWallet(walletPath);

async function main() {
  try {
    // connect to gateway and get network and contract
    console.log('connecting to gateway...');
    const gateway = new Gateway();

    await gateway.connect(connectionProfile, {
      discovery: {
        asLocalhost: true,
        enabled: true
      },
      identity: walletIdentity,
      wallet
    });
    console.log('connected to gateway...');
    console.log('retrieving network and contract...');
    const network = await gateway.getNetwork(channelName);
    const contract = network.getContract(chaincodeName, contractName);
    console.log('retrieved network and contract...');

    // Add contract listener
    /**
     * @param {String} listenerName the name of the event listener
     * @param {String} eventName the name of the event being listened to
     * @param {Function} callback the callback function with signature (error, event, blockNumber, transactionId, status)
     * @param {module:fabric-network.Network~EventListenerOptions} options
    **/
    console.log('adding contract listeners...');
    // Transfer events
    await contract.addContractListener(`${listenerName}-${eventName1}`, eventName1, async (err, event, blockNumber, transactionId, status) => {
      if (err) {
        console.error(err);
        return;
      }
      console.log(`Block Number: ${blockNumber} Transaction ID: ${transactionId} Status: ${status}`);
      console.log(event.payload.toString('utf8'));
      console.log("===================================================================")
    }, {
      replay: true // replay missed events on start up with the file system checkpointer
    });
    console.log('added contract listeners...');
    console.log('listening for events...');
  } catch (err) {
    console.error(err.message);
    process.exit(1);
  }
}

// add contract listener
main();