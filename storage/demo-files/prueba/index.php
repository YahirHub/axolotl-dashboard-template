<?php

declare(strict_types=1);

$name = "Axolotl";
$items = ["PHP", "Go", "JavaScript"];

function greet(string $name): string
{
    return "Hola, {$name}";
}

echo greet($name);

foreach ($items as $item) {
    echo PHP_EOL . "- " . $item;
}
