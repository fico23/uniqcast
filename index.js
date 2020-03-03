const NATS = require('nats')

const nc = NATS.connect({ url: 'nats://127.0.0.1:4222'})

if (process.argv.length < 3) {
  console.log('Please start with "node index.js {file path}"')
  process.exit(0)
}

nc.request('video', process.argv[2], { max: 1, timeout: 3000 }, (msg) => {
  if (msg instanceof NATS.NatsError && msg.code === NATS.REQ_TIMEOUT) {
    console.log('request timed out')
  } else if (msg.substr('/')) {
      console.log(`init segment path ${msg}`)
    } else {
      console.log(`error message: ${msg}`)
    }
    
    nc.close()
})
