# Unused Variables

```
About Unused Statement checking

We will need to analyze the program as a whole and somehow have a ref counter (bool: Used) that will allow us at
the end to make an unused statement (whether that is a variable function, return, etc) optimizer

- Assignment does not count effect the usage of the variable
```

```
I think in some regards, it may be easier to output an AST from the type checker that will use the same ident in all spots related to the referenced variable. This would allow us to save this ident along with a `used` attribute inside the VariableNode mapping and, throughout the processing of that tree, nullify the statement based on its usage. This would then make it entirely trivial to check for null idents when transpiling.
```