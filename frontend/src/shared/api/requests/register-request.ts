import { Result } from "../result";

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export class RegisterRequest {
  private readonly requestType = "register";

  readonly username: string;
  readonly email: string;
  readonly password: string;
  readonly confirmPassword: string;

  private constructor(
    username: string,
    email: string,
    password: string,
    confirmPassword: string,
  ) {
    this.username = username;
    this.email = email;
    this.password = password;
    this.confirmPassword = confirmPassword;
  }

  public static create(
    username: string,
    email: string,
    password: string,
    confirmPassword: string,
  ): Result<RegisterRequest> {
    const errors: Error[] = [];

    if (username.trim().length === 0)
      errors.push(new Error("Nazwa użytkownika jest wymagana."));

    if (email.trim().length === 0)
      errors.push(new Error("Adres e-mail jest wymagany."));
    else if (!EMAIL_PATTERN.test(email))
      errors.push(new Error("Adres e-mail ma nieprawidłowy format."));

    if (password.trim().length === 0)
      errors.push(new Error("Hasło jest wymagane."));

    if (confirmPassword.trim().length === 0)
      errors.push(new Error("Potwierdzenie hasła jest wymagane."));
    else if (password !== confirmPassword)
      errors.push(new Error("Hasła muszą być takie same."));

    if (errors.length > 0) return Result.failure<RegisterRequest>(errors);

    return Result.success(
      new RegisterRequest(username, email, password, confirmPassword),
    );
  }
}
