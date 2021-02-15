## go-inject in a nutshell

go-inject is a Golang library to allow simple creation of structs with complex dependency trees.
For example, imagine something like this:
```go
type Configuration struct {
    // contains some configuration values used in multiple program parts
}

type LowLevelComponent struct {
    Config *Configuration
}

type MidLevelComponent struct {
    Config *Configuration
    LowLevel *LowLevelComponent
    // Dependencies to other low level components
}

type HighLevelComponent struct {
    Config *Configuration
    MidLevelComponent *MidLevelComponent
    // Dependencies to other components
}
```

When this kind of dependency tree grows, it typically becomes hard to maintain since the structs must 
be initialized in the correct order, possibly they have custom initialization logic in each
struct that must be called, ...

go-inject aims to simply this kind of initialization:

```go
func createHighLevel() *HighLevelConfigUser {
    injector := &Injector{
        InjectableValues:[]interface{}{configuration}
    }
    var highLevel = &HighLevelConfigUser{}
    injector.Initialize(highLevel)
    return highLevel
}
```

In the background, go-inject walks through the structs, creates them, sets 
the injectable values, where applicable, and calls initializers.

