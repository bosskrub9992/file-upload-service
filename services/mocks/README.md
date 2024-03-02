#### how to generate mock

1. install mockery
```
brew install mockery
```

2. run command
```
mockery --all --dir=services --output=services/mocks --outpkg=mocks --with-expecter=true --case=snake
```
