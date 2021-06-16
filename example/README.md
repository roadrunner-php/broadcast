
## Installation

```sh
$ cd ~/broadcast # path to broadcast directory
$ composer install
$ cd example
$ rr get -s beta
$ rr serve
```

## Usage

```sh
$ cd ~/broadcast/example
$ php test.php
```

## API

### Connection

```php
$factory = new Factory(RPC::create('tcp://127.0.0.1:6001'));
```

### Publish Into Broker

```php
$broker = (new Factory(RPC::create('tcp://127.0.0.1:6001')))
    ->select('memory'); // or "redis"
// See "websockets:pubsubs: [ BROKER_NAME ]" section in ".rr.yaml" file

$broker->publish('topic', 'message');

// Other variants
$broker->publish('topic', ['message 1', 'message 2']);
$broker->publish(['topic 1', 'topic 2'], 'message');
$broker->publish(['topic 1', 'topic 2'], ['message 1', 'message 2']);
```

### Publish Into Topic

```php
$broker = (new Factory(RPC::create('tcp://127.0.0.1:6001')))
    ->select('memory');

$topic = $broker->join('topic');
// OR:
// $topic = $broker->join(['topic-1', 'topic-2', ...]);

$topic->publish('message'); // One message
$topic->publish(['message 1', 'message 2']); // Two messages
````
