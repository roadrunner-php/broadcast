<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast\Tests\Topic;

use Spiral\RoadRunner\Broadcast\DTO\V1\Message;
use Spiral\RoadRunner\Broadcast\DTO\V1\Request;
use Spiral\RoadRunner\Broadcast\Exception\InvalidArgumentException;
use Spiral\RoadRunner\Broadcast\Tests\TestCase;
use Spiral\RoadRunner\Broadcast\TopicInterface;

/**
 * @psalm-suppress PropertyNotSetInConstructor
 */
abstract class TopicTestCase extends TestCase
{
    public function setUp(): void
    {
        parent::setUp();
    }

    /**
     * @param array<string, mixed> $mapping
     */
    abstract protected function topic(array $mapping = []): TopicInterface;

    public function testPublishingSingleMessage(): void
    {
        $expected = \random_bytes(32);
        $topic = $this->topic([
            'broadcast.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $message) {
                    $this->assertSame($expected, $message->getPayload());
                }

                return $this->response();
            },
        ]);

        $topic->publish($expected);
    }

    public function testPublishingAnArrayOfMessages(): void
    {
        $expected = [\random_bytes(32), \random_bytes(32), \random_bytes(32)];

        $topic = $this->topic([
            'broadcast.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $i => $message) {
                    $this->assertArrayHasKey($i, $expected);
                    $this->assertSame($expected[$i], $message->getPayload());
                }

                return $this->response();
            },
        ]);

        $topic->publish($expected);
    }

    public function testPublishingIteratorOfMessages(): void
    {
        $expected = [\random_bytes(32), \random_bytes(32), \random_bytes(32)];

        $topic = $this->topic([
            'broadcast.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $i => $message) {
                    $this->assertArrayHasKey($i, $expected);
                    $this->assertSame($expected[$i], $message->getPayload());
                }

                return $this->response();
            },
        ]);

        $topic->publish((fn() => yield from $expected)());
    }

    public function testPublishZeroMessages(): void
    {
        $this->expectException(InvalidArgumentException::class);

        $this->topic()
            ->publish([])
        ;
    }
}
