# Composition Functions (xfn) for Crossplane

The composition function tutorial presented in KubeCon Amsterdam 2023 as part of
[Crossplane ContribFest](https://kccnceu2023.sched.com/event/1Hzcf).

This tutorial is a step-by-step guide to create composition functions that can
do things that would not be possible with the standard Composition such as:
* Creating resources conditionally,
* Generating a random string only once,
* Initialize infrastructure like creating a database schema or setting timezone
  for a database,
* Manipulating JSON values such as IAM policies in smarter ways.

## Index

* [Prerequisites](prerequisites.md)
* Create a no-op function.
* Create a function that creates firewall rule resources conditionally to open
  up the database to public.