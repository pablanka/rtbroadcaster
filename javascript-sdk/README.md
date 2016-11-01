# rtbroadcaster Javascript SDK

## How to use it

Create an instas of RTBroadcaster

```javascript
var broadcaster = new RTBroadcaster(url, onConnectionCallback, openCallback, closeCallback, errorCallback);
```

## API:

### RTBroadcaster.suscribeFunc(key, func)



```javascript
broadcaster.suscribeFunc(key, _func);
```

### RTBroadcaster.unsuscribeFunc(key)

```javascript
broadcaster.unsuscribeFunc(key);
```

### RTBroadcaster.sendAction(key, params, itsStateMessage)

```javascript
broadcaster.sendAction(key, params, itsStateMessage);
```

### RTBroadcaster.connected

```javascript
broadcaster.connected;
```

### RTBroadcaster.isOwner

```javascript
broadcaster.isOwner;
```

## Authors

* **Pablo Acuña**


