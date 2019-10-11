<?php

declare(strict_types=1);

// worker.php
ini_set('display_errors', 'stderr');
include __DIR__ . '/../vendor/autoload.php';

$relay = new Spiral\Goridge\StreamRelay(STDIN, STDOUT);
$psr7 = new Spiral\RoadRunner\PSR7Client(new Spiral\RoadRunner\Worker($relay));

while ($req = $psr7->acceptRequest()) {
    try {
        $resp = new \Zend\Diactoros\Response();

        if ($req->getAttribute('broadcast:joinServer')) {
            $psr7->respond($resp->withStatus(200));
        } else {
            $psr7->respond($resp->withStatus(401));
        }
    } catch (\Throwable $e) {
        $psr7->getWorker()->error((string)$e);
    }
}
