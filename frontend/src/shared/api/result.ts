type NonEmptyErrors = readonly [Error, ...Error[]];

type ResultState<TResult> =
  | {
      readonly isError: false;
      readonly value: TResult;
    }
  | {
      readonly isError: true;
      readonly errors: NonEmptyErrors;
    };

const NO_ERRORS: readonly Error[] = Object.freeze([]);

export class Result<TResult> {
  private readonly state: ResultState<TResult>;

  private constructor(state: ResultState<TResult>) {
    this.state = state;
  }

  public static success<TResult>(value: TResult): Result<TResult> {
    return new Result<TResult>({ isError: false, value });
  }

  public static failure<TResult = never>(error: Error): Result<TResult>;

  public static failure<TResult = never>(errors: readonly Error[]): Result<TResult>;

  public static failure<TResult = never>(errors: Error | readonly Error[]): Result<TResult> {
    const normalizedErrors = errors instanceof Error ? [errors] : [...errors];

    if (normalizedErrors.length === 0)
      throw new Error("A failed Result must contain at least one error.");

    return new Result<TResult>({
      isError: true,
      errors: Object.freeze(normalizedErrors) as NonEmptyErrors,
    });
  }

  public get isError(): boolean {
    return this.state.isError;
  }

  public get isSuccess(): boolean {
    return !this.state.isError;
  }

  public get value(): TResult {
    if (this.state.isError) throw new Error(
        "Cannot access the value of a failed Result. Check isError first.",
      );

    return this.state.value;
  }

  public get errors(): NonEmptyErrors {
    if (!this.state.isError) throw new Error(
        "Cannot access errors of a successful Result. Check isError first.",
      );

    return this.state.errors;
  }

  public get errorsOrEmptyList(): readonly Error[] {
    return this.state.isError ? this.state.errors : NO_ERRORS;
  }

  public get firstError(): Error {
    return this.errors[0];
  }

  public match<TMatch>(
    onSuccess: (value: TResult) => TMatch,
    onFailure: (errors: NonEmptyErrors) => TMatch,
  ): TMatch {
    return this.state.isError
      ? onFailure(this.state.errors)
      : onSuccess(this.state.value);
  }

  public switch(
    onSuccess: (value: TResult) => void,
    onFailure: (errors: NonEmptyErrors) => void,
  ): void {
    if (this.state.isError) {
      onFailure(this.state.errors);
      return;
    }

    onSuccess(this.state.value);
  }
}
