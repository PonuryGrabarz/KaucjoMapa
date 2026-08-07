export const COMMUNITY_REPORT_STATUSES = [
  "dostepny",
  "nieczynny",
  "brak_informacji",
] as const;

export type CommunityReportStatus = (typeof COMMUNITY_REPORT_STATUSES)[number];

export class CommunityReport {
  id!: string;
  user_id!: string;
  return_point_id!: string;
  status!: CommunityReportStatus;
  description!: string;
  confirmations_count!: number;
  rejections_count!: number;
  created_at!: string;
  updated_at!: string;
}

export class CreateCommunityReportResponse {
  report!: CommunityReport;
  reputation!: number;
  username!: string;
}
