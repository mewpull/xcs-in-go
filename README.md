# XCS (eXtended Classifier System) in Go #

Development stage:
```
    Early thoughts
--> Prototype/incomplete code
    Initial working version
    Solidly tested
    Release candidate
    End-user ready
```

## Overview ##

This project provides a Go implementation of the XCS (eXtended
Classifier System) algorithm as described by Butz and Wilson (2000).

This code should not be used in production without additional testing.
Achieving 80% test coverage is a high priority task remaining.

## Acknowledgements ##

Permission to use the 'algorithmic description' within the above paper
as a base for the present work was sought and received from Springer
(personal communication, 2020-07-07). The author would therefore like
to thank both the original paper authors and Springer for making this
project possible.

The Boolean multiplexer implementation was guided by the description
provided by Wilson (1998). For a very clear explanation of Boolean
multiplexers, please see Figure 2 [on this
page](https://preview.tinyurl.com/ybut2zvn) or consult the original
figure by Urbanowicz and Browne (2017).

## Implementation Notes ##

The implementation herein is very close to the original specification.
There have been a few small modifications to better suit the language
used.

Caveat: The interfaces within the `mli` package should _not_ be
considered fully stable at the present time. Addition of multi-step
problems such as woods2 may require slight adjustments to these
interfaces.

## Parameters ##

At the present time, please consult Wilson (1998) and Butz and Wilson
(2000) for details of the parameters.

## Build Instructions ##

- Make sure golang is installed (version 1.11 was used here)
- Navigate to the project directory
- Navigate to the `./cmd` sub-directory
- Run `go clean` then `go build -o xcs-on-multiplexer .`
- Execute `./xcs-on-multiplexer`

## High Priority Tasks Remaining ##

- 80% test coverage
- Command line options
- External configuration file
- More idiomatic code
- Document the parameters
- Only export identifiers where necessary (appropriate 'visibility')
- Smaller methods

## Contributors ##

- Matthew R. Karlsen

## Documentation Note ##

All of the documents for my XCS repositories are based upon a common
template. For this reason, there may be subsantial overlap between
documents.

## References ##

Butz, M. V., & Wilson, S. W. (2000, September). An algorithmic
description of XCS. In International Workshop on Learning Classifier
Systems (pp. 253-272). Springer, Berlin, Heidelberg.

Urbanowicz, R. J., & Browne, W. N. (2017). Introduction to learning
classifier systems. Springer, Berlin, Heidelberg.

Wilson, S. W. (1998). Generalization in the XCS classifier system. In
J. R. Koza et al. (Eds.), Genetic Programming 1998: Proceedings of the
Third Annual Conference (pp. 665-674). Morgan Kaufmann.
