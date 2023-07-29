# dev_agent

## Use

```bash
dev_agent execute plan.yaml data.yaml?
ask for all params in template

what params from file?
args design
execute plan.yaml -p "file=file(xxx.go) function=MyHero"
```

# args design

```
# example
file=file(xxx.go) function=MyHero
```

- split by whitespace
- split by `=`
- left store as map key
- right part as template action, and execute with Funcs, output store as map value



## Tools/Actions
template function provider - like xxx api or pure http get

too limited