# Code generation

ССKit use code generation with [Buf](https://docs.buf.build/)

## Using generators

###  1. Install Buf 

https://docs.buf.build/installation

### 2. Install generators

With go package versions from CCKit [go.mod](../go.mod):

```
cd generators
./install.sh
```

### 3. Run buf

Generation command defined in [Makefile](../Makefile)
```
make proto
```
