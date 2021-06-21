<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast;

use Spiral\RoadRunner\Broadcast\Exception\BroadcastException;

/**
 * @psalm-import-type MessagesList from BroadcastInterface
 */
interface TopicInterface
{
    /**
     * Method to send messages to the concrete selected topic.
     *
     * <code>
     *  $topic->publish('message');
     *  $topic->publish(['message 1', 'message 2']);
     * </code>
     *
     * Note: In future major releases, the signature of this method will be
     * changed to include follow type-hints.
     *
     * <code>
     *  public function publish(iterable|string $messages): void;
     * </code>
     *
     * @param MessagesList $messages
     * @throws BroadcastException
     */
    public function publish($messages): void;
}
