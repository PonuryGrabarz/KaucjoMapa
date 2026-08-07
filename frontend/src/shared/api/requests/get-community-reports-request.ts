import { z } from "zod";

import { Result } from "../result";

export interface GetCommunityReportsParams {
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

export class GetCommunityReportsRequest {
  public readonly kind = "get-community-reports" as const;

  private readonly params: Readonly<z.infer<typeof SCHEMA>>;

  private constructor(params: z.infer<typeof SCHEMA>) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: GetCommunityReportsParams): Result<GetCommunityReportsRequest> {
    const parsed = SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<GetCommunityReportsRequest>(toErrors(parsed.error));

    return Result.success(new GetCommunityReportsRequest(parsed.data));
  }

  public get returnPointId(): string {
    return this.params.return_point_id;
  }
}
