<?php

/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

declare(strict_types=1);

namespace Spiral\Broadcast;

use Spiral\Broadcast\Exception\BroadcastException;

/**
 * Provides the ability to broadcast messages to users.
 */
interface BroadcastInterface
{
    /**
     * Publish one or multiple messages.
     *
     * @param Message ...$message
     *
     * @throws BroadcastException
     */
    public function publish(Message ...$message): void;
}
