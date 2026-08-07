import { Result } from "../result";

export class GetAvatarRequest {
  public readonly kind = "get-avatar" as const;

  private constructor() {}

  public static create(): Result<GetAvatarRequest> {
    return Result.success(new GetAvatarRequest());
  }
}
