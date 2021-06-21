<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast\Tests\Topic;

use Spiral\RoadRunner\Broadcast\Broadcast;
use Spiral\RoadRunner\Broadcast\Exception\InvalidArgumentException;
use Spiral\RoadRunner\Broadcast\Topic\Group;

class TopicGroupTestCase extends TopicTestCase
{
    /**
     * @var non-empty-array<string>
     */
    protected array $names;

    public function setUp(): void
    {
        $this->names = [\bin2hex(\random_bytes(32)), \bin2hex(\random_bytes(32))];
        parent::setUp();
    }

    public function testNames(): void
    {
        $topic = $this->topic();
        $this->assertSame($this->names, $topic->getNames());
    }

    /**
     * @param array<string, mixed> $mapping
     * @throws InvalidArgumentException
     */
    protected function topic(array $mapping = []): Group
    {
        $broadcast = new Broadcast($this->rpc($mapping));

        return new Group($broadcast, $this->names);
    }

    public function testCreatingWithZeroTopicNames(): void
    {
        $this->expectException(InvalidArgumentException::class);
        new Group(new Broadcast($this->rpc()), []);
    }
}
