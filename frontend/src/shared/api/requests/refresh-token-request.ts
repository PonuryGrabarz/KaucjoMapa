import { z } from "zod";

import { Result } from "../result";

export interface RefreshTokenParams {
  refresh_token: string;
}

const SCHEMA = z.strictObject(
  {
    refresh_token: z
      .string({ error: "Token odświeżający: wartość musi być tekstem." })
      .trim()
      .min(1, { error: "Token odświeżający: pole jest wymagane." }),
  },
  { error: "Żądanie zawiera nieznany parametr." },
);

const toErrors = (error: z.ZodError): Error[] =>
  error.issues.map((issue) => new Error(issue.message));

export class RefreshTokenRequest {
  public readonly kind = "refresh-token" as const;

  private readonly params: Readonly<z.infer<typeof SCHEMA>>;

  private constructor(params: z.infer<typeof SCHEMA>) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: RefreshTokenParams): Result<RefreshTokenRequest> {
    const parsed = SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<RefreshTokenRequest>(toErrors(parsed.error));

    return Result.success(new RefreshTokenRequest(parsed.data));
  }

  public get refreshToken(): string {
    return this.params.refresh_token;
  }
}
