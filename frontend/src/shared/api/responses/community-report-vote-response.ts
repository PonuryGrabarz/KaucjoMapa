export const COMMUNITY_REPORT_VOTE_VALUES = ["confirm", "reject"] as const;

export type CommunityReportVoteValue = (typeof COMMUNITY_REPORT_VOTE_VALUES)[number];

export class CreateCommunityReportVoteResponse {
  CommunityReportID!: string;
  UserID!: string;
  Vote!: CommunityReportVoteValue;
  CreatedAt!: string;
  UpdatedAt!: string;
}
