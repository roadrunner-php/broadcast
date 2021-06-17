<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast\Topic;

use Spiral\RoadRunner\Broadcast\BrokerInterface;
use Spiral\RoadRunner\Broadcast\Exception\InvalidArgumentException;
use Spiral\RoadRunner\Broadcast\TopicInterface;

abstract class Topic implements TopicInterface
{
    /**
     * @var BrokerInterface
     */
    private BrokerInterface $broker;

    /**
     * @psalm-suppress InvalidPropertyAssignmentValue
     * @var non-empty-list<string>
     */
    protected array $topics = [];

    /**
     * @param BrokerInterface $broker
     * @param iterable<string> $topics
     * @throws InvalidArgumentException
     */
    public function __construct(BrokerInterface $broker, iterable $topics)
    {
        $this->broker = $broker;

        foreach ($topics as $topic) {
            $this->topics[] = $topic;
        }

        /** @psalm-suppress TypeDoesNotContainType */
        if ($this->topics === []) {
            throw new InvalidArgumentException('Unable to create instance for 0 topic names');
        }
    }

    /**
     * {@inheritDoc}
     */
    public function publish($messages): void
    {
        $this->broker->publish($this->topics, $messages);
    }
}
