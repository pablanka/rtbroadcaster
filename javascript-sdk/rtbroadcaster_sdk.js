function RTBroadcaster(url, onConnectionCallback, openCallback, closeCallback, errorCallback){
    var conn;
    var ref = this;

    this.suscribedFuncs = new Map();
    this.connected = false;
    this.isOwner = false;
    this.sendMessage = sendMessageToServer;

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
        var statusValue = 2;
        var paramUUID = getUrlParam("uuid");
        if(!paramUUID){
            ref.isOwner = true;
            statusValue = 1;
        }
        var message = {
            uuid: paramUUID,
            status: {
                value: statusValue, // Connetion status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = closed
                text: "new connection"
            },
            funcKey: "",
            funcParams: [],
            stateMessage: false
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
            objMessage = JSON.parse(message);
            if(objMessage.funcKey !== "CameraRot"){
                console.log("MESSAGE: ", objMessage);
            }
            if(objMessage.uuid){
                ref.uuid = objMessage.uuid;
                switch(objMessage.status.value){
                    case 0:
                        //
                        break;
                    case 3:
                        if(!ref.connected){
                            ref.connected = true;
                            if(ref.isOwner){
                                console.log("ROOM UUID: ", objMessage.uuid);
                            }
                            if(onConnectionCallback) {onConnectionCallback(objMessage.uuid);}
                        }
                        if(objMessage.funcKey){
                            var _func = ref.suscribedFuncs.get(objMessage.funcKey);
                            _func(objMessage.funcParams);
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
        console.error(error);
        if(errorCallback) {errorCallback(error);}
    }
}

// API
RTBroadcaster.prototype.suscribeFunc = function(key, func){
    this.suscribedFuncs.set(key, func);
}

RTBroadcaster.prototype.unsuscribeFunc = function(key){
    this.suscribedFuncs.delete(key);
}

RTBroadcaster.prototype.sendAction = function(key, params, itsStateMessage){
    if(!itsStateMessage){
        itsStateMessage = false;
    }
    var uuid = this.uuid;
     var message = {
        uuid: uuid,
        status: {
            value: 3, // Connetion status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = closed
            'text': "connected"
        },
        funcKey: key,
        funcParams: params,
        sateMessage: itsStateMessage
    }
    this.sendMessage(JSON.stringify(message));
}