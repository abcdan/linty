<?php

function doSomething()
{
    // Simulate some work
    return true;
}

function main()
{
    if (doSomething()) {
        echo "Success!\n";
    } else {
        echo "Error!\n";
    }
}

main();
?>
