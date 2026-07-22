import Styles from "./logout-button.module.sass";
import Button from "@components/button/button";

import { actions } from "astro:actions";
import { useState } from "preact/hooks";

const DEFAULT_ERROR = "Nie udało się wylogować. Spróbuj ponownie później.";

export default function LogoutButton() {
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [error, setError] = useState<string | undefined>(undefined);

  const handleLogout = async (): Promise<void> => {
    if (isLoggingOut) return;

    setError(undefined);
    setIsLoggingOut(true);

    try {
      const { data, error: actionError } = await actions.auth.destroySession();

      if (actionError) throw new Error(actionError.message);
      if (!data?.success) throw new Error(DEFAULT_ERROR);

      window.location.replace("/logowanie");
    } catch (logoutError) {
      setError(logoutError instanceof Error ? logoutError.message : DEFAULT_ERROR);
      setIsLoggingOut(false);
    }
  };

  return (
    <div class={Styles["logout-button"]}>
      <Button
        class={Styles["logout-button__action"]}
        type="button"
        variant="contained"
        startIcon="mdi:logout"
        fullWidth
        loading={isLoggingOut}
        loadingPosition="start"
        aria-describedby={error ? "logout-button-error" : undefined}
        onClick={handleLogout}
      >
        {isLoggingOut ? "Wylogowywanie…" : "Wyloguj się"}
      </Button>

      {error && (
        <p id="logout-button-error" class={Styles["logout-button__error"]} role="alert">
          {error}
        </p>
      )}
    </div>
  );
}
