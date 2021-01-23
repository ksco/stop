# stop

`stop` is a **s**imple remote server monitor, just like **top**, but based on the web, with a neat UI.


### How to run?

```bash
# Backend
docker run --name stop -p 5566:5566 --env-file=<env file path> -v ~/.ssh:/root/.ssh -d --rm <image>
```
