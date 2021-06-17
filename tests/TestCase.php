<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast\Tests;

use PHPUnit\Framework\TestCase as BaseTestCase;
use Spiral\RoadRunner\Broadcast\DTO\V1\Response;
use Spiral\RoadRunner\Broadcast\Tests\Stub\RPCConnectionStub;

abstract class TestCase extends BaseTestCase
{
    /**
     * @param array<string, mixed> $mapping
     * @return RPCConnectionStub
     */
    protected function rpc(array $mapping = []): RPCConnectionStub
    {
        return new RPCConnectionStub($mapping);
    }

    /**
     * @param bool $success
     * @return string
     */
    protected function response(bool $success = true): string
    {
        return (new Response(['ok' => $success]))->serializeToString();
    }
}
