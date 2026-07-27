import { Result } from "../result";

export class LoginRequest {
  private readonly requestType = "login";

  readonly username: string;
  readonly password: string;

  private constructor(username: string, password: string) {
    this.username = username;
    this.password = password;
  }

  public static create(
    username: string,
    password: string,
  ): Result<LoginRequest> {
    const errors: Error[] = [];

    if (username.trim().length === 0) {
      errors.push(new Error("Nazwa użytkownika jest wymagana."));
    }

    if (password.trim().length === 0) {
      errors.push(new Error("Hasło jest wymagane."));
    }

    if (errors.length > 0) {
      return Result.failure<LoginRequest>(errors);
    }

    return Result.success(new LoginRequest(username, password));
  }
}
