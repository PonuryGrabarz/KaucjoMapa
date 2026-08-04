import { z } from "zod";

import { Result } from "../result";

export interface CreateNotificationParams {
  name: string;
  description: string;
}

const SCHEMA = z.strictObject(
  {
    name: z
      .string({ error: "Tytuł powiadomienia: wartość musi być tekstem." })
      .trim()
      .min(1, { error: "Tytuł powiadomienia: pole jest wymagane." })
      .max(200, { error: "Tytuł powiadomienia: dozwolone są najwyżej 200 znaków." }),
    description: z
      .string({ error: "Treść powiadomienia: wartość musi być tekstem." })
      .trim()
      .min(1, { error: "Treść powiadomienia: pole jest wymagane." })
      .max(2000, { error: "Treść powiadomienia: dozwolone są najwyżej 2000 znaków." }),
  },
  { error: "Żądanie zawiera nieznany parametr." },
);

const toErrors = (error: z.ZodError): Error[] =>
  error.issues.map((issue) => new Error(issue.message));

export class CreateNotificationRequest {
  public readonly kind = "create-notification" as const;

  private readonly params: Readonly<z.infer<typeof SCHEMA>>;

  private constructor(params: z.infer<typeof SCHEMA>) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: CreateNotificationParams): Result<CreateNotificationRequest> {
    const parsed = SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<CreateNotificationRequest>(toErrors(parsed.error));

    return Result.success(new CreateNotificationRequest(parsed.data));
  }

  public toBody(): Record<string, string> {
    return { name: this.params.name, description: this.params.description };
  }
}
