import { z } from "zod";

import { Result } from "../result";

export interface GetBookmarkByReturnPointParams {
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

export class GetBookmarkByReturnPointRequest {
  public readonly kind = "get-bookmark-by-return-point" as const;

  private readonly params: Readonly<z.infer<typeof SCHEMA>>;

  private constructor(params: z.infer<typeof SCHEMA>) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: GetBookmarkByReturnPointParams): Result<GetBookmarkByReturnPointRequest> {
    const parsed = SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<GetBookmarkByReturnPointRequest>(toErrors(parsed.error));

    return Result.success(new GetBookmarkByReturnPointRequest(parsed.data));
  }

  public get returnPointId(): string {
    return this.params.return_point_id;
  }
}
