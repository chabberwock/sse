# Simple SSE Server

This application helps to route and send Server-Side events to clients.

## Installation

* Install Golang (https://golang.org/doc/install#install)
* Build app: `cd sse/cmd/server; go build`
* Configure app autoload via supervisor or similar service

## Configuration

### Command-line arguments example:
```
./server -bind 127.0.0.1:8000 -bindctl 127.0.0.1:8001 -secret MySuperPassword42
```
* **-bind**: IP and port to accept client connections
* **-bindctl**: IP and port to send control requests. Normally this interface should not be availabke from outside
* **-secret**: used to sign and verify JWT session tokens

### Setting up Apache proxy

Suppose server runs on 127.0.0.1:8000, then just add ProxyPass directive to
your VirtualHost section

```
<VirtualHost *:443>
    ServerName yoursite.com
    ProxyPass /api/sse/ http://127.0.0.1:8000/events/
    ProxyPassReverse /api/sse/ http://127.0.0.1:8000/events/
    ProxyPreserveHost On
</VirtualHost>
```


### Subscribing to server-side events from client

Simple SSE Server uses JWT tokens to authorize client sessions, so your EventSource request should
pass Authorization header. This can be achieved by using https://github.com/Yaffle/EventSource/ or similar library

```
<script src="/js/eventsource.min.js"></script>
<script>
        var sse = new EventSourcePolyfill("/api/sse/", {
            headers: {
                'Authorization': "Bearer YOUR_JWT_TOKEN"
            }
        });
        
        sse.addEventListener('message', e => {
            var data = JSON.parse(e.data);
            console.log(data);
        });
</script>
```

### Generating JWT token

JWT is used to store userId, channel and to authorize access to event stream
You can manually generate token with following params:
`{"alg": "HS256", "typ": "JWT"}.{"userId":"your_user_id","channel":"channel name","exp":"expiretimestamp"}`

* **userId** is just a string used to identify recipient. user `*` for all recipients
* **channel** used to filter messages, client will only receive events sent to specified channel
* **exp** Unix timestamp after which token becomes invalid

You may also generate token by sending POST request to control endpoint. Example below creates token
for user _alex_, subscribed to channel _chat_ and sets token to expire in 24 hours
:
````
POST http://127.0.0.1:8001/token/
userId=alex&channel=chat&exp=86400
````

### Sending events

````
POST http://127.0.0.1:8001/emit/
userId=alex&channel=chat&payload=your_data
````

Note, payload sends message as is and does not care about validation, and in order to allow browser EventSource to properly parse it, it has to be sent in compatible format.

Example:
````
event: message
data: My important data for client
````

More about message format can be read here: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format



