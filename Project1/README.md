# CS523-Project1
This is the repository for the project 1 of the CS-523 course @ EPFL, spring 2020

### Team members
- Romain Mendez
- Julien Heitmann

This project aims at designing a N-party multiparty computations (MPC) engine in a semi-honest (passive) adversarial setting in the Go programming language. The MPC framework works with circuits that contain basic arithmetic operations over additively secret-shared data. Multiplications are also supported, but require a pre-processing phase using BFV homomorphic encryption with Lattigo library to produce Beaver triplets. 

### Executing the code
The code can be tested as follows:

```
go test -v -run ^TestEval$
```

This will check that all test-circuits defined in `test_circuits.go` execute correctly, including a circuit we added that contains all the supported gates.

Finally, `main.go` will execute the dating-circuit we describe in the report, and expects binary inputs:

```
./mpc 0 1 & ./mpc 1 1 & ./mpc 2 1 & ./mpc 3 1
```

### Executing the benchmarks
If you wish to test the benchmarks for the different circuits that exists, the circuits and code for the benchmark have all been compiled to `benchmark_test.go`.
To execute all the benchmarks, you can use this command :
```
go test -run=XXX -bench=. -benchtime=T -timeout 99999s
```

Where `T` is the amount of time you want to use for each benchmarks, for example `10s` will execute the code in a loop for 10 seconds and then show you the results.
Please note that the print statements inside the code may render the results for multiple benchmarks unreadable, if you wish to run them properly, either remove all the prints inside the code being benchmarked, or you can select a particular benchmark :
```
go test -run=XXX -bench=MulCircuit1 -benchtime=600s -timeout 99999s
```

This command will for example only run the circuit with 1 multiplication for 10 min.