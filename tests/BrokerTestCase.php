<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast\Tests;

use Google\Protobuf\Internal\RepeatedField;
use Psr\Log\LoggerInterface;
use Spiral\RoadRunner\Broadcast\Broker;
use Spiral\RoadRunner\Broadcast\BrokerInterface;
use Spiral\RoadRunner\Broadcast\DTO\V1\Message;
use Spiral\RoadRunner\Broadcast\DTO\V1\Request;
use Spiral\RoadRunner\Broadcast\DTO\V1\Response;
use Spiral\RoadRunner\Broadcast\Exception\InvalidArgumentException;
use Spiral\RoadRunner\Broadcast\Topic\Group;
use Spiral\RoadRunner\Broadcast\Topic\SingleTopic;

class BrokerTestCase extends TestCase
{
    /** @psalm-suppress PropertyNotSetInConstructor */
    private string $name;

    public function setUp(): void
    {
        $this->name = \bin2hex(\random_bytes(32));
        parent::setUp();
    }

    /**
     * @param array<string, mixed> $mapping
     */
    private function broker(array $mapping = []): BrokerInterface
    {
        return new Broker($this->rpc($mapping), $this->name);
    }

    public function testName(): void
    {
        $broker = $this->broker();

        $this->assertSame($this->name, $broker->getName());
    }

    public function testPublishingSingleMessage(): void
    {
        $expected = \random_bytes(32);

        $broker = $this->broker([
            'websockets.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $message) {
                    $this->assertSame($expected, $message->getPayload());
                }

                return $this->response();
            },
        ]);

        $broker->publish('topic', $expected);
    }

    public function testPublishingAnArrayOfMessages(): void
    {
        $expected = [\random_bytes(32), \random_bytes(32), \random_bytes(32)];

        $broker = $this->broker([
            'websockets.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $i => $message) {
                    $this->assertArrayHasKey($i, $expected);
                    $this->assertSame($expected[$i], $message->getPayload());
                }

                return $this->response();
            },
        ]);

        $broker->publish('topic', $expected);
    }

    public function testPublishingIteratorOfMessages(): void
    {
        $expected = [\random_bytes(32), \random_bytes(32), \random_bytes(32)];

        $broker = $this->broker([
            'websockets.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $i => $message) {
                    $this->assertArrayHasKey($i, $expected);
                    $this->assertSame($expected[$i], $message->getPayload());
                }

                return $this->response();
            },
        ]);

        $broker->publish('topic', (fn() => yield from $expected)());
    }

    public function testPublishingToSingleTopic(): void
    {
        $expected = \bin2hex(\random_bytes(32));

        $broker = $this->broker([
            'websockets.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $message) {
                    $this->assertSame([$expected], [...$message->getTopics()]);
                }

                return $this->response();
            },
        ]);

        $broker->publish([$expected], '');
    }

    public function testPublishingToMultipleTopicsUsingArray(): void
    {
        $expected = [
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32))
        ];

        $broker = $this->broker([
            'websockets.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $message) {
                    $this->assertSame($expected, [...$message->getTopics()]);
                }

                return $this->response();
            },
        ]);

        $broker->publish($expected, 'test');
    }

    public function testPublishingToMultipleTopicsUsingIterator(): void
    {
        $expected = [
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32)),
            \bin2hex(\random_bytes(32))
        ];

        $broker = $this->broker([
            'websockets.Publish' => function (Request $request) use ($expected) {
                /** @var Message $message */
                foreach ($request->getMessages() as $message) {
                    $this->assertSame($expected, [...$message->getTopics()]);
                }

                return $this->response();
            },
        ]);

        $broker->publish((fn () => yield from $expected)(), 'test');
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

        $broker = $this->broker([
            'websockets.Publish' => function (Request $request) use ($expectedTopics, $expectedMessages) {
                /** @var Message $message */
                foreach ($request->getMessages() as $i => $message) {
                    $this->assertArrayHasKey($i, $expectedMessages);
                    $this->assertSame($expectedMessages[$i], $message->getPayload());
                    $this->assertSame($expectedTopics, [...$message->getTopics()]);
                }

                return $this->response();
            },
        ]);

        $broker->publish($expectedTopics, $expectedMessages);
    }

    public function testErrorOnBinaryTopicName(): void
    {
        $this->expectExceptionMessage('Expect utf-8 encoding');

        $broker = $this->broker(['websockets.Publish' => $this->response()]);
        $broker->publish(\random_bytes(32), 'msg');
    }

    public function testJoinToSingleTopic(): void
    {
        $topic = $this->broker()
            ->join(\bin2hex(\random_bytes(32)))
        ;

        $this->assertInstanceOf(SingleTopic::class, $topic);
    }

    public function testJoinToMultipleTopics(): void
    {
        $topic = $this->broker()
            ->join([\bin2hex(\random_bytes(32)), \bin2hex(\random_bytes(32))])
        ;

        $this->assertInstanceOf(Group::class, $topic);
    }

    public function testJoinErrorTopics(): void
    {
        $topic = $this->broker()
            ->join([\bin2hex(\random_bytes(32)), \bin2hex(\random_bytes(32))])
        ;

        $this->assertInstanceOf(Group::class, $topic);
    }

    public function testJoinToZeroTopics(): void
    {
        $this->expectException(InvalidArgumentException::class);

        $this->broker()
            ->join([])
        ;
    }

    public function testPublishIntoZeroTopics(): void
    {
        $this->expectException(InvalidArgumentException::class);

        $this->broker()
            ->publish([], 'message')
        ;
    }

    public function testPublishZeroMessages(): void
    {
        $this->expectException(InvalidArgumentException::class);

        $this->broker()
            ->publish('topic', [])
        ;
    }
}
