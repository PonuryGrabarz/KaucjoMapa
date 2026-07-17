import Styles from "./register-form.module.sass";
import Button from "@components/button/button";
import FormField from "../form_field/form-field";

import { apiClient } from "@api/api-client-factory";
import { RegisterRequest } from "@api/requests";

import { actions } from "astro:actions";

import SocialAuth, {
  type SocialAuthPlatform,
} from "../social_auth/social-auth";

import { useState } from "preact/hooks";

import type { TargetedSubmitEvent } from "preact";

type RegisterField =
  "username" | "email" | "password" | "confirmPassword" | "terms";
type RegisterFieldErrors = Partial<Record<RegisterField, string>>;

const mapValidationErrors = (errors: readonly Error[]): RegisterFieldErrors => {
  const fieldErrors: RegisterFieldErrors = {};

  for (const error of errors) {
    if (error.message === "Nazwa użytkownika jest wymagana.") {
      fieldErrors.username = error.message;
    } else if (error.message.startsWith("Adres e-mail")) {
      fieldErrors.email = error.message;
    } else if (error.message === "Hasło jest wymagane.") {
      fieldErrors.password = error.message;
    } else if (
      error.message === "Potwierdzenie hasła jest wymagane." ||
      error.message === "Hasła muszą być takie same."
    ) {
      fieldErrors.confirmPassword = error.message;
    }
  }

  return fieldErrors;
};

export default function RegisterForm() {
  const [fieldErrors, setFieldErrors] = useState<RegisterFieldErrors>({});
  const [formError, setFormError] = useState<string | undefined>(undefined);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const clearFieldError = (field: RegisterField): void => {
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
    // Placeholder do czasu udostępnienia endpointu rejestracji społecznościowej.
    setFormError("Rejestracja przez Google nie jest jeszcze dostępna.");
  };

  const handleSubmit = async (
    event: TargetedSubmitEvent<HTMLFormElement>,
  ): Promise<void> => {
    event.preventDefault();

    const formData = new FormData(event.currentTarget);
    const username = String(formData.get("username") ?? "").trim();
    const email = String(formData.get("email") ?? "").trim();
    const password = String(formData.get("password") ?? "");
    const confirmPassword = String(formData.get("confirmPassword") ?? "");
    const acceptedTerms = formData.get("terms") === "on";
    const requestResult = RegisterRequest.create(
      username,
      email,
      password,
      confirmPassword,
    );
    const validationErrors = requestResult.isError
      ? mapValidationErrors(requestResult.errors)
      : {};

    if (!acceptedTerms) {
      validationErrors.terms = "Zaakceptuj regulamin i politykę prywatności.";
    }

    if (requestResult.isError || !acceptedTerms) {
      setFieldErrors(validationErrors);
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
          : "Nie udało się utworzyć konta. Spróbuj ponownie później.",
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
      aria-describedby={formError ? "register-form-error" : undefined}
    >
      <header class={Styles["auth-form__header"]}>
        <h1 id="register-form-title">Załóż konto</h1>
        <p>Odkrywaj punkty zwrotu w swojej okolicy.</p>
      </header>

      <SocialAuth type="register" onAuthenticate={handleSocialAuth} />

      <div class={Styles["auth-form__divider"]}>albo danymi konta</div>

      <FormField
        id="register-username"
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
        id="register-email"
        label="Adres e-mail"
        name="email"
        type="email"
        inputMode="email"
        placeholder="np. jan.kowalski@example.com"
        autocomplete="email"
        required
        error={fieldErrors.email}
        onInput={() => clearFieldError("email")}
      />

      <FormField
        id="register-password"
        label="Hasło"
        name="password"
        type="password"
        showPasswordToggle
        placeholder="Utwórz hasło"
        autocomplete="new-password"
        required
        error={fieldErrors.password}
        onInput={() => clearFieldError("password")}
      />

      <FormField
        id="register-confirm-password"
        label="Powtórz hasło"
        name="confirmPassword"
        type="password"
        showPasswordToggle
        placeholder="Powtórz hasło"
        autocomplete="new-password"
        required
        error={fieldErrors.confirmPassword}
        onInput={() => clearFieldError("confirmPassword")}
      />

      <div
        class={Styles["auth-form__terms"]}
        data-error={fieldErrors.terms ? "true" : undefined}
      >
        <label class={Styles["auth-form__checkbox"]}>
          <input
            name="terms"
            type="checkbox"
            required
            aria-invalid={fieldErrors.terms ? "true" : undefined}
            aria-describedby={
              fieldErrors.terms ? "register-terms-error" : undefined
            }
            onInput={() => clearFieldError("terms")}
          />
          <span>Akceptuję regulamin i politykę prywatności</span>
        </label>

        {fieldErrors.terms && (
          <span
            id="register-terms-error"
            class={Styles["auth-form__field-error"]}
            role="alert"
          >
            {fieldErrors.terms}
          </span>
        )}
      </div>

      {formError && (
        <p
          id="register-form-error"
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
        Załóż konto
      </Button>

      <p class={Styles["auth-form__switch-line"]}>
        Masz już konto? <a href="/logowanie">Zaloguj się</a>
      </p>
    </form>
  );
}
