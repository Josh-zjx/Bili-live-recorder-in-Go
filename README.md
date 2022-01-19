# Bili-live-recorder-in-Go
A simple Bilibili livestream recorder written in Go

This project is designed to be 

- A personal exercise of Golang
    - To gain experience in using Golang for real purpose
        - error handling, channel, concurrency and other topics
- A cli tool for recording livestream
    - Dockerization and cloud deployment
    - Might add GUI if cli is fully implemented and tested
- An experimental field for software development practice
    - Documentation and Tests
    - Design Patterns
    - Performance Tuning

This project is primal and naive now. Anyone who is currently looking for a versatile and robust recorder should check [BililiveRecorder](https://github.com/Bililive/BililiveRecorder)

## Basic Usage Example (working)
As source code
```
go run . 123123
```

As installed binary
```
bilibili-livestream-recorder 123123
```
## Working on 
- [ ] Add cli argument parser and more wrapper functions to fully use the API
- [ ] Add auto slice functionality
- [x] Refine naming strategy
