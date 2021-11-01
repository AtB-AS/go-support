# Logging functions

This module contains functions to make the [Zerolog](https://github.com/rs/zerolog)
work well with GCP and AMP, and functions for setting parameters for 
[contextual logging](https://github.com/rs/zerolog#contextual-logging).

To make Zerolog format log entries correctly for GCP, call `ConfigureGcpLogging` as early as possible in 
application / function startup. Use the `XxxContext` functions to set the given context when it is available.
