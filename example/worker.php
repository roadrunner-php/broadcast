<?php

use Spiral\RoadRunner\Worker;
use Spiral\RoadRunner\Http\HttpWorker;

require __DIR__ . '/../vendor/autoload.php';

$http = new HttpWorker(Worker::create());

$script = <<<'html'
<!doctype html>
<html lang="en">
    <body>
    <code id="messages"></code>
    <script>
    window.addEventListener("load", function (evt) {
        let messages = document.getElementById("messages");

        let print = function (message) {
            messages.innerText += message + "\n";
        };

        let ws = new WebSocket('ws://127.0.0.1/ws');

        ws.onopen = function () {
            print("open");
            ws.send(`{"command":"join", "broker": "memory", "topics":["topic-1", "topic-2"]}`)
        };

        ws.onclose = function () {
            print("close");
        };

        ws.onmessage = function (evt) {
            print(`MESSAGE: ${evt.data}`);
        };

        ws.onerror = function (evt) {
            print("error: " + evt.data);
        };
    });
    </script>
    </body>
</html>
html;


while ($req = $http->waitRequest()) {
    $http->respond(200, $script);
}
