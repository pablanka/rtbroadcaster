function RTBroadcaster(url, onConnectionCallback, openCallback, closeCallback, errorCallback) {
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

    function sendMessageToServer(command) {
        if (!conn) {
            return false;
        }
        if (!command) {
            return false;
        }
        conn.send(command);
    }

    function onConnectionOpen() {
        var statusValue = 2;
        var paramUUID = getUrlParam("uuid");
        if (!paramUUID) {
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
            SateMessage: false,
            StateMessageID: -1
        }
        sendMessageToServer(JSON.stringify(message));
        if (openCallback) { openCallback(); }
    }

    function onConnectionClose(evt) {
        console.log("Connection closed.");
        if (closeCallback) { closeCallback(); }
    }

    function onReceiveConnectionMessage(evt) {
        var messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            var message = messages[i];
            objMessage = JSON.parse(message);
            if (objMessage.FuncKey !== "CameraRot") {
                console.log("MESSAGE: ", objMessage);
            }
            if (objMessage.UUID) {
                ref.uuid = objMessage.UUID;
                switch (objMessage.Status.Value) {
                    case 0:
                        //
                        break;
                    case 3:
                        if (!ref.connected) {
                            ref.connected = true;
                            if (ref.isOwner) {
                                console.log("ROOM UUID: ", objMessage.UUID);
                            }
                            if (onConnectionCallback) { onConnectionCallback(objMessage.UUID); }
                        }
                        if (objMessage.FuncKey) {
                            var _func = ref.suscribedFuncs.get(objMessage.FuncKey);
                            _func(objMessage.FuncParams);
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
        if (errorCallback) { errorCallback(error); }
    }
}

// API
RTBroadcaster.prototype.suscribeFunc = function (key, _func) {
    this.suscribedFuncs.set(key, _func);
}

RTBroadcaster.prototype.unsuscribeFunc = function (key) {
    this.suscribedFuncs.delete(key);
}

RTBroadcaster.prototype.sendAction = function (key, params, itsStateMessage, stateID) {
    if (!itsStateMessage) {
        itsStateMessage = false;
    }
    var uuid = this.uuid;
    var message = {
        UUID: uuid,
        Status: {
            Value: 3, // Connetion status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = closed
            'Text': "connected"
        },
        FuncKey: key,
        FuncParams: params,
        SateMessage: itsStateMessage,
        StateMessageID: (stateID ? stateID : -1)
    }
    this.sendMessage(JSON.stringify(message));
}