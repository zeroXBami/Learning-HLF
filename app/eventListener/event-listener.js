const { FileSystemWallet, Gateway } = require('fabric-network');

const connectionProfile = require('../../test-network/organizations/peerOrganizations/intage.example.com/connection-intage.json');

const walletIdentity = 'admin';
const walletPath = './wallet';
const channelName = 'mychannel';

const listenerName = 'ERC20-listener';
const erc20ChaincodeName = 'erc20';
const erc20ContractName = 'ERC20Chaincode';
const privDataChaincodeName = 'privateData';
const privDataContractName = 'DataMarket'
const transferEvent = 'transferEvent';
const uploadEvent = 'uploadEvent';
const shareEvent = 'shareEvent';


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
    const erc20Contract = network.getContract(erc20ChaincodeName, erc20ContractName);
    const privDatacontract = network.getContract(privDataChaincodeName, privDataContractName);
    console.log('retrieved network and contract...');

    // Add contract listener
    /**
     * @param {String} listenerName the name of the event listener
     * @param {String} eventName the name of the event being listened to
     * @param {Function} callback the callback function with signature (error, event, blockNumber, transactionId, status)
     * @param {module:fabric-network.Network~EventListenerOptions} options
    **/
    console.log('adding contract listeners...');
    // Upload Data Event
    await privDatacontract.addContractListener(`${listenerName}-${uploadEvent}`, uploadEvent, async (err, event, blockNumber, transactionId, status) => {
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
    // Share Data Event
    await privDatacontract.addContractListener(`${listenerName}-${shareEvent}`, shareEvent, async (err, event, blockNumber, transactionId, status) => {
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
    // Transfer Token Event
    await erc20Contract.addContractListener(`${listenerName}-${transferEvent}`, transferEvent, async (err, event, blockNumber, transactionId, status) => {
      if (err) {
        console.error(err);
        return;
      }
      console.log("ERC20 Event")
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