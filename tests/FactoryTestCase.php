<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast\Tests;

use Spiral\RoadRunner\Broadcast\Factory;
use Spiral\RoadRunner\Broadcast\FactoryInterface;

class FactoryTestCase extends TestCase
{
    /**
     * @param array<string, mixed> $mapping
     */
    private function factory(array $mapping = []): FactoryInterface
    {
        return new Factory($this->rpc($mapping));
    }

    public function testFactoryCreation(): void
    {
        $this->expectNotToPerformAssertions();
        $this->factory();
    }

    public function testIsAvailable(): void
    {
        $factory = $this->factory(['informer.List' => '["websockets"]']);
        $this->assertTrue($factory->isAvailable());
    }

    public function testNotAvailable(): void
    {
        $factory = $this->factory(['informer.List' => '[]']);
        $this->assertFalse($factory->isAvailable());
    }

    public function testNotAvailableOnNonArrayResponse(): void
    {
        $factory = $this->factory(['informer.List' => '42']);
        $this->assertFalse($factory->isAvailable());
    }

    public function testNotAvailableOnErrorResponse(): void
    {
        $factory = $this->factory(['informer.List' => (static function () {
            throw new \Exception();
        })]);

        $this->assertFalse($factory->isAvailable());
    }

    public function testBrokerSelection(): void
    {
        $name = \random_bytes(32);

        $broker = $this->factory()
            ->select($name)
        ;

        $this->assertSame($name, $broker->getName());
    }
}
