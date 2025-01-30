Currently supported kinds (in term of [reflect] package) (including pointers and aliases): bool, string, int64, float64, struct and time.Duration type.
A double pointer returns an error without the possibility of traversal.
Any other kinds return errors, but can be skipped with provided [OptionContinueOnUnknownKind] option.

If the "efName" tag exists for a nested struct, it will be added as a prefix to the first level nested struct fields