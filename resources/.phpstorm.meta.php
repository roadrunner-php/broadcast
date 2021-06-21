<?php

namespace PHPSTORM_META {

    //
    // RPC Methods
    //

    registerArgumentsSet('goridge_rpc_methods_broadcast',
        'broadcast.Publish',
        'broadcast.PublishAsync',
    );

    expectedArguments(\Spiral\Goridge\RPC\RPCInterface::call(), 0,
        argumentsSet('goridge_rpc_methods_broadcast'));
    expectedArguments(\Spiral\Goridge\RPC\RPC::call(), 0,
        argumentsSet('goridge_rpc_methods_broadcast'));

}
