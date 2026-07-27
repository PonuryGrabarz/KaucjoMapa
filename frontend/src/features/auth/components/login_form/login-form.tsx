import Styles from "./login-form.module.sass";
import Button from "@components/button/button";
import FormField from "../form_field/form-field";

import { apiClient } from "@api/api-client-factory";
import { LoginRequest } from "@api/requests";

import { actions } from "astro:actions";

import SocialAuth, {
  type SocialAuthPlatform,
} from "../social_auth/social-auth";

import { useState } from "preact/hooks";

import type { TargetedSubmitEvent } from "preact";

type LoginField = "username" | "password";
type LoginFieldErrors = Partial<Record<LoginField, string>>;

const mapValidationErrors = (errors: readonly Error[]): LoginFieldErrors => {
  const fieldErrors: LoginFieldErrors = {};

  for (const error of errors) {
    if (error.message === "Nazwa użytkownika jest wymagana.") {
      fieldErrors.username = error.message;
    } else if (error.message === "Hasło jest wymagane.") {
      fieldErrors.password = error.message;
    }
  }

  return fieldErrors;
};

export default function LoginForm() {
  const [fieldErrors, setFieldErrors] = useState<LoginFieldErrors>({});
  const [formError, setFormError] = useState<string | undefined>(undefined);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const clearFieldError = (field: LoginField): void => {
    setFieldErrors((currentErrors) => {
      if (!currentErrors[field]) return currentErrors;

      const nextErrors = { ...currentErrors };
      delete nextErrors[field];
      return nextErrors;
    });
    setFormError(undefined);
  };

  const handleSocialAuth = (platform: SocialAuthPlatform): void => {
    void platform;
    // Placeholder do czasu udostępnienia endpointu logowania społecznościowego.
    setFormError("Logowanie przez Google nie jest jeszcze dostępne.");
  };

  const handleSubmit = async (
    event: TargetedSubmitEvent<HTMLFormElement>,
  ): Promise<void> => {
    event.preventDefault();

    const formData = new FormData(event.currentTarget);
    const username = String(formData.get("username") ?? "").trim();
    const password = String(formData.get("password") ?? "");
    const requestResult = LoginRequest.create(username, password);

    if (requestResult.isError) {
      setFieldErrors(mapValidationErrors(requestResult.errors));
      setFormError(undefined);
      return;
    }

    setFieldErrors({});
    setFormError(undefined);
    setIsSubmitting(true);

    try {
      const responseResult = await apiClient.send(requestResult.value);

      if (responseResult.isError) {
        setFormError(responseResult.firstError.message);
        return;
      }

      const response = responseResult.value;
      const { error: sessionError } = await actions.auth.createSession({
        createdAt: response.createdAt,
        email: response.email,
        id: response.id,
        isAdmin: response.isAdmin,
        isBanned: response.isBanned,
        reputation: response.reputation,
        token: response.token.rawValue,
        updatedAt: response.updatedAt,
        username: response.username,
      });

      if (sessionError) {
        throw new Error(sessionError.message);
      }

      window.location.replace("/");
    } catch (error) {
      setFormError(
        error instanceof Error
          ? error.message
          : "Nie udało się zalogować. Spróbuj ponownie później.",
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form
      class={Styles["auth-form"]}
      noValidate
      onSubmit={handleSubmit}
      aria-busy={isSubmitting ? "true" : undefined}
      aria-describedby={formError ? "login-form-error" : undefined}
    >
      <header class={Styles["auth-form__header"]}>
        <h1 id="login-form-title">Zaloguj się</h1>
        <p>Wróć do swojej mapy punktów zwrotu.</p>
      </header>

      <SocialAuth type="login" onAuthenticate={handleSocialAuth} />

      <div class={Styles["auth-form__divider"]}>albo danymi konta</div>

      <FormField
        id="login-username"
        label="Nazwa użytkownika"
        name="username"
        type="text"
        placeholder="Wpisz nazwę użytkownika"
        autocomplete="username"
        required
        error={fieldErrors.username}
        onInput={() => clearFieldError("username")}
      />

      <FormField
        id="login-password"
        label="Hasło"
        name="password"
        type="password"
        showPasswordToggle
        placeholder="Wpisz hasło"
        autocomplete="current-password"
        required
        error={fieldErrors.password}
        onInput={() => clearFieldError("password")}
      />

      <div class={Styles["auth-form__row-between"]}>
        <label class={Styles["auth-form__checkbox"]}>
          <input name="remember" type="checkbox" />
          <span>Zapamiętaj mnie</span>
        </label>

        <span class={Styles["auth-form__soon"]}>
          Odzyskiwanie hasła wkrótce
        </span>
      </div>

      {formError && (
        <p
          id="login-form-error"
          class={Styles["auth-form__error"]}
          role="alert"
        >
          {formError}
        </p>
      )}

      <Button
        class={Styles["auth-form__submit"]}
        type="submit"
        variant="contained"
        fullWidth
        loading={isSubmitting}
      >
        Zaloguj się
      </Button>

      <p class={Styles["auth-form__switch-line"]}>
        Nie masz konta? <a href="/rejestracja">Zarejestruj się</a>
      </p>
    </form>
  );
}
