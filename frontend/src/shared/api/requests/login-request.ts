import { z } from "zod";

import { Result } from "../result";

export interface LoginParams {
  username: string;
  password: string;
}

const SCHEMA = z.strictObject(
  {
    username: z
      .string({ error: "Nazwa użytkownika: wartość musi być tekstem." })
      .trim()
      .min(1, { error: "Nazwa użytkownika: pole jest wymagane." }),
    password: z
      .string({ error: "Hasło: wartość musi być tekstem." })
      .min(1, { error: "Hasło: pole jest wymagane." }),
  },
  { error: "Żądanie zawiera nieznany parametr." },
);

const toErrors = (error: z.ZodError): Error[] =>
  error.issues.map((issue) => new Error(issue.message));

export class LoginRequest {
  public readonly kind = "login" as const;

  private readonly params: Readonly<z.infer<typeof SCHEMA>>;

  private constructor(params: z.infer<typeof SCHEMA>) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: LoginParams): Result<LoginRequest> {
    const parsed = SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<LoginRequest>(toErrors(parsed.error));

    return Result.success(new LoginRequest(parsed.data));
  }

  public get username(): string {
    return this.params.username;
  }

  public toBody(): Record<string, string> {
    return { username: this.params.username, password: this.params.password };
  }
}
