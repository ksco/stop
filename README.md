# stop

Simple and extendable server monitor, it's `top` for the web.


### How to run?

```bash
# Backend
docker run --name stop -p 5566:5566 --env-file=<env file path> -v ~/.ssh:/root/.ssh -d --rm <image>
```


### License

The MIT License (MIT) - see [LICENSE](LICENSE) for more details