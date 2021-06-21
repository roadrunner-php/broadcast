<?php

use Spiral\Goridge\RPC\RPC;
use Spiral\RoadRunner\Broadcast\Broadcast;

require __DIR__ . '/../vendor/autoload.php';

$broadcast = new Broadcast(RPC::create('tcp://127.0.0.1:6001'));

// Send into broker
$broadcast->publish(['topic-1', 'topic-2'], 'message from broker');

// Send into topic
$topic = $broadcast->join(['topic-1', 'topic-2']);
$topic->publish('message from topic');
$topic->publish(['message 1 from topic', 'message 2 from topic']);
