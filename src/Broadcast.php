<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast;

use Spiral\Goridge\RPC\Codec\JsonCodec;
use Spiral\Goridge\RPC\Codec\ProtobufCodec;
use Spiral\Goridge\RPC\RPCInterface;
use Spiral\RoadRunner\Broadcast\DTO\V1\Message;
use Spiral\RoadRunner\Broadcast\DTO\V1\Request;
use Spiral\RoadRunner\Broadcast\DTO\V1\Response;
use Spiral\RoadRunner\Broadcast\Exception\BroadcastException;
use Spiral\RoadRunner\Broadcast\Exception\InvalidArgumentException;
use Spiral\RoadRunner\Broadcast\Topic\Group;
use Spiral\RoadRunner\Broadcast\Topic\SingleTopic;

final class Broadcast implements BroadcastInterface
{
    /**
     * @var RPCInterface
     */
    private RPCInterface $rpc;

    /**
     * @param RPCInterface $rpc
     */
    public function __construct(RPCInterface $rpc)
    {
        $this->rpc = $rpc->withCodec(new ProtobufCodec());
    }

    /**
     * @deprecated Information about RoadRunner plugins is not available since RoadRunner version 2.2
     */
    public function isAvailable(): bool
    {
        throw new \RuntimeException(\sprintf('%s::isAvailable method is deprecated.', self::class));
    }

    /**
     * @param non-empty-array<string> $topics
     * @param string $message
     * @return Message
     */
    private function createMessage(string $message, array $topics): Message
    {
        return new Message([
            'topics'  => $topics,
            'payload' => $message,
        ]);
    }

    /**
     * @param non-empty-list<Message> $messages
     * @return void
     * @throws BroadcastException
     */
    private function request(iterable $messages): void
    {
        $request = new Request(['messages' => $this->toArray($messages)]);

        /** @var Response $response */
        $response = $this->rpc->call('broadcast.Publish', $request, Response::class);

        if (! $response->getOk()) {
            throw new BroadcastException('An error occurred while publishing message');
        }
    }

    /**
     * @template T of mixed
     * @param iterable<T>|T $entries
     * @return array<T>
     */
    private function toArray($entries): array
    {
        switch (true) {
            case \is_array($entries):
                return $entries;

            case $entries instanceof \Traversable:
                return \iterator_to_array($entries, false);

            default:
                return [$entries];
        }
    }

    /**
     * {@inheritDoc}
     */
    public function publish($topics, $messages): void
    {
        assert(
            \is_string($topics) || \is_iterable($topics),
            '$topics argument must be a type of iterable<string>|string'
        );
        assert(
            \is_string($messages) || \is_iterable($messages),
            '$messages argument must be a type of iterable<string>|string'
        );

        /** @var array<string> $topics */
        $topics = $this->toArray($topics);

        if ($topics === []) {
            throw new InvalidArgumentException('Unable to publish message to 0 topics');
        }

        $request = [];

        /** @var string $message */
        foreach ($this->toArray($messages) as $message) {
            assert(\is_string($message), 'Message argument must be a type of string');

            $request[] = $this->createMessage($message, $topics);
        }

        if ($request === []) {
            throw new InvalidArgumentException('Unable to publish 0 messages');
        }

        $this->request($request);
    }

    /**
     * {@inheritDoc}
     */
    public function join($topics): TopicInterface
    {
        assert(
            \is_string($topics) || \is_iterable($topics),
            '$topics argument must be type of iterable<string>|string'
        );

        /** @var array<string> $topics */
        $topics = $this->toArray($topics);

        switch (\count($topics)) {
            case 0:
                throw new InvalidArgumentException('Unable to connect to 0 topics');

            case 1:
                return new SingleTopic($this, \reset($topics));

            default:
                return new Group($this, $topics);
        }
    }
}
