<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast\Topic;

use Spiral\RoadRunner\Broadcast\BroadcastInterface;

final class SingleTopic extends Topic
{
    /**
     * @param BroadcastInterface $broadcast
     * @param string $topic
     */
    public function __construct(BroadcastInterface $broadcast, string $topic)
    {
        parent::__construct($broadcast, [$topic]);
    }

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->topics[0];
    }
}
