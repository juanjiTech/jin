# Jǐn

> 瑾，瑾瑜，美玉也。

Jin is a HTTP web framework written in [Go](https://go.dev/) (Golang) 
with a slim core but limitless extensibility.

## Feature

- Middleware support
- Dependency injection
- Integrate non-intrusively

## Performance

Due to the speed of `reflect.Call`, every inject process will take about
200ns, which means if the handler in handler-chain didn't support fast-invoke
will take about 200ns for dependency inject (on mac m2).

## Status

Alpha. Expect API changes and bug fixes.

## License

[MIT](./LICENSE)