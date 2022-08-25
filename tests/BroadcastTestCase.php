<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast\Tests;

use Spiral\RoadRunner\Broadcast\Broadcast;
use Spiral\RoadRunner\Broadcast\BroadcastInterface;
use Spiral\RoadRunner\Broadcast\DTO\V1\Message;
use Spiral\RoadRunner\Broadcast\DTO\V1\Request;
use Spiral\RoadRunner\Broadcast\Exception\InvalidArgumentException;
use Spiral\RoadRunner\Broadcast\Topic\Group;
use Spiral\RoadRunner\Broadcast\Topic\SingleTopic;

class BroadcastTestCase extends TestCase
{
    /**
     * @param array<string, mixed> $mapping
     */
    private function broadcast(array $mapping = []): BroadcastInterface
    {
        return new Broadcast($this->rpc($mapping));
    }

    public function testFactoryCreation(): void
    {
        $this->expectNotToPerformAssertions();
        $this->broadcast();
    }

    public function testIsAvailable(): void
    {
        $this->expectException(\RuntimeException::class);
        $this->expectErrorMessage('Spiral\RoadRunner\Broadcast\Broadcast::isAvailable method is deprecated.');

        $this->broadcast()->isAvailable();
    }

    public function testPublishingSingleMessage(): void
    {
        $expected = \random_bytes(32);

        $broadcast = $this->broadcast([
            'broadcast.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $message) {
                    $this->assertSame($expected, $message->getPayload());
                }

                return $this->response();
            },
        ]);

        $broadcast->publish('topic', $expected);
    }

    public function testPublishingAnArrayOfMessages(): void
    {
        $expected = [\random_bytes(32), \random_bytes(32), \random_bytes(32)];

        $broadcast = $this->broadcast([
            'broadcast.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $i => $message) {
                    $this->assertArrayHasKey($i, $expected);
                    $this->assertSame($expected[$i], $message->getPayload());
                }

                return $this->response();
            },
        ]);

        $broadcast->publish('topic', $expected);
    }

    public function testPublishingIteratorOfMessages(): void
    {
        $expected = [\random_bytes(32), \random_bytes(32), \random_bytes(32)];

        $broadcast = $this->broadcast([
            'broadcast.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $i => $message) {
                    $this->assertArrayHasKey($i, $expected);
                    $this->assertSame($expected[$i], $message->getPayload());
                }

                return $this->response();
            },
        ]);

        $broadcast->publish('topic', (fn() => yield from $expected)());
    }

    public function testPublishingToSingleTopic(): void
    {
        $expected = \bin2hex(\random_bytes(32));

        $broadcast = $this->broadcast([
            'broadcast.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $message) {
                    $this->assertSame([$expected], [...$message->getTopics()]);
                }

                return $this->response();
            },
        ]);

        $broadcast->publish([$expected], '');
    }

    public function testPublishingToMultipleTopicsUsingArray(): void
    {
        $expected = [
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32))
        ];

        $broadcast = $this->broadcast([
            'broadcast.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $message) {
                    $this->assertSame($expected, [...$message->getTopics()]);
                }

                return $this->response();
            },
        ]);

        $broadcast->publish($expected, 'test');
    }

    public function testPublishingToMultipleTopicsUsingIterator(): void
    {
        $expected = [
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32))
        ];

        $broadcast = $this->broadcast([
            'broadcast.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $message) {
                    $this->assertSame($expected, [...$message->getTopics()]);
                }

                return $this->response();
            },
        ]);

        $broadcast->publish((fn () => yield from $expected)(), 'test');
    }

    public function testPublishingMultipleMessagesToMultipleTopics(): void
    {
        $expectedTopics = [
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32))
        ];

        $expectedMessages = [
            \random_bytes(32),
            \random_bytes(32)
        ];

        $broadcast = $this->broadcast([
            'broadcast.Publish' => function (Request $request) use ($expectedTopics, $expectedMessages) {
                /** @var Message $message */
                foreach ($request->getMessages() as $i => $message) {
                    $this->assertArrayHasKey($i, $expectedMessages);
                    $this->assertSame($expectedMessages[$i], $message->getPayload());
                    $this->assertSame($expectedTopics, [...$message->getTopics()]);
                }

                return $this->response();
            },
        ]);

        $broadcast->publish($expectedTopics, $expectedMessages);
    }

    public function testErrorOnBinaryTopicName(): void
    {
        $this->expectExceptionMessage('Expect utf-8 encoding');

        $broadcast = $this->broadcast(['broadcast.Publish' => $this->response()]);
        $broadcast->publish(\random_bytes(32), 'msg');
    }

    public function testJoinToSingleTopic(): void
    {
        $topic = $this->broadcast()
            ->join(\bin2hex(\random_bytes(32)))
        ;

        $this->assertInstanceOf(SingleTopic::class, $topic);
    }

    public function testJoinToMultipleTopics(): void
    {
        $topic = $this->broadcast()
            ->join([\bin2hex(\random_bytes(32)), \bin2hex(\random_bytes(32))])
        ;

        $this->assertInstanceOf(Group::class, $topic);
    }

    public function testJoinErrorTopics(): void
    {
        $topic = $this->broadcast()
            ->join([\bin2hex(\random_bytes(32)), \bin2hex(\random_bytes(32))])
        ;

        $this->assertInstanceOf(Group::class, $topic);
    }

    public function testJoinToZeroTopics(): void
    {
        $this->expectException(InvalidArgumentException::class);

        $this->broadcast()
            ->join([])
        ;
    }

    public function testPublishIntoZeroTopics(): void
    {
        $this->expectException(InvalidArgumentException::class);

        $this->broadcast()
            ->publish([], 'message')
        ;
    }

    public function testPublishZeroMessages(): void
    {
        $this->expectException(InvalidArgumentException::class);

        $this->broadcast()
            ->publish('topic', [])
        ;
    }
}
