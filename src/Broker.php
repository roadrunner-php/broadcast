<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast;

use Spiral\Goridge\RPC\Codec\ProtobufCodec;
use Spiral\Goridge\RPC\RPCInterface;
use Spiral\RoadRunner\Broadcast\DTO\V1\Message;
use Spiral\RoadRunner\Broadcast\DTO\V1\Request;
use Spiral\RoadRunner\Broadcast\DTO\V1\Response;
use Spiral\RoadRunner\Broadcast\Topic\Group;
use Spiral\RoadRunner\Broadcast\Topic\SingleTopic;

class Broker implements BrokerInterface
{
    /**
     * @var RPCInterface
     */
    private RPCInterface $rpc;

    /**
     * @var string
     */
    private string $broker;

    /**
     * @param RPCInterface $rpc
     * @param string $name
     */
    public function __construct(RPCInterface $rpc, string $name)
    {
        $this->rpc = $rpc->withCodec(new ProtobufCodec());
        $this->broker = $name;
    }

    /**
     * {@inheritDoc}
     */
    public function getName(): string
    {
        return $this->broker;
    }

    /**
     * @param non-empty-array<string> $topics
     * @param string $message
     * @return Message
     */
    private function createMessage(string $message, array $topics): Message
    {
        return new Message([
            'broker'  => $this->broker,
            'topics'  => $topics,
            'payload' => $message,
        ]);
    }

    /**
     * @param non-empty-array<Message> $messages
     * @return Response
     */
    private function request(array $messages): Response
    {
        $request = new Request(['messages' => $messages]);

        return $this->rpc->call('websockets.Publish', $request, Response::class);
    }

    /**
     * @param iterable<string>|string $entry
     * @return array
     */
    private function toArrayOfStrings($entry): array
    {
        switch (true) {
            case \is_string($entry):
                return [$entry];

            case \is_array($entry):
                return $entry;

            case $entry instanceof \Traversable:
                return \iterator_to_array($entry, false);

            default:
                throw new \InvalidArgumentException(\sprintf(
                    'Argument must be a string or iterable<string>, but %s passed',
                    \get_debug_type($entry)
                ));
        }
    }

    /**
     * {@inheritDoc}
     */
    public function publish($topics, $messages): void
    {
        $topics = $this->toArrayOfStrings($topics);

        $map = fn(string $message): Message => $this->createMessage($message, $topics);
        $this->request(\array_map($map, $this->toArrayOfStrings($messages)));
    }

    /**
     * {@inheritDoc}
     */
    public function join($topics): TopicInterface
    {
        $topics = $this->toArrayOfStrings($topics);

        switch (\count($topics)) {
            case 0:
                throw new \InvalidArgumentException('Unable to connect to 0 topics');

            case 1:
                return new SingleTopic($this, \reset($topics));

            default:
                return new Group($this, $topics);
        }
    }
}
