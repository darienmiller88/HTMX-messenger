const centrifuge = new Centrifuge('ws://localhost:8000/connection/websocket')

centrifuge.on('connected', function(ctx){
    console.log("ctx:", ctx)
});

centrifuge.connect();