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
 * @psalm-type TopicsList = non-empty-list<string> | string
 * @psalm-type MessagesList = non-empty-list<string> | string
 */
interface BrokerInterface
{
    /**
     * Returns name of the target broker.
     *
     * @return string
     */
    public function getName(): string;

    /**
     * Method to send messages to the required topic (channel).
     * <code>
     *  $broker->send('topic', 'message');
     *  $broker->send('topic', ['message 1', 'message 2']);
     *
     *  $broker->send(['topic 1', 'topic 2'], 'message');
     *  $broker->send(['topic 1', 'topic 2'], ['message 1', 'message 2']);
     * </code>
     *
     * In future major releases, the signature of this method will match.
     * <code>
     *  public function publish(iterable|string $message, iterable|string $channels): void;
     * </code>
     *
     * @param TopicsList $topics
     * @param MessagesList $messages
     * @throws \InvalidArgumentException
     */
    public function publish($topics, $messages): void;

    /**
     * @param TopicsList $topics
     * @return TopicInterface
     * @throws \InvalidArgumentException
     */
    public function join($topics): TopicInterface;
}
