# monkey-go

This repo contains my code from reading the following:

- [Writing An Interpreter In Go](https://interpreterbook.com/) (finished 2024-07-23)
- [The Lost Chapter](https://thorstenball.com/blog/2017/06/28/the-lost-chapter-a-macro-system-for-monkey/) (finished 2024-07-27)
- [Writing A Compiler In Go](https://compilerbook.com/) (finished 2024-10-25)

## Benchmark Results

```zsh
❯ go build -o fib ./benchmark
❯ ./fib -engine=eval
engine=eval, result=9227465, duration=14.605212132s
❯ ./fib -engine=vm
engine=vm, result=9227465, duration=4.792953424s
```

thanks Thorsten, this was fun :)
