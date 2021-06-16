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

final class SingleTopic extends Topic
{
    /**
     * @param BrokerInterface $broker
     * @param string $topic
     */
    public function __construct(BrokerInterface $broker, string $topic)
    {
        parent::__construct($broker, [$topic]);
    }

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->topics[0];
    }
}
