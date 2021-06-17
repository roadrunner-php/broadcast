<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Broadcast\Tests\Topic;

use Spiral\RoadRunner\Broadcast\Broker;
use Spiral\RoadRunner\Broadcast\Topic\SingleTopic;

class SingleTopicTestCase extends TopicTestCase
{
    /**
     * @var string
     */
    protected string $name;

    public function setUp(): void
    {
        $this->name = \bin2hex(\random_bytes(32));

        parent::setUp();
    }

    public function testName(): void
    {
        $topic = $this->topic();
        $this->assertSame($this->name, $topic->getName());
    }

    /**
     * @param array<string, mixed> $mapping
     */
    protected function topic(array $mapping = []): SingleTopic
    {
        $broker = new Broker($this->rpc($mapping), $this->broker);

        return new SingleTopic($broker, $this->name);
    }
}
