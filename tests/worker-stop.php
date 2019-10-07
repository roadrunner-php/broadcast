<?php
// worker.php
ini_set('display_errors', 'stderr');
include __DIR__ . "/../vendor/autoload.php";

$relay = new Spiral\Goridge\StreamRelay(STDIN, STDOUT);
$psr7 = new Spiral\RoadRunner\PSR7Client(new Spiral\RoadRunner\Worker($relay));

while ($req = $psr7->acceptRequest()) {
    try {
        $resp = new \Zend\Diactoros\Response();
        $psr7->respond($resp->withStatus(401));
    } catch (\Throwable $e) {
        $psr7->getWorker()->error((string)$e);
    }
}