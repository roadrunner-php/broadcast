<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast;

interface FactoryInterface
{
    /**
     * Returns information about whether a broadcast plugin is available.
     *
     * @return bool
     */
    public function isAvailable(): bool;

    /**
     * @param string $broker
     * @return BrokerInterface
     */
    public function select(string $broker): BrokerInterface;
}
