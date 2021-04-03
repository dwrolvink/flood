# Start situation
Interval without skips: 31 ms

## full run
```
Avg time per Update:  5942 us
Max time per Update:  19884 us
Min time per Update:  1523 us

Avg time per Update:  5897 us
Max time per Update:  21005 us
Min time per Update:  1436 us

Avg time per Update:  5336 us
Max time per Update:  20405 us
Min time per Update:  1574 us

Avg time per Draw Frame:  2004 us
Max time per Draw Frame:  6500 us
Min time per Draw Frame:  658 us

Avg time per Draw Frame:  2055 us
Max time per Draw Frame:  6892 us
Min time per Draw Frame:  633 us

Avg time per Draw Frame:  1955 us
Max time per Draw Frame:  5788 us
Min time per Draw Frame:  624 us


Avg time per Battle:  573 us
Max time per Battle:  1408 us
Min time per Battle:  237 us

Avg time per Battle:  701 us
Max time per Battle:  2239 us
Min time per Battle:  265 us


```

## till collision
```
Avg time per Update:  3373 us
Max time per Update:  11240 us
Min time per Update:  1586 us

Avg time per Update:  4536 us
Max time per Update:  16496 us
Min time per Update:  1604 us

Avg time per Update:  4679 us
Max time per Update:  18480 us
Min time per Update:  1680 us

Avg time per Draw Frame:  1259 us
Max time per Draw Frame:  5759 us
Min time per Draw Frame:  653 us

Avg time per Draw Frame:  1566 us
Max time per Draw Frame:  4448 us
Min time per Draw Frame:  616 us

Avg time per Draw Frame:  1574 us
Max time per Draw Frame:  5800 us
Min time per Draw Frame:  644 us
```

# Changes
## Use array referencing instead of direct access + extra go routines
Interval without skips: 20 ms (30 --> 18)

```


Old:
Avg time per Update:  5336 us
Max time per Update:  20405 us
Min time per Update:  1574 us

Avg time per Draw Frame:  1955 us
Max time per Draw Frame:  5788 us
Min time per Draw Frame:  624 us

Avg time per Battle:  573 us
Max time per Battle:  1408 us
Min time per Battle:  237 us

New:
Avg time per Update:  4583 us
Max time per Update:  13158 us
Min time per Update:  1586 us

Avg time per Draw Frame:  1643 us
Max time per Draw Frame:  5723 us
Min time per Draw Frame:  577 us

Avg time per Battle:  266 us
Max time per Battle:  881 us
Min time per Battle:  116 us