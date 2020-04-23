const error = (tag, ...message) => {
    console.log(`[${new Date().toLocaleString()}] - [ERROR] [${tag}]  - ${JSON.stringify(message)}`);
  };
  
  const info = (tag, ...message) => {
    console.log(`[${new Date().toLocaleString()}] - [INFO] [${tag}]  - ${JSON.stringify(message)}`);
  };
  
  
  export default {
    error,
    info,
  };