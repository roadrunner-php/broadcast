<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */
declare(strict_types=1);

namespace Spiral\Broadcast\Tests;

use PHPUnit\Framework\TestCase;
use Spiral\Broadcast\Broadcast;
use Spiral\Broadcast\Message;
use Spiral\Goridge\RPC;
use Spiral\Goridge\SocketRelay;
use Symfony\Component\Process\Process;

class BroadcastTest extends TestCase
{
    public function tearDown()
    {
        if (file_exists(__DIR__ . '/../log.txt')) {
            unlink(__DIR__ . '/../log.txt');
        }
    }

    public function testBroadcast()
    {
        $rpc = new RPC(new SocketRelay("localhost", 6001));
        $br = new Broadcast($rpc);

        $p = new Process("ws-client", dirname(__DIR__));
        $p->start();

        while (!file_exists(__DIR__ . '/../log.txt')) {
            usleep(1000);
            if ($p->getErrorOutput() !== "") {
                $this->fail($p->getErrorOutput());
            }
        }

        $br->broadcast(
            new Message("topic", "hello"),
            new Message("topic", ["key" => "value"])
        );

        while ($p->isRunning()) {
            usleep(1000);
        }

        $this->assertSame('{"topic":"@join","payload":["topic"]}
{"topic":"topic","payload":"hello"}
{"topic":"topic","payload":{"key":"value"}}
{"topic":"@leave","payload":["topic"]}
', file_get_contents(__DIR__ . '/../log.txt'));
    }
}
