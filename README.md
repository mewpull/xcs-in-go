# XCS (eXtended Classifier System) in Go #

Development stage:
```
    Early thoughts
--> Prototype/working code
    Initial stable version
    Solidly tested
    Release candidate
    Production ready
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
page](https://ryanurbanowicz.com/index.php/resources-2/multiplexer-problem/)
or consult the original figure by Urbanowicz and Browne (2017).

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

- Make sure Go is installed (version 1.11 was used here)
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

## Contributing ##

Contributions are welcome, providing they are small and focused (and
thus can be code reviewed relatively easily).

Please remember to update the contributors section in the README.md
with your name as part of your first commit to the project.

Any contributions submitted will be dual-licensed, under the supplied
MIT and MIT-0 licenses.

Should you be unable or unwilling to comply with the dual-licensing
condition above, we are unable to accept your code for the project.
This particularly applies to code that has previously been licensed
with an incompatible license (just the Apache 2.0 license, just the MIT
license, any GPL license, etc).

Please do not submit code that has already been licensed with an
incompatible license.

## Documentation Note ##

All of the documents for the XCS repositories are based upon a common
template. For this reason, there may be substantial overlap between
documents.

## Licensing ##

Please read this file fully before copying or using the code anywhere.
Questions can be raised through creation of an 'issue' on GitHub.

The code written for this project is dual-licensed, under the MIT
and MIT-0 licenses.

This means that you can select and use one of these licenses, i.e. MIT
or MIT-0. You do not have to use both licenses at once.

Clear attribution is strongly preferred.

This project uses additional external packages. The code of these
external packages is clearly not covered by the above license(s). Each
package pulled in is subject to its own license.

As explained within the README, this implementation of the XCS
(eXtended Classifier System) algorithm is based on [Butz, M. V., &
Wilson, S. W. (2000, September). An algorithmic description of XCS. In
International Workshop on Learning Classifier Systems (pp. 253-272).
Springer, Berlin, Heidelberg]. Permission to use the algorithm
description within the above paper as a base for the present work was
sought and received from Springer (personal communication, 2020-07-07).

## References ##

Butz, M. V., & Wilson, S. W. (2000, September). An algorithmic
description of XCS. In International Workshop on Learning Classifier
Systems (pp. 253-272). Springer, Berlin, Heidelberg.

Urbanowicz, R. J., & Browne, W. N. (2017). Introduction to learning
classifier systems. Springer, Berlin, Heidelberg.

Wilson, S. W. (1998). Generalization in the XCS classifier system. In
J. R. Koza et al. (Eds.), Genetic Programming 1998: Proceedings of the
Third Annual Conference (pp. 665-674). Morgan Kaufmann.
