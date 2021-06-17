<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast;

use Spiral\Goridge\RPC\Codec\JsonCodec;
use Spiral\Goridge\RPC\RPCInterface;

final class Factory implements FactoryInterface
{
    /**
     * @var RPCInterface
     */
    private RPCInterface $rpc;

    /**
     * @param RPCInterface $rpc
     */
    public function __construct(RPCInterface $rpc)
    {
        $this->rpc = $rpc;
    }

    /**
     * {@inheritDoc}
     */
    public function isAvailable(): bool
    {
        try {
            /** @var array<string>|mixed $result */
            $result = $this->rpc
                ->withCodec(new JsonCodec())
                ->call('informer.List', true);

            if (! \is_array($result)) {
                return false;
            }

            return \in_array('websockets', $result, true);
        } catch (\Throwable $e) {
            return false;
        }
    }

    /**
     * {@inheritDoc}
     */
    public function select(string $broker): Broker
    {
        return new Broker($this->rpc, $broker);
    }
}
