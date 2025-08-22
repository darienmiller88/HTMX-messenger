const centrifuge = new Centrifuge('ws://localhost:8080/connection/websocket')

//Connects client to server using websockets powered through centrifuge API
centrifuge.on('connected', function(ctx){
    console.log("ctx:", ctx)
})

//Create a new 
const sub = centrifuge.newSubscription("chat")

//Receive messages from anyone subscribed to "chat" channel
sub.on('publication', function(ctx) {
    console.log("message:", ctx)
})

//subscribe to the channel (unnecessary second step in my opinion, but necessary)
sub.subscribe();

//Complete the server connection, allowing the client to send real time messages.
centrifuge.connect();