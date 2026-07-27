import Styles from "./social-auth.module.sass";
import Button from "@components/button/button";

type SocialAuthProps = {
  type: "login" | "register";
  onAuthenticate?: (platform: SocialAuthPlatform) => void;
};

export type SocialAuthPlatform = "google";

export default function SocialAuth({ type, onAuthenticate }: SocialAuthProps) {
  const handleSocialAuth = (platform: SocialAuthPlatform): void => {
    onAuthenticate?.(platform);
  };

  return (
    <Button
      class={Styles["social-auth__btn"]}
      type="button"
      variant="contained"
      startIcon="material-icon-theme:google"
      startIconColored
      fullWidth
      onClick={() => handleSocialAuth("google")}
    >
      {type === "login" ? "Kontynuuj z Google" : "Zarejestruj się przez Google"}
    </Button>
  );
}
