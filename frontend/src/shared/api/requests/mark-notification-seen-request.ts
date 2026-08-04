import { z } from "zod";

import { Result } from "../result";

export interface MarkNotificationSeenParams {
  notification_id: string;
}

const SCHEMA = z.strictObject(
  {
    notification_id: z.uuid({
      error: "Powiadomienie: identyfikator musi być poprawnym UUID.",
    }),
  },
  { error: "Żądanie zawiera nieznany parametr." },
);

const toErrors = (error: z.ZodError): Error[] =>
  error.issues.map((issue) => new Error(issue.message));

export class MarkNotificationSeenRequest {
  public readonly kind = "mark-notification-seen" as const;

  private readonly params: Readonly<z.infer<typeof SCHEMA>>;

  private constructor(params: z.infer<typeof SCHEMA>) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: MarkNotificationSeenParams): Result<MarkNotificationSeenRequest> {
    const parsed = SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<MarkNotificationSeenRequest>(toErrors(parsed.error));

    return Result.success(new MarkNotificationSeenRequest(parsed.data));
  }

  public get notificationId(): string {
    return this.params.notification_id;
  }
}
