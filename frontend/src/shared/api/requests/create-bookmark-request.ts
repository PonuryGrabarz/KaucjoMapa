import { z } from "zod";

import { Result } from "../result";

export interface CreateBookmarkParams {
  return_point_id: string;
}

const SCHEMA = z.strictObject(
  {
    return_point_id: z.uuid({
      error: "Punkt zwrotu: identyfikator musi być poprawnym UUID.",
    }),
  },
  { error: "Żądanie zawiera nieznany parametr." },
);

const toErrors = (error: z.ZodError): Error[] =>
  error.issues.map((issue) => new Error(issue.message));

export class CreateBookmarkRequest {
  public readonly kind = "create-bookmark" as const;

  private readonly params: Readonly<z.infer<typeof SCHEMA>>;

  private constructor(params: z.infer<typeof SCHEMA>) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: CreateBookmarkParams): Result<CreateBookmarkRequest> {
    const parsed = SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<CreateBookmarkRequest>(toErrors(parsed.error));

    return Result.success(new CreateBookmarkRequest(parsed.data));
  }

  public get returnPointId(): string {
    return this.params.return_point_id;
  }

  public toBody(): Record<string, string> {
    return { return_point_id: this.params.return_point_id };
  }
}
