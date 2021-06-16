<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast;

/**
 * @psalm-import-type MessagesList from BrokerInterface
 */
interface TopicInterface
{
    /**
     * @param MessagesList $messages
     */
    public function publish($messages): void;
}
