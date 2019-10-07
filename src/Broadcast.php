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
use Spiral\Goridge\Exceptions\ServiceException;
use Spiral\Goridge\RPC;

final class Broadcast implements BroadcastInterface
{
    // RPC service name
    private const SERVICE = 'broadcast';

    /** @var RPC */
    private $rpc;

    /**
     * @param RPC $rpc
     */
    public function __construct(RPC $rpc)
    {
        $this->rpc = $rpc;
    }

    /**
     * @inheritDoc
     */
    public function broadcast(Message ...$message): void
    {
        try {
            $this->rpc->call(
                sprintf("%s.Publish", self::SERVICE),
                $message
            );
        } catch (ServiceException $e) {
            throw new BroadcastException($e->getMessage(), $e->getCode(), $e);
        }
    }
}