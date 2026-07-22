import Styles from "./alerts-summary.module.sass";
import Button from "@components/button/button";

interface Props {
    listId: string;
    unreadCount: number;
    onMarkAllRead: () => void;
}

export default function AlertsSummary({ listId, unreadCount, onMarkAllRead }: Props) {
    const unreadLabel = unreadCount === 1 ? "1 nowe" : `${unreadCount} nowych`;

    return (
        <header class={Styles["alerts-summary"]}>
            <div class={Styles["alerts-summary__copy"]}>
                <p class={Styles["alerts-summary__count"]} aria-live="polite" aria-atomic="true">
                    <strong>{unreadLabel}</strong>
                </p>

                <p class={Styles["alerts-summary__description"]}>
                    Nieprzeczytane powiadomienia są oznaczone zieloną kropką.
                </p>
            </div>

            <Button
                class={Styles["alerts-summary__button"]}
                variant="outlined"
                type="button"
                aria-controls={listId}
                disabled={unreadCount === 0}
                onClick={onMarkAllRead}
            >
                {unreadCount === 0
                    ? "Brak nieprzeczytanych"
                    : "Oznacz wszystkie jako przeczytane"}
            </Button>
        </header>
    );
}
