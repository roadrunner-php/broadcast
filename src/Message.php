<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */
declare(strict_types=1);

namespace Spiral\Broadcast;

/**
 * Broadcast message.
 */
final class Message implements \JsonSerializable
{
    /** @var string */
    private $topic;

    /** @var mixed */
    private $payload;

    /**
     * @param string $topic
     * @param mixed  $payload
     */
    public function __construct(string $topic, $payload)
    {
        $this->topic = $topic;
        $this->payload = $payload;
    }

    /**
     * @return array
     */
    public function jsonSerialize(): array
    {
        return [
            'topic'   => $this->topic,
            'message' => $this->payload
        ];
    }
}