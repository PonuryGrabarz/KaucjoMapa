import { Result } from "../result";

export class GetNotificationsRequest {
  public readonly kind = "get-notifications" as const;

  private constructor() {}

  public static create(): Result<GetNotificationsRequest> {
    return Result.success(new GetNotificationsRequest());
  }
}
