# Fixturec

`fixturec` is a command-line utility that automatically generates test fixtures and mocks for Go structs.

🚀 Installation
```
go install github.com/Vypolor/fixturec@latest
```

🧠 How It Works

Use following commands for generate fixture:

```
cd path/to/struct/package
fixturec -t Impl
```


where Impl is the name of the struct for which you want to generate a fixture.

The tool performs the following steps:
- analyzes the specified struct and finds fields that are interfaces defined within the same module;
- checks for a //go:generate mockgen ... directive in the interface’s file and adds it if missing;
- runs go generate to create mocks using mockgen;
- generates a fixture_test.go file in the same package as the original struct.

⚙️ Flags

| Flag | Description                                                                                          |
|------|------------------------------------------------------------------------------------------------------|
| `-t` | **(required)** — name of the struct to generate the fixture for. Example: `-t Impl`.                 |
| `-g` | *(planned)* — disables automatic mock generation. Enabled by default. Currently **not implemented**. |
| `-e` | *(planned)* — enables mock generation for **external packages**. Currently **not implemented**.      |
