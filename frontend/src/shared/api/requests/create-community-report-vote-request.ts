import { z } from "zod";

import { Result } from "../result";
import { COMMUNITY_REPORT_VOTE_VALUES } from "../responses";

import type { CommunityReportVoteValue } from "../responses";

export interface CreateCommunityReportVoteParams {
  report_id: string;
  vote: CommunityReportVoteValue;
}

const SCHEMA = z.strictObject(
  {
    report_id: z.uuid({ error: "Zgłoszenie: identyfikator musi być poprawnym UUID." }),
    vote: z.enum(COMMUNITY_REPORT_VOTE_VALUES, {
      error: `Głos musi być jedną z wartości: ${COMMUNITY_REPORT_VOTE_VALUES.join(", ")}.`,
    }),
  },
  { error: "Żądanie zawiera nieznany parametr." },
);

const toErrors = (error: z.ZodError): Error[] =>
  error.issues.map((issue) => new Error(issue.message));

export class CreateCommunityReportVoteRequest {
  public readonly kind = "create-community-report-vote" as const;

  private readonly params: Readonly<z.infer<typeof SCHEMA>>;

  private constructor(params: z.infer<typeof SCHEMA>) {
    this.params = Object.freeze({ ...params });
  }

  public static create(
    params: CreateCommunityReportVoteParams,
  ): Result<CreateCommunityReportVoteRequest> {
    const parsed = SCHEMA.safeParse(params);

    if (!parsed.success)
      return Result.failure<CreateCommunityReportVoteRequest>(toErrors(parsed.error));

    return Result.success(new CreateCommunityReportVoteRequest(parsed.data));
  }

  public get reportId(): string {
    return this.params.report_id;
  }

  public toBody(): Record<string, string> {
    return { vote: this.params.vote };
  }
}
