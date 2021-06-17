# RoadRunner Broadcast Plugin Bridge

[![Latest Stable Version](https://poser.pugx.org/spiral/roadrunner-broadcast/version)](https://packagist.org/packages/spiral/roadrunner-broadcast)
[![Build Status](https://github.com/spiral/roadrunner-broadcast/workflows/build/badge.svg)](https://github.com/spiral/roadrunner-broadcast/actions)
[![Codecov](https://codecov.io/gh/spiral/roadrunner-broadcast/branch/master/graph/badge.svg)](https://codecov.io/gh/spiral/roadrunner-broadcast/)

This repository contains the codebase bridge for broadcast RoadRunner plugin.

## Installation

To install application server and broadcast codebase

```bash
$ composer require spiral/roadrunner-broadcast
```

You can use the convenient installer to download the latest available compatible
version of RoadRunner assembly:

```bash
$ composer require spiral/roadrunner-cli --dev
$ vendor/bin/rr get
```

## Usage

First you need to add at least one pub/sub broker to your RoadRunner configuration
and create http worker. Each of the dependencies below is required.

For example, such a configuration would be quite feasible to run:

```yaml
rpc:
  listen: tcp://127.0.0.1:6001

server:
  # Don't forget to create a "worker.php" file
  command: "php worker.php" 
  relay: "pipes"

http:
  address: 127.0.0.1:80
  # Indicate that HTTP support ws protocol
  middleware: [ "websockets" ]

websockets:
  # Using an in-memory broker named "memory"
  pubsubs: [ "memory" ]
  path: "/broadcast"
```

> Read more about all available brokers on the
> [documentation](https://roadrunner.dev/docs) page.

After configuring and starting the RoadRunner server, the corresponding API
will become available to you.

```php
<?php

use Spiral\Goridge\RPC\RPC;
use Spiral\RoadRunner\Broadcast\Factory;

require __DIR__ . '/vendor/autoload.php';

$factory = new Factory(RPC::create('tcp://127.0.0.1:6001'));

if (!$factory->isAvailable()) {
    throw new \LogicException('The [broadcast] plugin not available');
}

//
// Select "memory" broker
//
$broker = $factory->select('memory');

//
// Now we can send a message to a specific topic
//
$broker->publish('channel-1', 'message for channel #1');
```

### Select Specific Topic

Alternatively, you can also use a specific topic (or set of topics) as a 
separate entity and post directly to it.

```php
// Let's say we have already selected a specific broker
$broker = $factory->select('memory');

// Now we can select the topic we need to work only with it
$topic = $broker->join(['channel-1', 'channel-2']);

// And send messages there
$topic->publish('message');
$topic->publish(['another message', 'third message']);
```

> Read more about all the possibilities in the
> [documentation](https://roadrunner.dev/docs) page.

## Client

In addition to the server (PHP) part, the client part is also present in most
projects. In most cases, this is a browser in which the connection to the server
is made using the [WebSocket](https://en.wikipedia.org/wiki/WebSocket) protocol.

```js
const ws = new WebSocket('ws://127.0.0.1/broadcast');

ws.onopen = e => {
    const message = {
        command: 'join',
        broker:  'memory',
        topics:  ['channel-1', 'channel-2']
    };

    ws.send(JSON.stringify(message));
};

ws.onmessage = e => {
    const message = JSON.parse(e.data);

    console.log(`${message.topic}: ${message.payload}`);
}
```


## Examples

Examples are available in the corresponding directory [./example](./example).

## License

The MIT License (MIT). Please see [`LICENSE`](./LICENSE) for more information. 
Maintained by [Spiral Scout](https://spiralscout.com).

