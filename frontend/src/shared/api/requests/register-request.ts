import { z } from "zod";

import { Result } from "../result";

export interface RegisterParams {
  username: string;
  email: string;
  password: string;
}

const SCHEMA = z.strictObject(
  {
    username: z
      .string({ error: "Nazwa użytkownika: wartość musi być tekstem." })
      .trim()
      .min(1, { error: "Nazwa użytkownika: pole jest wymagane." })
      .max(64, { error: "Nazwa użytkownika: dozwolone są najwyżej 64 znaki." }),
    email: z
      .email({ error: "Adres e-mail: wartość musi być poprawnym adresem." })
      .trim()
      .max(254, { error: "Adres e-mail: dozwolone są najwyżej 254 znaki." }),
    password: z
      .string({ error: "Hasło: wartość musi być tekstem." })
      .min(8, { error: "Hasło: wymagane jest co najmniej 8 znaków." })
      .max(256, { error: "Hasło: dozwolone jest najwyżej 256 znaków." }),
  },
  { error: "Żądanie zawiera nieznany parametr." },
);

const toErrors = (error: z.ZodError): Error[] =>
  error.issues.map((issue) => new Error(issue.message));

export class RegisterRequest {
  public readonly kind = "register" as const;

  private readonly params: Readonly<z.infer<typeof SCHEMA>>;

  private constructor(params: z.infer<typeof SCHEMA>) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: RegisterParams): Result<RegisterRequest> {
    const parsed = SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<RegisterRequest>(toErrors(parsed.error));

    return Result.success(new RegisterRequest(parsed.data));
  }

  public get username(): string {
    return this.params.username;
  }

  public get email(): string {
    return this.params.email;
  }

  public toBody(): Record<string, string> {
    return {
      username: this.params.username,
      email: this.params.email,
      password: this.params.password,
    };
  }
}
