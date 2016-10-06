// TODO: Add public method to send messages to socket connection.

function RTBroadcaster(url, onConnectionCallback, openCallback, closeCallback, errorCallback){
    var conn;
    var suscribedFuncs = new Map();
    var ref = this;

    // API
    this.connected = false;
    this.isOwner = false;

    function getUrlParam(param) {
        var paramValue = location.search.split(param + '=')[1];
        if (!paramValue) {
            return "";
        }
        return paramValue;
    }

    function sendMessageToServer(command){
        if (!conn) {
            return false;
        }
        if (!command) {
            return false;
        }
        conn.send(command);
    }

    function onConnectionOpen(){
        var paramUUID = getUrlParam("uuid");
        if(!paramUUID){
            ref.isOwner = true;
        }
        var message = {
            uuid: paramUUID,
            status: {
                value: 0, // 0 = not connected, 1 = connected, 2 = closed
                text: "not connected"
            },
            funcKey: "",
            funcParams: []
        }
        sendMessageToServer(JSON.stringify(message));
        if(openCallback) {openCallback();}
    }

    function onConnectionClose(evt){
        console.log("Connection closed.");
        if(closeCallback){ closeCallback();}
    }

    function onReceiveConnectionMessage(evt){
        var messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            var message = messages[i];
            console.log("MESSAGE: ", message);
            if(message.uuid && message.status == 1){
                switch(message.status){
                    case 0:
                        //
                        break;
                    case 1:
                        if(!ref.connected){
                            ref.connected = true;
                            if(ref.isOwner){
                                console.log("ROOM UUID: ", message.uuid);
                                if(onConnectionCallback) {onConnectionCallback(message.uuid);}
                            }
                        }
                        if(message.funcKey){
                            suscribedFuncs.get(key)(message.funcParams);
                        }
                        break;
                    case 2:
                        // TODO: Close broadcasting, close connection, call onclose callback.
                        break;
                }
            }
        }
    }

    if (window["WebSocket"]) {
        conn = new WebSocket(url); // "ws://localhost:8080/broadcasting"
        conn.onopen = onConnectionOpen;
        conn.onclose = onConnectionClose;
        conn.onmessage = onReceiveConnectionMessage;
    } else {
        var error = "Your browser does not support WebSockets.";
        console.log(error);
        if(errorCallback) {errorCallback(error);}
    }
}

// API
RTBroadcaster.prototype.suscribeFunc = function(key, _func){
    suscribedFuncs.set(key, _func);
}

RTBroadcaster.prototype.unsuscribeFunc = function(key){
    suscribedFuncs.delete(key);
}