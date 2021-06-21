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
use Spiral\RoadRunner\Broadcast\Exception\InvalidArgumentException;

/**
 * @psalm-type TopicsList = non-empty-list<string> | string
 * @psalm-type MessagesList = non-empty-list<string> | string
 */
interface BroadcastInterface
{
    /**
     * Returns information about whether a broadcast plugin is available.
     *
     * @return bool
     */
    public function isAvailable(): bool;

    /**
     * Method to send messages to the required topic (channel).
     * <code>
     *  $broadcast->publish('topic', 'message');
     *  $broadcast->publish('topic', ['message 1', 'message 2']);
     *
     *  $broadcast->publish(['topic 1', 'topic 2'], 'message');
     *  $broadcast->publish(['topic 1', 'topic 2'], ['message 1', 'message 2']);
     * </code>
     *
     * Note: In future major releases, the signature of this method will be
     * changed to include follow type-hints.
     *
     * <code>
     *  public function publish(iterable|string $topics, iterable|string $messages): void;
     * </code>
     *
     * @param TopicsList $topics
     * @param MessagesList $messages
     * @throws BroadcastException
     */
    public function publish($topics, $messages): void;

    /**
     * Join to concrete topic
     *
     * @param TopicsList $topics
     * @return TopicInterface
     * @throws InvalidArgumentException
     */
    public function join($topics): TopicInterface;
}
