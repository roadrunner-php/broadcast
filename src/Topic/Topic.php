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
use Spiral\RoadRunner\Broadcast\TopicInterface;

abstract class Topic implements TopicInterface
{
    /**
     * @var BrokerInterface
     */
    private BrokerInterface $broker;

    /**
     * @var non-empty-array<string>
     */
    protected array $topics = [];

    /**
     * @param BrokerInterface $broker
     * @param iterable<string> $topics
     */
    public function __construct(BrokerInterface $broker, iterable $topics)
    {
        $this->broker = $broker;

        foreach ($topics as $topic) {
            $this->topics[] = $topic;
        }

        if ($this->topics === []) {
            throw new \InvalidArgumentException(\sprintf(
                'Unable to create Topic instance for 0 topics'
            ));
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


