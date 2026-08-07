import { z } from "zod";

import { Result } from "../result";
import { COMMUNITY_REPORT_STATUSES } from "../responses";

import type { CommunityReportStatus } from "../responses";

export interface CreateCommunityReportParams {
  return_point_id: string;
  description: string;
  status: CommunityReportStatus;
}

const FIELD_LABELS: Record<string, string> = {
  return_point_id: "Punkt zwrotu",
  description: "Opis zgłoszenia",
  status: "Status punktu",
};

const listOptions = (values: readonly string[]): string => values.join(", ");

const CREATE_COMMUNITY_REPORT_SCHEMA = z.strictObject(
  {
    return_point_id: z
      .uuid({ error: "Identyfikator punktu zwrotu musi być poprawnym UUID." }),

    status: z.enum(COMMUNITY_REPORT_STATUSES, {
      error: `Status punktu musi być jedną z wartości: ${listOptions(COMMUNITY_REPORT_STATUSES)}.`,
    }),

    description: z
      .string({ error: "Opis zgłoszenia musi być tekstem." })
      .trim()
      .min(1, { error: "Opis zgłoszenia nie może być pusty." }),
  },
  { error: "Żądanie zawiera nieznany parametr." },
);

type ValidatedParams = z.infer<typeof CREATE_COMMUNITY_REPORT_SCHEMA>;

const toErrors = (error: z.ZodError): Error[] =>
  error.issues.map((issue) => {
    const field = issue.path[0];

    const label = typeof field === "string" ? FIELD_LABELS[field] : undefined;

    return new Error(
      label && !issue.message.startsWith(label) ? `${label}: ${issue.message}` : issue.message,
    );
  });

export class CreateCommunityReportRequest {
  public readonly kind = "create-community-report" as const;

  private readonly params: Readonly<ValidatedParams>;

  private constructor(params: ValidatedParams) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: CreateCommunityReportParams): Result<CreateCommunityReportRequest> {
    const parsed = CREATE_COMMUNITY_REPORT_SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<CreateCommunityReportRequest>(toErrors(parsed.error));

    return Result.success(new CreateCommunityReportRequest(parsed.data));
  }

  public get returnPointId(): string {
    return this.params.return_point_id;
  }

  public get status(): CommunityReportStatus {
    return this.params.status;
  }

  public get description(): string {
    return this.params.description;
  }

  public toBody(): Record<string, string> {
    return {
      report_id: this.params.return_point_id,
      status: this.params.status,
      description: this.params.description,
    };
  }
}
